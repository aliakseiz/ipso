package registry

const (
	SOURCE_URL  = "URL"
	SOURCE_FILE = "File"

	DefaultSource = SOURCE_URL
	DefaultPath   = ""
)

// Configuration to control the registry
type Configuration struct {
	Source string
	Path   string
}

// DefaultConfiguration creates a Configuration
// with the default settings
func DefaultConfiguration() *Configuration {
	return &Configuration{
		// TODO fill in
		Source: DefaultSource,
		Path:   DefaultPath,
	}
}
