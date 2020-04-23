package registry

import openapi "github.com/aliakseiz/ipso-registry/api/client"

// Resource structure represents OMA Resource entity
type Resource struct {
	ID               int32         `yaml:"ID"`
	Name             string        `yaml:"Name"`
	Operations       OperationType `yaml:"Operations"`
	MultipleInstance InstanceType  `yaml:"MultipleInstances"`
	Mandatory        MandatoryType `yaml:"Mandatory"`
	Type             ResourceType  `yaml:"Type"`
	RangeEnumeration string        `yaml:"RangeEnumeration,omitempty"`
	Units            string        `yaml:"Units,omitempty"`
	Description      string        `yaml:"Description,omitempty"`
}

func mapResources(omaResources []openapi.Resource) []Resource {
	resources := make([]Resource, len(omaResources))

	for i, omaRes := range omaResources {
		resources[i] = mapResource(omaRes)
	}

	return resources
}

func mapResource(omaResource openapi.Resource) Resource {
	return Resource{
		ID:               omaResource.ID,
		Name:             omaResource.Name,
		Operations:       OperationType(omaResource.Operations),
		MultipleInstance: InstanceType(omaResource.MultipleInstances),
		Mandatory:        MandatoryType(omaResource.Mandatory),
		Type:             ResourceType(omaResource.Type),
		RangeEnumeration: omaResource.RangeEnumeration,
		Units:            omaResource.Units,
		Description:      omaResource.Description,
	}
}
