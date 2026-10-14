package ipso_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aliakseiz/ipso"
)

func TestEnums_String(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		e    fmt.Stringer
		want string
	}{
		"InstanceType/single":             {ipso.InstanceTypeSingle, "Single"},
		"InstanceType/multiple":           {ipso.InstanceTypeMultiple, "Multiple"},
		"MandatoryType/mandatory":         {ipso.MandatoryTypeMandatory, "Mandatory"},
		"MandatoryType/optional":          {ipso.MandatoryTypeOptional, "Optional"},
		"OperationType/read":              {ipso.OperationTypeRead, "R"},
		"OperationType/write":             {ipso.OperationTypeWrite, "W"},
		"OperationType/read-write":        {ipso.OperationTypeReadWrite, "RW"},
		"OperationType/execute":           {ipso.OperationTypeExecute, "E"},
		"OperationType/empty":             {ipso.OperationTypeEmpty, ""},
		"ResourceType/boolean":            {ipso.ResourceTypeBoolean, "Boolean"},
		"ResourceType/string":             {ipso.ResourceTypeString, "String"},
		"ResourceType/integer":            {ipso.ResourceTypeInteger, "Integer"},
		"ResourceType/unsigned-integer":   {ipso.ResourceTypeUnsignedInteger, "Unsigned Integer"},
		"ResourceType/float":              {ipso.ResourceTypeFloat, "Float"},
		"ResourceType/object-link":        {ipso.ResourceTypeObjLink, "Objlnk"},
		"ResourceType/core-link":          {ipso.ResourceTypeCoreLink, "Corelnk"},
		"ResourceType/opaque":             {ipso.ResourceTypeOpaque, "Opaque"},
		"ResourceType/time":               {ipso.ResourceTypeTime, "Time"},
		"ResourceType/empty (executable)": {ipso.ResourceTypeEmpty, ""},
		"ResourceType/none":               {ipso.ResourceTypeNone, "none"},
		"DifferenceType/new-object":       {ipso.DifferenceTypeNewObject, "New object added"},
		"DifferenceType/object-removed":   {ipso.DifferenceTypeObjectRemoved, "Object was removed"},
		"DifferenceType/unknown":          {ipso.DifferenceTypeUnknown, "Unknown reason"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, tt.e.String())
		})
	}
}
