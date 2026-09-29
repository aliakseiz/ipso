package ipso

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var registry Registry

func TestMain(m *testing.M) {
	reg, err := newTestRegistry()
	if err != nil {
		os.Exit(1)
	}

	registry = reg

	ret := m.Run()

	os.Exit(ret)
}

func newTestRegistry() (Registry, error) {
	reg, err := New(Configuration{})
	if err != nil {
		return nil, err
	}

	err = reg.Import("share/registry.yaml")
	if err != nil {
		return nil, err
	}

	return reg, nil
}

func TestReg_FindResourceByOIR(t *testing.T) {
	t.Parallel()

	res, err := registry.FindResourceByOIR("3303/0/5700")
	require.NoError(t, err)

	assert.Equal(t, int32(5700), res.ID)
}

func TestReg_FindResourceByOIRParseError(t *testing.T) {
	t.Parallel()

	_, err := registry.FindResourceByOIR("3303/05700")
	require.Error(t, err)

	assert.ErrorIs(t, err, errInvalidOIR)
}

func TestReg_FindResourceByOIRObjNotFound(t *testing.T) {
	t.Parallel()

	_, err := registry.FindResourceByOIR("-1/0/0")
	require.Error(t, err)

	assert.ErrorIs(t, err, errObjNotFound)
}

func TestReg_FindResourceByOIRResNotFound(t *testing.T) {
	t.Parallel()

	_, err := registry.FindResourceByOIR("3303/0/0")
	require.Error(t, err)

	assert.ErrorIs(t, err, errResNotFound)
}

func TestReg_ExportImportLeadingWhitespace(t *testing.T) {
	t.Parallel()

	objects := []Object{{
		Name:          "\tObject",
		Description1:  "\t\tLeading tabs\nsecond line  ",
		Description2:  "  leading spaces",
		ObjectID:      1,
		ObjectVersion: DefaultObjectVersion,
		Resources: ObjectResources{Item: []Resource{{
			ID:               0,
			Name:             " Resource ",
			Description:      "\t\tResource description\nsecond line",
			RangeEnumeration: "\t0-255",
			Units:            " m ",
		}}},
	}}
	reg := &Reg{Objects: objects}
	filename := filepath.Join(t.TempDir(), "registry.yaml")

	require.NoError(t, reg.Export(filename))

	// Export must not modify the registry.
	assert.Equal(t, "\tObject", reg.Objects[0].Name)
	assert.Equal(t, " Resource ", reg.Objects[0].Resources.Item[0].Name)

	imported := &Reg{}
	require.NoError(t, imported.Import(filename))
	require.Len(t, imported.Objects, 1)

	got := imported.Objects[0]
	assert.Equal(t, "Object", got.Name)
	assert.Equal(t, "Leading tabs\nsecond line", got.Description1)
	assert.Equal(t, "leading spaces", got.Description2)
	assert.Equal(t, "Resource", got.Resources.Item[0].Name)
	assert.Equal(t, "Resource description\nsecond line", got.Resources.Item[0].Description)
	assert.Equal(t, "0-255", got.Resources.Item[0].RangeEnumeration)
	assert.Equal(t, "m", got.Resources.Item[0].Units)
}
