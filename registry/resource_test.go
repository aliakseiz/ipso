package registry

import (
	"reflect"
	"testing"

	"github.com/aliakseiz/lwm2m-registry/api/client"
)

func Test_mapResource(t *testing.T) {
	type args struct {
		omaResource openapi.Resource
	}
	tests := []struct {
		name string
		args args
		want Resource
	}{
		{
			name: "Valid resource",
			args: args{openapi.Resource{
				ID:                2,
				Name:              "TestName",
				Operations:        "RW",
				MultipleInstances: "Single",
				Mandatory:         "Mandatory",
				Type:              "Float",
				RangeEnumeration:  "1-2",
				Units:             "C",
				Description:       "Test description",
			}},
			want: Resource{
				ID:               2,
				Name:             "TestName",
				Operations:       OperationTypeReadWrite,
				MultipleInstance: InstanceTypeSingle,
				Mandatory:        MandatoryTypeMandatory,
				Type:             ResourceTypeFloat,
				RangeEnumeration: "1-2",
				Units:            "C",
				Description:      "Test description",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapResource(tt.args.omaResource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapResources(t *testing.T) {
	type args struct {
		omaResources []openapi.Resource
	}
	tests := []struct {
		name string
		args args
		want []Resource
	}{
		{
			name: "Single resource",
			args: args{[]openapi.Resource{
				{
					ID:                2,
					Name:              "TestName",
					Operations:        "R",
					MultipleInstances: "Single",
					Mandatory:         "Mandatory",
					Type:              "Integer",
					RangeEnumeration:  "0",
					Units:             "F",
					Description:       "Test description",
				},
			}},
			want: []Resource{
				{
					ID:               2,
					Name:             "TestName",
					Operations:       OperationTypeRead,
					MultipleInstance: InstanceTypeSingle,
					Mandatory:        MandatoryTypeMandatory,
					Type:             ResourceTypeInteger,
					RangeEnumeration: "0",
					Units:            "F",
					Description:      "Test description",
				},
			},
		},
		{
			name: "Multiple resources",
			args: args{[]openapi.Resource{
				{
					ID:                3,
					Name:              "TestName3",
					Operations:        "E",
					MultipleInstances: "Multiple",
					Mandatory:         "Optional",
					Type:              "Opaque",
					RangeEnumeration:  "",
					Units:             "",
					Description:       "Test 3 description",
				},
				{
					ID:                4,
					Name:              "TestName4",
					Operations:        "W",
					MultipleInstances: "Single",
					Mandatory:         "Optional",
					Type:              "Boolean",
					RangeEnumeration:  "On|Off",
					Units:             "",
					Description:       "Test 4 description",
				},
			}},
			want: []Resource{
				{
					ID:               3,
					Name:             "TestName3",
					Operations:       OperationTypeExecute,
					MultipleInstance: InstanceTypeMultiple,
					Mandatory:        MandatoryTypeOptional,
					Type:             ResourceTypeOpaque,
					RangeEnumeration: "",
					Units:            "",
					Description:      "Test 3 description",
				},
				{
					ID:               4,
					Name:             "TestName4",
					Operations:       OperationTypeWrite,
					MultipleInstance: InstanceTypeSingle,
					Mandatory:        MandatoryTypeOptional,
					Type:             ResourceTypeBoolean,
					RangeEnumeration: "On|Off",
					Units:            "",
					Description:      "Test 4 description",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapResources(tt.args.omaResources); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapResources() = %v, want %v", got, tt.want)
			}
		})
	}
}
