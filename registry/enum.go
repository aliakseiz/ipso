package registry

// InstanceType constants
const (
	InstanceTypeSingle   InstanceType = "Single"
	InstanceTypeMultiple InstanceType = "Multiple"
)

// InstanceType represents the type of data
type InstanceType string

// String for InstanceType
func (e InstanceType) String() string {
	return string(e)
}

// MandatoryType constants
const (
	MandatoryTypeOptional  MandatoryType = "Optional"
	MandatoryTypeMandatory MandatoryType = "Mandatory"
)

// MandatoryType represents the type of data
type MandatoryType string

// String for MandatoryType
func (e MandatoryType) String() string {
	return string(e)
}

// OperationType constants
const (
	OperationTypeRead      OperationType = "R"
	OperationTypeWrite     OperationType = "W"
	OperationTypeReadWrite OperationType = "RW"
	OperationTypeExecute   OperationType = "E"
)

// OperationType represents the type of data
type OperationType string

// String for OperationType
func (e OperationType) String() string {
	return string(e)
}

// ResourceType constants
const (
	ResourceTypeString  ResourceType = "String"
	ResourceTypeBoolean ResourceType = "Boolean"
	ResourceTypeInteger ResourceType = "Integer"
	ResourceTypeFloat   ResourceType = "Float"
	ResourceTypeObjLink ResourceType = "ObjLink"
	ResourceTypeOpaque  ResourceType = "Opaque"
)

// ResourceType represents the type of data
type ResourceType string

// String for ResourceType
func (e ResourceType) String() string {
	return string(e)
}
