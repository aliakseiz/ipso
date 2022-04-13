package ipso

import (
	"os"
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
