package ipso_test

import (
	"testing"

	"github.com/aliakseiz/ipso"
)

func TestInstanceType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    ipso.InstanceType
		want string
	}{
		{
			name: "Single instance",
			e:    ipso.InstanceTypeSingle,
			want: "Single",
		},
		{
			name: "Multiple instances",
			e:    ipso.InstanceTypeMultiple,
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
		e    ipso.MandatoryType
		want string
	}{
		{
			name: "Mandatory type",
			e:    ipso.MandatoryTypeMandatory,
			want: "Mandatory",
		},
		{
			name: "Optional type",
			e:    ipso.MandatoryTypeOptional,
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
		e    ipso.OperationType
		want string
	}{
		{
			name: "Read operation",
			e:    ipso.OperationTypeRead,
			want: "R",
		},
		{
			name: "Write operation",
			e:    ipso.OperationTypeWrite,
			want: "W",
		},
		{
			name: "Read/write operations",
			e:    ipso.OperationTypeReadWrite,
			want: "RW",
		},
		{
			name: "Execute operation",
			e:    ipso.OperationTypeExecute,
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
		e    ipso.ResourceType
		want string
	}{
		{
			name: "Type boolean",
			e:    ipso.ResourceTypeBoolean,
			want: "Boolean",
		},
		{
			name: "Type string",
			e:    ipso.ResourceTypeString,
			want: "String",
		},
		{
			name: "Type integer",
			e:    ipso.ResourceTypeInteger,
			want: "Integer",
		},
		{
			name: "Type unsigned integer",
			e:    ipso.ResourceTypeUnsignedInteger,
			want: "Unsigned Integer",
		},
		{
			name: "Type float",
			e:    ipso.ResourceTypeFloat,
			want: "Float",
		},
		{
			name: "Type object link",
			e:    ipso.ResourceTypeObjLink,
			want: "Objlnk",
		},
		{
			name: "Type core link",
			e:    ipso.ResourceTypeCoreLink,
			want: "Corelnk",
		},
		{
			name: "Type opaque",
			e:    ipso.ResourceTypeOpaque,
			want: "Opaque",
		},
		{
			name: "Type time",
			e:    ipso.ResourceTypeTime,
			want: "Time",
		},
		{
			name: "Type none (executable)",
			e:    ipso.ResourceTypeNone,
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
		e    ipso.DifferenceType
		want string
	}{
		{
			name: "Object added",
			e:    ipso.DifferenceTypeNewObject,
			want: "New object added",
		},
		{
			name: "Object removed",
			e:    ipso.DifferenceTypeObjectRemoved,
			want: "Object was removed",
		},
		{
			name: "Unknown difference type",
			e:    ipso.DifferenceTypeUnknown,
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
