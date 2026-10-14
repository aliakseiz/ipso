package embedded_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aliakseiz/ipso/embedded"
)

func TestNew(t *testing.T) {
	t.Parallel()

	reg, err := embedded.New()
	require.NoError(t, err)
	require.Greater(t, len(reg.GetObjects()), 500)

	res, err := reg.FindResourceByOIR("3303/0/5700")
	require.NoError(t, err)
	require.Equal(t, "Sensor Value", res.Name)

	other, err := embedded.New()
	require.NoError(t, err)

	objects := other.GetObjects()
	require.NotEmpty(t, objects)
	require.Equal(t, reg.GetObjects()[0], objects[0], "calls give registries with equal data")

	other.Objects[0].Name = "changed"
	require.NotEqual(t, "changed", reg.Objects[0].Name, "registries are independent")
}

func TestEmbeddedMatchesShare(t *testing.T) {
	t.Parallel()

	want, err := os.ReadFile("../share/registry_sanitized.yaml")
	require.NoError(t, err)

	got, err := os.ReadFile("registry.yaml")
	require.NoError(t, err)

	require.True(t, bytes.Equal(want, got), "embedded/registry.yaml differs from share/registry_sanitized.yaml")
}
