package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding    = 2
	maxWidth   = 80
	listHeight = 14
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	listhelpStyle     = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

/* progress bar */
var p *tea.Program

type progressWriter struct {
	total      int
	downloaded int
	file       *os.File
	reader     io.Reader
	onProgress func(float64)
}

func (pw *progressWriter) Start() {
	// TeeReader calls pw.Write() each time a new response is received
	if _, err := io.Copy(pw.file, io.TeeReader(pw.reader, pw)); err != nil {
		p.Send(progressErrMsg{err})
	}
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.downloaded += len(p)
	if pw.total > 0 && pw.onProgress != nil {
		pw.onProgress(float64(pw.downloaded) / float64(pw.total))
	}
	return len(p), nil
}

type progressMsg float64

type progressErrMsg struct{ err error }

func finalPause() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg {
		return nil
	})
}

/* list */

type item struct {
	Name string
	URL  string
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

/* model */

type model struct {
	list     list.Model
	choice   item
	quitting bool

	formula  *Formula
	pw       *progressWriter
	progress progress.Model
	err      error
}

func initialModel(formula *Formula) model {
	var items []list.Item

	if formula.Bottle.Stable.Files.Arm64Sonoma.URL != "" {
		items = append(items, item{
			Name: "macOS Sonoma (arm64)",
			URL:  formula.Bottle.Stable.Files.Arm64Sonoma.URL,
		})
	}
	if formula.Bottle.Stable.Files.Arm64Ventura.URL != "" {
		items = append(items, item{
			Name: "macOS Ventura (arm64)",
			URL:  formula.Bottle.Stable.Files.Arm64Ventura.URL,
		})
	}
	if formula.Bottle.Stable.Files.Arm64Monterey.URL != "" {
		items = append(items, item{
			Name: "macOS Monterey (arm64)",
			URL:  formula.Bottle.Stable.Files.Arm64Monterey.URL,
		})
	}
	if formula.Bottle.Stable.Files.Sonoma.URL != "" {
		items = append(items, item{
			Name: "macOS Sonoma (x86_64)",
			URL:  formula.Bottle.Stable.Files.Sonoma.URL,
		})
	}
	if formula.Bottle.Stable.Files.Ventura.URL != "" {
		items = append(items, item{
			Name: "macOS Ventura (x86_64)",
			URL:  formula.Bottle.Stable.Files.Ventura.URL,
		})
	}
	if formula.Bottle.Stable.Files.Monterey.URL != "" {
		items = append(items, item{
			Name: "macOS Monterey (x86_64)",
			URL:  formula.Bottle.Stable.Files.Monterey.URL,
		})
	}
	if formula.Bottle.Stable.Files.Arm64Linux.URL != "" {
		items = append(items, item{
			Name: "Linux (arm64)",
			URL:  formula.Bottle.Stable.Files.Arm64Linux.URL,
		})
	}
	if formula.Bottle.Stable.Files.X8664Linux.URL != "" {
		items = append(items, item{
			Name: "Linux (x86_64)",
			URL:  formula.Bottle.Stable.Files.X8664Linux.URL,
		})
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = fmt.Sprintf("Download '%s' ?", formula.Name)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = listhelpStyle

	return model{
		formula: formula,
		list:    l,
		// pw:       pw,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i
			}
			return m, m.downloadBottle
		}

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case progressErrMsg:
		m.err = msg.err
		return m, tea.Quit

	case progressMsg:
		var cmds []tea.Cmd

		if msg >= 1.0 {
			cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
		}

		cmds = append(cmds, m.progress.SetPercent(float64(msg)))
		return m, tea.Batch(cmds...)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return "Error downloading: " + m.err.Error() + "\n"
	}
	if m.choice.Name != "" {
		pad := strings.Repeat(" ", padding)
		return "\n" +
			pad + m.progress.View() + "\n\n" +
			pad + helpStyle("Press any key to quit")
	}
	if m.quitting {
		return quitTextStyle.Render("🍺 Bottle dud? That’s cool.")
	}
	return "\n" + m.list.View()
}

func (m model) downloadBottle() tea.Msg {
	req, err := http.NewRequest("GET", m.choice.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer QQ==")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(m.formula.Name + ".tar.gz")
	if err != nil {
		return err
	}
	defer f.Close()

	m.pw = &progressWriter{
		total:  int(resp.ContentLength),
		file:   f,
		reader: resp.Body,
		onProgress: func(ratio float64) {
			p.Send(progressMsg(ratio))
		},
	}
	// Start the download
	m.pw.Start()

	return progressMsg(100.0)
}
