package registry

// Lwm2M struct for Lwm2M wrapper entry.
type Lwm2M struct {
	Object *Object `xml:"Object"`
}

// ObjectMeta struct for ObjectMeta.
type ObjectMeta struct {
	ObjectID          int64  `xml:"ObjectID"`
	Ver               string `xml:"Ver"`
	URN               string `xml:"URN"`
	Name              string `xml:"Name"`
	Description       string `xml:"Description"`
	Owner             string `xml:"Owner"`
	Label             string `xml:"Label"`
	ObjectLink        string `xml:"ObjectLink"`
	ObjectLinkVisible string `xml:"ObjectLinkVisible"`
	SpecLink          string `xml:"SpecLink"`
	SpecLinkVisible   string `xml:"SpecLinkVisible"`
	VortoLink         string `xml:"VortoLink"`
}

// Object structure represents OMA Object entity.
type Object struct {
	Name              string          `xml:"Name" yaml:"Name"`
	Description1      string          `xml:"Description1" yaml:"Description1"`
	Description2      string          `xml:"Description2" yaml:"Description2,omitempty"`
	ObjectID          int64           `xml:"ObjectID" yaml:"ObjectID"`
	ObjectURN         string          `xml:"ObjectURN" yaml:"ObjectURN"`
	LWM2MVersion      string          `xml:"LWM2MVersion" yaml:"LWM2MVersion,omitempty"`
	ObjectVersion     string          `xml:"ObjectVersion" yaml:"ObjectVersion"`
	MultipleInstances InstanceType    `xml:"MultipleInstances" yaml:"MultipleInstances"`
	Mandatory         MandatoryType   `xml:"Mandatory" yaml:"Mandatory"`
	Resources         ObjectResources `xml:"Resources" yaml:"Resources,inline"`
}

// ObjectResources struct for ObjectResources field.
type ObjectResources struct {
	Item []Resource `xml:"Item" yaml:"Resources"`
}

// ObjectComparison contains details of two objects comparison.
type ObjectComparison struct {
	Difference DifferenceType
	Object     *Object // object in existing registry
	ObjectComp *Object // object in compared registry, which is passed as parameter to Compare function
}
