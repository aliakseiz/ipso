package ipso_registry

// InstanceType constants.
const (
	InstanceTypeSingle   InstanceType = "Single"
	InstanceTypeMultiple InstanceType = "Multiple"
)

// InstanceType represents the type of data.
type InstanceType string

// String for InstanceType.
func (e InstanceType) String() string {
	return string(e)
}

// MandatoryType constants.
const (
	MandatoryTypeOptional  MandatoryType = "Optional"
	MandatoryTypeMandatory MandatoryType = "Mandatory"
)

// MandatoryType represents the type of data.
type MandatoryType string

// String for MandatoryType.
func (e MandatoryType) String() string {
	return string(e)
}

// OperationType constants.
const (
	OperationTypeRead      OperationType = "R"
	OperationTypeWrite     OperationType = "W"
	OperationTypeReadWrite OperationType = "RW"
	OperationTypeExecute   OperationType = "E"
)

// OperationType represents the type of data.
type OperationType string

// String for OperationType.
func (e OperationType) String() string {
	return string(e)
}

// ResourceType constants according to OMA-TS-LightweightM2M_Core-V1_2-20201110-A Appendix C. Data types (Normative).
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
	ResourceTypeNone            ResourceType = "none"
)

// ResourceType represents the type of data.
type ResourceType string

// String for ResourceType.
func (e ResourceType) String() string {
	return string(e)
}

// DifferenceType constants.
const (
	DifferenceTypeUnknown       DifferenceType = "Unknown reason"
	DifferenceTypeNewObject     DifferenceType = "New object added"
	DifferenceTypeObjectRemoved DifferenceType = "Object was removed"
)

// DifferenceType represents the objects difference explanation.
type DifferenceType string

// String for ResourceType.
func (e DifferenceType) String() string {
	return string(e)
}
