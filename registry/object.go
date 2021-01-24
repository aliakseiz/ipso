package registry

// Lwm2M struct for Lwm2M wrapper entry.
type Lwm2M struct {
	Object *Object `json:"Object,omitempty" xml:"Object"`
}

// ObjectMeta struct for ObjectMeta.
type ObjectMeta struct {
	ObjectID          int64  `json:"ObjectID" xml:"ObjectID"`
	Ver               string `json:"Ver,omitempty" xml:"Ver"`
	URN               string `json:"URN,omitempty" xml:"URN"`
	Name              string `json:"Name" xml:"Name"`
	Description       string `json:"Description,omitempty" xml:"Description"`
	Owner             string `json:"Owner,omitempty" xml:"Owner"`
	Label             string `json:"Label,omitempty" xml:"Label"`
	ObjectLink        string `json:"ObjectLink" xml:"ObjectLink"`
	ObjectLinkVisible string `json:"ObjectLinkVisible,omitempty" xml:"ObjectLinkVisible"`
	SpecLink          string `json:"SpecLink,omitempty" xml:"SpecLink"`
	SpecLinkVisible   string `json:"SpecLinkVisible,omitempty" xml:"SpecLinkVisible"`
	VortoLink         string `json:"VortoLink,omitempty" xml:"VortoLink"`
}

// Object structure represents OMA Object entity.
type Object struct {
	Name              string          `json:"Name" xml:"Name"`
	Description1      string          `json:"Description1,omitempty" xml:"Description1"`
	Description2      string          `json:"Description2,omitempty" xml:"Description2"`
	ObjectID          int64           `json:"ObjectID" xml:"ObjectID"`
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
	Object     *Object // object in existing registry
	ObjectComp *Object // object in compared registry, which is passed as parameter to Compare function
}
