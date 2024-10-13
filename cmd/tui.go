package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh" // Add the 'huh' package
	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 100

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type state int

const (
	statusNormal state = iota
	stateDownloading
	stateQuitting
	stateDone
)

/* progress bar */

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

/* model */

type Model struct {
	state  state
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form // Replace 'list' with 'form'
	width  int

	selectedURL string // To store the selected URL
	formula     *Formula
	err         error

	pw       *progressWriter
	progress progress.Model
	created  bool
}

func initialModel(formula *Formula) Model {
	m := Model{
		formula: formula,
	}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	var options []huh.Option[string]

	// Build the options from the available files
	if formula.Bottle.Stable.Files.Arm64Sonoma.URL != "" {
		options = append(options, huh.NewOption("macOS Sonoma (arm64)", formula.Bottle.Stable.Files.Arm64Sonoma.URL))
	}
	if formula.Bottle.Stable.Files.Arm64Ventura.URL != "" {
		options = append(options, huh.NewOption("macOS Ventura (arm64)", formula.Bottle.Stable.Files.Arm64Ventura.URL))
	}
	if formula.Bottle.Stable.Files.Arm64Monterey.URL != "" {
		options = append(options, huh.NewOption("macOS Monterey (arm64)", formula.Bottle.Stable.Files.Arm64Monterey.URL))
	}
	if formula.Bottle.Stable.Files.Sonoma.URL != "" {
		options = append(options, huh.NewOption("macOS Sonoma (x86_64)", formula.Bottle.Stable.Files.Sonoma.URL))
	}
	if formula.Bottle.Stable.Files.Ventura.URL != "" {
		options = append(options, huh.NewOption("macOS Ventura (x86_64)", formula.Bottle.Stable.Files.Ventura.URL))
	}
	if formula.Bottle.Stable.Files.Monterey.URL != "" {
		options = append(options, huh.NewOption("macOS Monterey (x86_64)", formula.Bottle.Stable.Files.Monterey.URL))
	}
	if formula.Bottle.Stable.Files.Arm64Linux.URL != "" {
		options = append(options, huh.NewOption("Linux (arm64)", formula.Bottle.Stable.Files.Arm64Linux.URL))
	}
	if formula.Bottle.Stable.Files.X8664Linux.URL != "" {
		options = append(options, huh.NewOption("Linux (x86_64)", formula.Bottle.Stable.Files.X8664Linux.URL))
	}

	// Create the form
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("'%s' Bottles", formula.Name)).
				Options(options...).
				Value(&m.selectedURL),
		),
	).
		WithWidth(30).
		WithShowHelp(false).
		WithShowErrors(false)

	m.progress = progress.New(
		progress.WithDefaultGradient(),
		// progress.WithWidth(50),
	)

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.form.Init(),
		m.progress.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
		m.progress.Width = m.width - 10 // Adjust width as needed
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			m.state = stateQuitting
			return m, tea.Quit
		}
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
	formModel, cmd := m.form.Update(msg)
	if f, ok := formModel.(*huh.Form); ok {
		m.form = f
		if m.form.State == huh.StateCompleted && m.state != stateDownloading {
			m.state = stateDownloading
			return m, m.downloadBottle()
		} else if m.form.State == huh.StateCompleted && m.state == stateDone {
			return m, tea.Quit
		}
		cmds = append(cmds, cmd)
	}

	// Keep the progress bar ticking
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := m.styles
	switch m.state {
	case stateDownloading:
		header := m.appBoundaryView("🍺 Bottle Downloader")
		progressView := m.lg.NewStyle().Margin(1, 1, 0, 4).Render(m.progress.View())
		footer := m.appBoundaryView("Downloading... Press 'q' to quit")
		// return s.Base.Render(form + "\n\n" + progressView + "\n\n" + footer)
		return s.Base.Render(header + "\n" + progressView + "\n\n" + footer)

	case stateQuitting:
		// title := s.Highlight.Render("Bottle Downloader")
		// var b strings.Builder
		// fmt.Fprintf(&b, "Congratulations, you’re Charm’s newest\n%s!\n\n", title)
		// return s.Status.Margin(0, 1).Padding(1, 2).Width(48).Render(b.String()) + "\n\n"
		// return s.Base.Render("🍺 Bottle dud? That's cool.")
		return lipgloss.NewStyle().Margin(1, 0, 2, 4).Render("🍺 Bottle dud? That's cool.")

	case stateDone:
		if m.err != nil {
			return s.Base.Render(
				m.appErrorBoundaryView("Error downloading: " + m.err.Error()),
			)
		}
		header := m.appBoundaryView("🍾 Download Complete! 💥")
		progressView := m.lg.NewStyle().Margin(1, 1, 0, 4).Render(m.progress.View())
		footer := m.appBoundaryView("Downloading... Press 'q' to quit")
		// return s.Base.Render(form + "\n\n" + progressView + "\n\n" + footer)
		return s.Base.Render(header + "\n" + progressView + "\n\n" + footer)

	default:
		v := strings.TrimSuffix(m.form.View(), "\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		var status string
		{
			var deps string
			if len(m.formula.Dependencies) > 0 {
				deps = "\n\n" + s.StatusHeader.Render("Dependencies") + "\n"
				for _, dep := range m.formula.Dependencies {
					deps += "  • " + dep + "\n"
				}
			}
			const statusWidth = 60
			statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
			status = s.Status.
				Height(lipgloss.Height(form)).
				Width(statusWidth).
				MarginLeft(statusMarginLeft).
				Render(s.StatusHeader.Render(m.formula.Name)+"\n"+
					"Version: "+m.formula.Versions.Stable+"\n"+
					"Homepage: "+m.formula.Homepage+"\n"+
					"Description: "+m.formula.Desc,
					deps,
				)
		}
		errors := m.errorView()
		header := m.appBoundaryView("🍺 Bottle Downloader")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func (m *Model) downloadBottle() tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("GET", m.selectedURL, nil)
		if err != nil {
			return m.appErrorBoundaryView(err.Error())
		}
		req.Header.Add("Authorization", "Bearer QQ==")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return m.appErrorBoundaryView(err.Error())
		}
		defer resp.Body.Close()

		f, err := os.Create(m.formula.Name + ".tar.gz")
		if err != nil {
			return m.appErrorBoundaryView(err.Error())
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

		m.created = true
		m.state = stateDone

		return progressMsg(100.0)
	}
}
