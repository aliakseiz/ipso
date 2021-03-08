package registry_test

import (
	"reflect"
	"testing"

	"github.com/aliakseiz/ipso-registry/registry"
)

func TestDefaultConfiguration(t *testing.T) {
	tests := []struct {
		name string
		want registry.Configuration
	}{
		{name: "ValidDefaultConfiguration", want: registry.Configuration{
			InitOnNew:      true,
			SkipInitErrors: true,
			Sanitize:       true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := registry.DefaultConfiguration(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}
