package registry

// Configuration to control the registry
type Configuration struct {
	InitOnNew      bool // indicates whether registry should be initialized from API on creation
	SkipInitErrors bool // indicates whether to ignore particular resource or object initialization errors
	// TODO add `Throttle` field to decrease load on OMAs API
}

// DefaultConfiguration creates a Configuration with the default settings
func DefaultConfiguration() *Configuration {
	return &Configuration{
		InitOnNew:      true,
		SkipInitErrors: true, // true by default, because OMA API returns many objects without ObjectLink filled in,
		// which make it impossible to initialize Object with Resources
	}
}
