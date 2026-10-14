package ipso

// Resource is an OMA resource: a single readable, writable or executable value of an Object, such as
// resource 5700 Sensor Value of object 3303. The ID is unique within one object version.
// Type is empty for resources without data, such as executable ones. RangeEnumeration is free text describing
// the valid values, such as "0-100". Units is free text, usually a unit symbol such as "s" or "W", empty when the value has none.
// The YAML keys are the lowercase field names, the format of the shipped registry files.
type Resource struct {
	ID                int32         `json:"ID"                         xml:"ID,attr"           yaml:"id"`
	Name              string        `json:"Name"                       xml:"Name"              yaml:"name"`
	Operations        OperationType `json:"Operations"                 xml:"Operations"        yaml:"operations"`
	MultipleInstances InstanceType  `json:"MultipleInstances"          xml:"MultipleInstances" yaml:"multipleinstances"`
	Mandatory         MandatoryType `json:"Mandatory"                  xml:"Mandatory"         yaml:"mandatory"`
	Type              ResourceType  `json:"Type"                       xml:"Type"              yaml:"type"`
	RangeEnumeration  string        `json:"RangeEnumeration,omitempty" xml:"RangeEnumeration"  yaml:"rangeenumeration"`
	Units             string        `json:"Units,omitempty"            xml:"Units"             yaml:"units"`
	Description       string        `json:"Description,omitempty"      xml:"Description"       yaml:"description"`
}
