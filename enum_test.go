package ipso_registry_test

import (
	"testing"

	ipso_registry "github.com/aliakseiz/ipso-registry"
)

func TestInstanceType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    ipso_registry.InstanceType
		want string
	}{
		{
			name: "Single instance",
			e:    ipso_registry.InstanceTypeSingle,
			want: "Single",
		},
		{
			name: "Multiple instances",
			e:    ipso_registry.InstanceTypeMultiple,
			want: "Multiple",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMandatoryType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    ipso_registry.MandatoryType
		want string
	}{
		{
			name: "Mandatory type",
			e:    ipso_registry.MandatoryTypeMandatory,
			want: "Mandatory",
		},
		{
			name: "Optional type",
			e:    ipso_registry.MandatoryTypeOptional,
			want: "Optional",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperationType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    ipso_registry.OperationType
		want string
	}{
		{
			name: "Read operation",
			e:    ipso_registry.OperationTypeRead,
			want: "R",
		},
		{
			name: "Write operation",
			e:    ipso_registry.OperationTypeWrite,
			want: "W",
		},
		{
			name: "Read/write operations",
			e:    ipso_registry.OperationTypeReadWrite,
			want: "RW",
		},
		{
			name: "Execute operation",
			e:    ipso_registry.OperationTypeExecute,
			want: "E",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    ipso_registry.ResourceType
		want string
	}{
		{
			name: "Type boolean",
			e:    ipso_registry.ResourceTypeBoolean,
			want: "Boolean",
		},
		{
			name: "Type string",
			e:    ipso_registry.ResourceTypeString,
			want: "String",
		},
		{
			name: "Type integer",
			e:    ipso_registry.ResourceTypeInteger,
			want: "Integer",
		},
		{
			name: "Type unsigned integer",
			e:    ipso_registry.ResourceTypeUnsignedInteger,
			want: "Unsigned Integer",
		},
		{
			name: "Type float",
			e:    ipso_registry.ResourceTypeFloat,
			want: "Float",
		},
		{
			name: "Type object link",
			e:    ipso_registry.ResourceTypeObjLink,
			want: "Objlnk",
		},
		{
			name: "Type core link",
			e:    ipso_registry.ResourceTypeCoreLink,
			want: "Corelnk",
		},
		{
			name: "Type opaque",
			e:    ipso_registry.ResourceTypeOpaque,
			want: "Opaque",
		},
		{
			name: "Type time",
			e:    ipso_registry.ResourceTypeTime,
			want: "Time",
		},
		{
			name: "Type none (executable)",
			e:    ipso_registry.ResourceTypeNone,
			want: "none",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDifferenceType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    ipso_registry.DifferenceType
		want string
	}{
		{
			name: "Object added",
			e:    ipso_registry.DifferenceTypeNewObject,
			want: "New object added",
		},
		{
			name: "Object removed",
			e:    ipso_registry.DifferenceTypeObjectRemoved,
			want: "Object was removed",
		},
		{
			name: "Unknown difference type",
			e:    ipso_registry.DifferenceTypeUnknown,
			want: "Unknown reason",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
