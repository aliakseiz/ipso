package registry

import openapi "github.com/aliakseiz/lwm2m-registry/api/client"

type Object struct {
	Name             string
	Description1     string
	Description2     string
	ObjectID         int32
	ObjectURN        string
	MultipleInstance InstanceType
	Mandatory        MandatoryType
	Resources        []*Resource
}

func mapObject(omaObject *openapi.Object) *Object {
	return &Object{
		Name:             omaObject.Name,
		Description1:     omaObject.Description1,
		Description2:     omaObject.Description2,
		ObjectID:         omaObject.ObjectID,
		ObjectURN:        omaObject.ObjectURN,
		MultipleInstance: InstanceType(omaObject.MultipleInstances),
		Mandatory:        MandatoryType(omaObject.Mandatory),
		Resources:        mapResources(omaObject.Resources.Item),
	}
}
