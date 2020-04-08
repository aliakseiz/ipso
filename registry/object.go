package registry

import openapi "github.com/aliakseiz/lwm2m-registry/api/client"

// Object structure represents OMA Object entity
type Object struct {
	Name             string        `yaml:"Name"`
	Description1     string        `yaml:"Description1,omitempty"`
	Description2     string        `yaml:"Description2,omitempty"`
	ObjectID         int32         `yaml:"ObjectID"`
	ObjectURN        string        `yaml:"ObjectURN,omitempty"`
	LwM2MVersion     string        `yaml:"LWM2MVersion,omitempty"`
	ObjectVersion    string        `yaml:"ObjectVersion,omitempty"`
	MultipleInstance InstanceType  `yaml:"MultipleInstances"`
	Mandatory        MandatoryType `yaml:"Mandatory"`
	Resources        []Resource    `yaml:"Resources"`
}

// ObjectComparison contains details of two objects comparison
type ObjectComparison struct {
	Difference DifferenceType
	Object     *Object // object in existing registry
	ObjectComp *Object // object in compared registry, which is passed as parameter to Compare function
}

func mapObject(omaObject *openapi.Object) Object {
	return Object{
		Name:             omaObject.Name,
		Description1:     omaObject.Description1,
		Description2:     omaObject.Description2,
		ObjectID:         omaObject.ObjectID,
		ObjectURN:        omaObject.ObjectURN,
		LwM2MVersion:     omaObject.LWM2MVersion,
		ObjectVersion:    omaObject.ObjectVersion,
		MultipleInstance: InstanceType(omaObject.MultipleInstances),
		Mandatory:        MandatoryType(omaObject.Mandatory),
		Resources:        mapResources(omaObject.Resources.Item),
	}
}
