package ipso_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aliakseiz/ipso"
)

func TestDefaultConfiguration(t *testing.T) {
	t.Parallel()

	want := ipso.Configuration{
		InitOnNew:      true,
		SkipInitErrors: true,
		Sanitize:       true,
	}

	require.Equal(t, want, ipso.DefaultConfiguration())
}
