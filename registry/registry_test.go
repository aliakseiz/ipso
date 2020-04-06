package registry

import (
	"reflect"
	"testing"
)

func TestDefaultRegistryConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Configuration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultConfiguration(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		cfg *Configuration
	}
	tests := []struct {
		name string
		args args
		want *Registry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}