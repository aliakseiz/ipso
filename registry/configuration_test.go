package registry

import (
	"reflect"
	"testing"
)

func TestDefaultConfiguration(t *testing.T) {
	tests := []struct {
		name string
		want *Configuration
	}{
		{name: "ValidDefaultConfiguration", want: &Configuration{
			InitOnNew:      true,
			SkipInitErrors: true,
			Sanitize:       true,
			Sanitizer:      DefaultSanitizer(),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultConfiguration(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}
