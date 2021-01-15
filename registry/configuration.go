package registry

// Configuration to control the registry.
type Configuration struct {
	InitOnNew      bool     // indicates whether registry should be initialized from API on creation
	SkipInitErrors bool     // indicates whether to ignore particular resource or object initialization errors
	Sanitize       bool     // indicates whether bjects and Resources description should be cleaned up on registry initialization
	Sanitizer      []string // strings that should be removed from resource and object description
	// TODO add `Throttle` field to decrease load on OMAs API
}

// DefaultConfiguration creates a Configuration with the default settings.
func DefaultConfiguration() *Configuration {
	return &Configuration{
		InitOnNew:      true,
		SkipInitErrors: true, // true by default, because OMA API returns some objects without ObjectLink filled in,
		// which makes it impossible to initialize objects and resources with strict validation
		Sanitize:  true,
		Sanitizer: DefaultSanitizer(),
	}
}

// DefaultSanitizer returns an array with default sanitizer strings.
func DefaultSanitizer() []string {
	return []string{
		"\n", "\\\n", "\t", "\\t", "\\", "•", "  ",
	}
}
