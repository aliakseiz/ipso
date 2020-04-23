package registry

import (
	"reflect"
	"testing"

	openapi "github.com/aliakseiz/ipso-registry/api/client"
)

func Test_mapObject(t *testing.T) {
	type args struct {
		omaObject *openapi.Object
	}
	tests := []struct {
		name string
		args args
		want Object
	}{
		{
			name: "Object without resources",
			args: args{&openapi.Object{
				Name:              "TestName",
				Description1:      "TestDescription1",
				Description2:      "TestDescription2",
				ObjectID:          2,
				ObjectURN:         "object:urn",
				LWM2MVersion:      "1.1",
				ObjectVersion:     "1.2",
				MultipleInstances: "Multiple",
				Mandatory:         "Optional",
				Resources:         openapi.ObjectResources{},
			}},
			want: Object{
				Name:             "TestName",
				Description1:     "TestDescription1",
				Description2:     "TestDescription2",
				ObjectID:         2,
				ObjectURN:        "object:urn",
				LwM2MVersion:     "1.1",
				ObjectVersion:    "1.2",
				MultipleInstance: InstanceTypeMultiple,
				Mandatory:        MandatoryTypeOptional,
				Resources:        []Resource{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapObject(tt.args.omaObject); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
