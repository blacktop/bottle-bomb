package cmd

type Formula struct {
	Name              string        `json:"name"`
	FullName          string        `json:"full_name"`
	Tap               string        `json:"tap"`
	Oldnames          []interface{} `json:"oldnames"`
	Aliases           []interface{} `json:"aliases"`
	VersionedFormulae []interface{} `json:"versioned_formulae"`
	Desc              string        `json:"desc"`
	License           string        `json:"license"`
	Homepage          string        `json:"homepage"`
	Versions          struct {
		Stable string `json:"stable"`
		Head   string `json:"head"`
		Bottle bool   `json:"bottle"`
	} `json:"versions"`
	Urls struct {
		Stable struct {
			URL      string      `json:"url"`
			Tag      interface{} `json:"tag"`
			Revision interface{} `json:"revision"`
			Using    interface{} `json:"using"`
			Checksum string      `json:"checksum"`
		} `json:"stable"`
		Head struct {
			URL    string      `json:"url"`
			Branch string      `json:"branch"`
			Using  interface{} `json:"using"`
		} `json:"head"`
	} `json:"urls"`
	Revision      int `json:"revision"`
	VersionScheme int `json:"version_scheme"`
	Bottle        struct {
		Stable struct {
			Rebuild int    `json:"rebuild"`
			RootURL string `json:"root_url"`
			Files   struct {
				Arm64Sonoma struct {
					Cellar string `json:"cellar"`
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"arm64_sonoma"`
				Arm64Ventura struct {
					Cellar string `json:"cellar"`
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"arm64_ventura"`
				Arm64Monterey struct {
					Cellar string `json:"cellar"`
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"arm64_monterey"`
				Sonoma struct {
					Cellar string `json:"cellar"`
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"sonoma"`
				Ventura struct {
					Cellar string `json:"cellar"`
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"ventura"`
				Monterey struct {
					Cellar string `json:"cellar"`
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"monterey"`
				Arm64Linux struct {
					Cellar string `json:"cellar"`
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"arm64_linux"`
				X8664Linux struct {
					Cellar string `json:"cellar"`
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"x86_64_linux"`
			} `json:"files"`
		} `json:"stable"`
	} `json:"bottle"`
	PourBottleOnlyIf        interface{}   `json:"pour_bottle_only_if"`
	KegOnly                 bool          `json:"keg_only"`
	KegOnlyReason           interface{}   `json:"keg_only_reason"`
	Options                 []interface{} `json:"options"`
	BuildDependencies       []string      `json:"build_dependencies"`
	Dependencies            []string      `json:"dependencies"`
	TestDependencies        []interface{} `json:"test_dependencies"`
	RecommendedDependencies []interface{} `json:"recommended_dependencies"`
	OptionalDependencies    []interface{} `json:"optional_dependencies"`
	UsesFromMacos           []interface{} `json:"uses_from_macos"`
	UsesFromMacosBounds     []interface{} `json:"uses_from_macos_bounds"`
	Requirements            []interface{} `json:"requirements"`
	ConflictsWith           []interface{} `json:"conflicts_with"`
	ConflictsWithReasons    []interface{} `json:"conflicts_with_reasons"`
	LinkOverwrite           []interface{} `json:"link_overwrite"`
	Caveats                 interface{}   `json:"caveats"`
	Installed               []interface{} `json:"installed"`
	LinkedKeg               interface{}   `json:"linked_keg"`
	Pinned                  bool          `json:"pinned"`
	Outdated                bool          `json:"outdated"`
	Deprecated              bool          `json:"deprecated"`
	DeprecationDate         interface{}   `json:"deprecation_date"`
	DeprecationReason       interface{}   `json:"deprecation_reason"`
	Disabled                bool          `json:"disabled"`
	DisableDate             interface{}   `json:"disable_date"`
	DisableReason           interface{}   `json:"disable_reason"`
	PostInstallDefined      bool          `json:"post_install_defined"`
	Service                 interface{}   `json:"service"`
	TapGitHead              string        `json:"tap_git_head"`
	RubySourcePath          string        `json:"ruby_source_path"`
	RubySourceChecksum      struct {
		Sha256 string `json:"sha256"`
	} `json:"ruby_source_checksum"`
	Variations struct {
	} `json:"variations"`
	Analytics struct {
		Install          map[string]any `json:"install"`
		InstallOnRequest map[string]any `json:"install_on_request"`
		BuildError       map[string]any `json:"build_error"`
	} `json:"analytics"`
	GeneratedDate string `json:"generated_date"`
}

type Bottle struct {
	SchemaVersion int `json:"schemaVersion"`
	Manifests     []struct {
		MediaType string `json:"mediaType"`
		Digest    string `json:"digest"`
		Size      int    `json:"size"`
		Platform  struct {
			Architecture string `json:"architecture"`
			Os           string `json:"os"`
			OsVersion    string `json:"os.version"`
		} `json:"platform"`
		Annotations struct {
			OrgOpencontainersImageRefName string `json:"org.opencontainers.image.ref.name"`
			ShBrewBottleCPUVariant        string `json:"sh.brew.bottle.cpu.variant"`
			ShBrewBottleDigest            string `json:"sh.brew.bottle.digest"`
			ShBrewBottleGlibcVersion      string `json:"sh.brew.bottle.glibc.version"`
			ShBrewTab                     string `json:"sh.brew.tab"`
		} `json:"annotations"`
	} `json:"manifests"`
	Annotations struct {
		ComGithubPackageType                string `json:"com.github.package.type"`
		OrgOpencontainersImageCreated       string `json:"org.opencontainers.image.created"`
		OrgOpencontainersImageDescription   string `json:"org.opencontainers.image.description"`
		OrgOpencontainersImageDocumentation string `json:"org.opencontainers.image.documentation"`
		OrgOpencontainersImageLicense       string `json:"org.opencontainers.image.license"`
		OrgOpencontainersImageRefName       string `json:"org.opencontainers.image.ref.name"`
		OrgOpencontainersImageRevision      string `json:"org.opencontainers.image.revision"`
		OrgOpencontainersImageSource        string `json:"org.opencontainers.image.source"`
		OrgOpencontainersImageTitle         string `json:"org.opencontainers.image.title"`
		OrgOpencontainersImageURL           string `json:"org.opencontainers.image.url"`
		OrgOpencontainersImageVendor        string `json:"org.opencontainers.image.vendor"`
		OrgOpencontainersImageVersion       string `json:"org.opencontainers.image.version"`
	} `json:"annotations"`
}
