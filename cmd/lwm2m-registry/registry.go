package main

type ObjectMeta struct {
	ObjectID          int64  `json:"ObjectID"`
	Ver               string `json:"Ver"`
	URN               string `json:"URN"`
	Name              string `json:"Name"`
	Description       string `json:"Description"`
	Owner             string `json:"Owner"`
	Label             string `json:"Label"`
	ObjectLink        string `json:"ObjectLink"`
	ObjectLinkVisible string `json:"ObjectLinkVisible"`
	SpecLink          string `json:"SpecLink"`
	SpecLinkVisible   string `json:"SpecLinkVisible"`
	VortoLink         string `json:"VortoLink"`
}

type LwM2M struct {
	Object Object
}

type Object struct {
	Name             string
	Description1     string
	Description2     string
	ObjectID         int64
	ObjectURN        string
	MultipleInstance string
	Mandatory        string
	Resources        []Resource
}

type Resource struct {
	Name             string
	Operations       string
	MultipleInstance string
	Mandatory        string
	Type             string
	RangeEnumeration string
	Units            string
	Description      string
}

// Registry an API to control the registry
// type Registry interface {
//
// }

// RegistryOptions to configure the registry
type RegistryOptions struct {
}

// Registry holds registry settings
type Registry struct {
	Objects []Object
}

// New creates a new registry, using custom configuration
func New(cfg *RegistryOptions) *Registry {
	if cfg == nil {
		cfg = DefaultRegistryConfig()
	}

	return nil // TODO initialize registry
}

// DefaultRegistryConfig creates a RegistryOptions
// with the default settings
func DefaultRegistryConfig() *RegistryOptions {
	return &RegistryOptions{
		// TODO fill in
	}
}
