package ipso_registry_test

import (
	"reflect"
	"testing"

	ipso_registry "github.com/aliakseiz/ipso-registry"
)

func TestDefaultConfiguration(t *testing.T) {
	tests := []struct {
		name string
		want ipso_registry.Configuration
	}{
		{name: "ValidDefaultConfiguration", want: ipso_registry.Configuration{
			InitOnNew:      true,
			SkipInitErrors: true,
			Sanitize:       true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ipso_registry.DefaultConfiguration(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}
