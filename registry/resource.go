package registry

import openapi "github.com/aliakseiz/lwm2m-registry/api/client"

type Resource struct {
	Name             string
	Operations       OperationType
	MultipleInstance InstanceType
	Mandatory        MandatoryType
	Type             ResourceType
	RangeEnumeration string
	Units            string
	Description      string
}

func mapResources(omaResources []openapi.Resource) []*Resource {
	var resources []*Resource

	for _, omaRes := range omaResources {
		if res := mapResource(omaRes); res != nil {
			resources = append(resources, res)
		}
		// TODO log mapping error?
	}

	return resources
}

func mapResource(omaResource openapi.Resource) *Resource {
	// TODO validate omaResource, return nil or error
	return &Resource{
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
