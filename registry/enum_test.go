package registry

import "testing"

func TestInstanceType_String(t *testing.T) {
	tests := []struct {
		name string
		e    InstanceType
		want string
	}{
		{
			name: "Single instance",
			e:    InstanceTypeSingle,
			want: "Single",
		},
		{
			name: "Multiple instances",
			e:    InstanceTypeMultiple,
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
		e    MandatoryType
		want string
	}{
		{
			name: "Mandatory type",
			e:    MandatoryTypeMandatory,
			want: "Mandatory",
		},
		{
			name: "Optional type",
			e:    MandatoryTypeOptional,
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
		e    OperationType
		want string
	}{
		{
			name: "Read operation",
			e:    OperationTypeRead,
			want: "R",
		},
		{
			name: "Write operation",
			e:    OperationTypeWrite,
			want: "W",
		},
		{
			name: "Read/write operations",
			e:    OperationTypeReadWrite,
			want: "RW",
		},
		{
			name: "Execute operation",
			e:    OperationTypeExecute,
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
		e    ResourceType
		want string
	}{
		{
			name: "Type boolean",
			e:    ResourceTypeBoolean,
			want: "Boolean",
		},
		{
			name: "Type string",
			e:    ResourceTypeString,
			want: "String",
		},
		{
			name: "Type integer",
			e:    ResourceTypeInteger,
			want: "Integer",
		},
		{
			name: "Type float",
			e:    ResourceTypeFloat,
			want: "Float",
		},
		{
			name: "Type object link",
			e:    ResourceTypeObjLink,
			want: "ObjLink",
		},
		{
			name: "Type opaque",
			e:    ResourceTypeOpaque,
			want: "Opaque",
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
