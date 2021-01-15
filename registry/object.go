package registry

// Lwm2M struct for Lwm2M wrapper entry.
type Lwm2M struct {
	Object Object `json:"Object,omitempty" xml:"Object"`
}

// Object structure represents OMA Object entity.
type Object struct {
	Name              string          `json:"Name" xml:"Name"`
	Description1      string          `json:"Description1,omitempty" xml:"Description1"`
	Description2      string          `json:"Description2,omitempty" xml:"Description2"`
	ObjectID          int32           `json:"ObjectID" xml:"ObjectID"`
	ObjectURN         string          `json:"ObjectURN,omitempty" xml:"ObjectURN"`
	LWM2MVersion      string          `json:"LWM2MVersion,omitempty" xml:"LWM2MVersion"`
	ObjectVersion     string          `json:"ObjectVersion,omitempty" xml:"ObjectVersion"`
	MultipleInstances InstanceType    `json:"MultipleInstances" xml:"MultipleInstances"`
	Mandatory         MandatoryType   `json:"Mandatory" xml:"Mandatory"`
	Resources         ObjectResources `json:"Resources" xml:"Resources"`
}

// ObjectResources struct for ObjectResources field.
type ObjectResources struct {
	Item []Resource `json:"Item,omitempty" xml:"Item"`
}

// ObjectComparison contains details of two objects comparison.
type ObjectComparison struct {
	Difference DifferenceType
	Object     Object // object in existing registry
	ObjectComp Object // object in compared registry, which is passed as parameter to Compare function
}
