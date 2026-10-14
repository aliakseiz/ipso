package ipso

import "slices"

const (
	// DefaultLwM2MVersion is the LwM2M version set on objects that the OMA API lists without one.
	DefaultLwM2MVersion = "1.0"
	// DefaultObjectVersion is the object version set on objects that the OMA API lists without one.
	DefaultObjectVersion = "1.0"
)

// Lwm2M is the root element of an OMA object definition XML, the file that ObjectMeta.ObjectLink points to.
type Lwm2M struct {
	Object Object `json:"Object" xml:"Object"`
}

// Object is an OMA object: a versioned definition of a device capability, such as 3303 Temperature,
// with the resources it consists of. An object is identified by ObjectID and ObjectVersion, the definition
// version such as "1.0" or "1.2"; LWM2MVersion is the LwM2M protocol version the object is defined for.
// The YAML keys are the field names, unlike the lowercase keys of Resource.
type Object struct {
	Name              string          `json:"Name"                    xml:"Name"              yaml:"Name"`
	Description1      string          `json:"Description1,omitempty"  xml:"Description1"      yaml:"Description1"`
	Description2      string          `json:"Description2,omitempty"  xml:"Description2"      yaml:"Description2"`
	ObjectID          int32           `json:"ObjectID"                xml:"ObjectID"          yaml:"ObjectID"`
	ObjectURN         string          `json:"ObjectURN,omitempty"     xml:"ObjectURN"         yaml:"ObjectURN"`
	LWM2MVersion      string          `json:"LWM2MVersion,omitempty"  xml:"LWM2MVersion"      yaml:"LWM2MVersion"`
	ObjectVersion     string          `json:"ObjectVersion,omitempty" xml:"ObjectVersion"     yaml:"ObjectVersion"`
	MultipleInstances InstanceType    `json:"MultipleInstances"       xml:"MultipleInstances" yaml:"MultipleInstances"`
	Mandatory         MandatoryType   `json:"Mandatory"               xml:"Mandatory"         yaml:"Mandatory"`
	Resources         ObjectResources `json:"Resources"               xml:"Resources"         yaml:"Resources"`
}

// ObjectResources wraps the list of resources of an Object, matching the Resources/Item structure of the OMA XML.
type ObjectResources struct {
	Item []Resource `json:"Item,omitempty" xml:"Item" yaml:"Item,omitempty"`
}

// ObjectMeta is an entry of the objects list served by the OMA API; ObjectLink is the URL of the object definition.
type ObjectMeta struct {
	ObjectID          int32  `json:"ObjectID"                    xml:"ObjectID"          yaml:"ObjectID"`
	Ver               string `json:"Ver,omitempty"               xml:"Ver"               yaml:"Ver"`
	URN               string `json:"URN,omitempty"               xml:"URN"               yaml:"URN"`
	Name              string `json:"Name"                        xml:"Name"              yaml:"Name"`
	Description       string `json:"Description,omitempty"       xml:"Description"       yaml:"Description"`
	Owner             string `json:"Owner,omitempty"             xml:"Owner"             yaml:"Owner"`
	Label             string `json:"Label,omitempty"             xml:"Label"             yaml:"Label"`
	ObjectLink        string `json:"ObjectLink"                  xml:"ObjectLink"        yaml:"ObjectLink"`
	ObjectLinkVisible string `json:"ObjectLinkVisible,omitempty" xml:"ObjectLinkVisible" yaml:"ObjectLinkVisible"`
	SpecLink          string `json:"SpecLink,omitempty"          xml:"SpecLink"          yaml:"SpecLink"`
	SpecLinkVisible   string `json:"SpecLinkVisible,omitempty"   xml:"SpecLinkVisible"   yaml:"SpecLinkVisible"`
	VortoLink         string `json:"VortoLink,omitempty"         xml:"VortoLink"         yaml:"VortoLink"`
}

// ObjectComparison is one difference found by Reg.Compare. The object missing on one side is the zero Object.
type ObjectComparison struct {
	Difference DifferenceType
	Object     Object // object in the registry Compare is called on
	ObjectComp Object // object in the registry passed to Compare
}

// clone returns a copy of the object that shares no memory with the original.
func (o *Object) clone() Object {
	dup := *o
	dup.Resources.Item = slices.Clone(o.Resources.Item)

	return dup
}

// resource returns the resource with the ID, or ErrResourceNotFound.
func (o *Object) resource(id int32) (Resource, error) {
	for _, res := range o.Resources.Item {
		if res.ID == id {
			return res, nil
		}
	}

	return Resource{}, ErrResourceNotFound
}
