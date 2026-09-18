package manager

type Manager interface {
	Name() string
	IsAvailable() bool
	Install(pkgName string) error
	Remove(pkgName string) error
	Search(pkgName string) ([]SearchResult, error)
	Update() error
}

type SearchResult struct {
	Name        string
	Description string
	Version     string
	Source      string
}
