package ipso_test

import (
	"reflect"
	"testing"

	"github.com/aliakseiz/ipso"
)

func TestDefaultConfiguration(t *testing.T) {
	tests := []struct {
		name string
		want ipso.Configuration
	}{
		{name: "ValidDefaultConfiguration", want: ipso.Configuration{
			InitOnNew:      true,
			SkipInitErrors: true,
			Sanitize:       true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ipso.DefaultConfiguration(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}
