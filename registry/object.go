package registry

const (
	DefaultLwM2MVersion  = "1.0"
	DefaultObjectVersion = "1.0"
)

// Lwm2M struct for Lwm2M wrapper entry.
type Lwm2M struct {
	Object Object `json:"Object,omitempty" xml:"Object"`
}

// Object structure represents OMA Object entity.
type Object struct {
	Name              string          `json:"Name" xml:"Name" yaml:"Name"`
	Description1      string          `json:"Description1,omitempty" xml:"Description1" yaml:"Description1"`
	Description2      string          `json:"Description2,omitempty" xml:"Description2" yaml:"Description2"`
	ObjectID          int32           `json:"ObjectID" xml:"ObjectID" yaml:"ObjectID"`
	ObjectURN         string          `json:"ObjectURN,omitempty" xml:"ObjectURN" yaml:"ObjectURN"`
	LWM2MVersion      string          `json:"LWM2MVersion,omitempty" xml:"LWM2MVersion" yaml:"LWM2MVersion"`
	ObjectVersion     string          `json:"ObjectVersion,omitempty" xml:"ObjectVersion" yaml:"ObjectVersion"`
	MultipleInstances InstanceType    `json:"MultipleInstances" xml:"MultipleInstances" yaml:"MultipleInstances"`
	Mandatory         MandatoryType   `json:"Mandatory" xml:"Mandatory" yaml:"Mandatory"`
	Resources         ObjectResources `json:"Resources" xml:"Resources" yaml:"Resources"`
}

// ObjectResources struct for ObjectResources field.
type ObjectResources struct {
	Item []Resource `json:"Item,omitempty" xml:"Item" yaml:"Item,omitempty"`
}

// ObjectMeta struct for ObjectMeta.
type ObjectMeta struct {
	ObjectID          int32  `json:"ObjectID" xml:"ObjectID" yaml:"ObjectID"`
	Ver               string `json:"Ver,omitempty" xml:"Ver" yaml:"Ver"`
	URN               string `json:"URN,omitempty" xml:"URN" yaml:"URN"`
	Name              string `json:"Name" xml:"Name" yaml:"Name"`
	Description       string `json:"Description,omitempty" xml:"Description" yaml:"Description"`
	Owner             string `json:"Owner,omitempty" xml:"Owner" yaml:"Owner"`
	Label             string `json:"Label,omitempty" xml:"Label" yaml:"Label"`
	ObjectLink        string `json:"ObjectLink" xml:"ObjectLink" yaml:"ObjectLink"`
	ObjectLinkVisible string `json:"ObjectLinkVisible,omitempty" xml:"ObjectLinkVisible" yaml:"ObjectLinkVisible"`
	SpecLink          string `json:"SpecLink,omitempty" xml:"SpecLink" yaml:"SpecLink"`
	SpecLinkVisible   string `json:"SpecLinkVisible,omitempty" xml:"SpecLinkVisible" yaml:"SpecLinkVisible"`
	VortoLink         string `json:"VortoLink,omitempty" xml:"VortoLink" yaml:"VortoLink"`
}

// ObjectComparison contains details of two objects comparison.
type ObjectComparison struct {
	Difference DifferenceType
	Object     Object // object in existing registry
	ObjectComp Object // object in compared registry, which is passed as parameter to Compare function
}
