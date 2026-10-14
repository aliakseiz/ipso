package ipso

// InstanceType tells whether an object or a resource can have several instances.
type InstanceType string

// InstanceType values.
const (
	InstanceTypeSingle   InstanceType = "Single"
	InstanceTypeMultiple InstanceType = "Multiple"
)

// String returns the value as a string.
func (e InstanceType) String() string {
	return string(e)
}

// MandatoryType tells whether an object or a resource must be implemented by a device.
type MandatoryType string

// MandatoryType values.
const (
	MandatoryTypeOptional  MandatoryType = "Optional"
	MandatoryTypeMandatory MandatoryType = "Mandatory"
)

// String returns the value as a string.
func (e MandatoryType) String() string {
	return string(e)
}

// OperationType is the set of operations a resource supports.
type OperationType string

// OperationType values.
const (
	OperationTypeRead      OperationType = "R"
	OperationTypeWrite     OperationType = "W"
	OperationTypeReadWrite OperationType = "RW"
	OperationTypeExecute   OperationType = "E"
	// OperationTypeEmpty marks resources listed without operations; 106 resources in the shipped registry have it.
	OperationTypeEmpty OperationType = ""
)

// String returns the value as a string.
func (e OperationType) String() string {
	return string(e)
}

// ResourceType is the data type of a resource.
type ResourceType string

// ResourceType values according to OMA-TS-LightweightM2M_Core-V1_2-20201110-A Appendix C. Data types (Normative).
const (
	ResourceTypeString          ResourceType = "String"
	ResourceTypeInteger         ResourceType = "Integer"
	ResourceTypeUnsignedInteger ResourceType = "Unsigned Integer"
	ResourceTypeFloat           ResourceType = "Float"
	ResourceTypeBoolean         ResourceType = "Boolean"
	ResourceTypeOpaque          ResourceType = "Opaque"
	ResourceTypeTime            ResourceType = "Time"
	ResourceTypeObjLink         ResourceType = "Objlnk"
	ResourceTypeCoreLink        ResourceType = "Corelnk"
	// ResourceTypeEmpty is the type of resources without data, such as executable ones.
	// The shipped registry files store it as an empty string.
	ResourceTypeEmpty ResourceType = ""
	// ResourceTypeNone is a textual alias for resources without data.
	// No resource in the shipped registry files has this type, resources without data have ResourceTypeEmpty.
	ResourceTypeNone ResourceType = "none"
)

// String returns the value as a string.
func (e ResourceType) String() string {
	return string(e)
}

// DifferenceType explains why Reg.Compare lists an object.
type DifferenceType string

// DifferenceType values.
const (
	// DifferenceTypeUnknown is never produced by Reg.Compare.
	DifferenceTypeUnknown DifferenceType = "Unknown reason"
	// DifferenceTypeNewObject marks an object present only in the compared registry.
	DifferenceTypeNewObject DifferenceType = "New object added"
	// DifferenceTypeObjectRemoved marks an object present only in the registry Compare is called on.
	DifferenceTypeObjectRemoved DifferenceType = "Object was removed"
)

// String returns the value as a string.
func (e DifferenceType) String() string {
	return string(e)
}
