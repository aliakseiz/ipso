package registry_test

import (
	"testing"

	"github.com/aliakseiz/ipso-registry/registry"
)

func TestInstanceType_String(t *testing.T) {
	tests := []struct {
		name string
		e    registry.InstanceType
		want string
	}{
		{
			name: "Single instance",
			e:    registry.InstanceTypeSingle,
			want: "Single",
		},
		{
			name: "Multiple instances",
			e:    registry.InstanceTypeMultiple,
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
	tests := []struct {
		name string
		e    registry.MandatoryType
		want string
	}{
		{
			name: "Mandatory type",
			e:    registry.MandatoryTypeMandatory,
			want: "Mandatory",
		},
		{
			name: "Optional type",
			e:    registry.MandatoryTypeOptional,
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
	tests := []struct {
		name string
		e    registry.OperationType
		want string
	}{
		{
			name: "Read operation",
			e:    registry.OperationTypeRead,
			want: "R",
		},
		{
			name: "Write operation",
			e:    registry.OperationTypeWrite,
			want: "W",
		},
		{
			name: "Read/write operations",
			e:    registry.OperationTypeReadWrite,
			want: "RW",
		},
		{
			name: "Execute operation",
			e:    registry.OperationTypeExecute,
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
	tests := []struct {
		name string
		e    registry.ResourceType
		want string
	}{
		{
			name: "Type boolean",
			e:    registry.ResourceTypeBoolean,
			want: "Boolean",
		},
		{
			name: "Type string",
			e:    registry.ResourceTypeString,
			want: "String",
		},
		{
			name: "Type integer",
			e:    registry.ResourceTypeInteger,
			want: "Integer",
		},
		{
			name: "Type unsigned integer",
			e:    registry.ResourceTypeUnsignedInteger,
			want: "Unsigned Integer",
		},
		{
			name: "Type float",
			e:    registry.ResourceTypeFloat,
			want: "Float",
		},
		{
			name: "Type object link",
			e:    registry.ResourceTypeObjLink,
			want: "Objlnk",
		},
		{
			name: "Type core link",
			e:    registry.ResourceTypeCoreLink,
			want: "Corelnk",
		},
		{
			name: "Type opaque",
			e:    registry.ResourceTypeOpaque,
			want: "Opaque",
		},
		{
			name: "Type time",
			e:    registry.ResourceTypeTime,
			want: "Time",
		},
		{
			name: "Type none (executable)",
			e:    registry.ResourceTypeNone,
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
	tests := []struct {
		name string
		e    registry.DifferenceType
		want string
	}{
		{
			name: "Object added",
			e:    registry.DifferenceTypeNewObject,
			want: "New object added",
		},
		{
			name: "Object removed",
			e:    registry.DifferenceTypeObjectRemoved,
			want: "Object was removed",
		},
		{
			name: "Unknown difference type",
			e:    registry.DifferenceTypeUnknown,
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
