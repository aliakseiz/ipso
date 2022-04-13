package ipso

// Resource structure represents OMA Resource entity.
type Resource struct {
	ID                int32         `json:"ID" xml:"ID,attr"`
	Name              string        `json:"Name" xml:"Name"`
	Operations        OperationType `json:"Operations" xml:"Operations"`
	MultipleInstances InstanceType  `json:"MultipleInstances" xml:"MultipleInstances"`
	Mandatory         MandatoryType `json:"Mandatory" xml:"Mandatory"`
	Type              ResourceType  `json:"Type" xml:"Type"`
	RangeEnumeration  string        `json:"RangeEnumeration,omitempty" xml:"RangeEnumeration"`
	Units             string        `json:"Units,omitempty" xml:"Units"`
	Description       string        `json:"Description,omitempty" xml:"Description"`
}
