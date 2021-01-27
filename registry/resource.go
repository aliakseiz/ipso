package registry

// Resource structure represents OMA Resource entity
type Resource struct {
	ID                int64         `xml:"ID,attr" yaml:"ID"`
	Name              string        `xml:"Name" yaml:"Name"`
	Operations        OperationType `xml:"Operations" yaml:"Operations"`
	MultipleInstances InstanceType  `xml:"MultipleInstances" yaml:"MultipleInstances"`
	Mandatory         MandatoryType `xml:"Mandatory" yaml:"Mandatory"`
	Type              ResourceType  `xml:"Type" yaml:"Type"`
	RangeEnumeration  string        `xml:"RangeEnumeration" yaml:"RangeEnumeration,omitempty"`
	Units             string        `xml:"Units" yaml:"Units,omitempty"`
	Description       string        `xml:"Description" yaml:"Description,omitempty"`
}
