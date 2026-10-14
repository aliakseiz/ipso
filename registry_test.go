package ipso_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aliakseiz/ipso"
)

var registry *ipso.Reg

func TestMain(m *testing.M) {
	reg, err := ipso.New(ipso.Configuration{})
	if err != nil {
		os.Exit(1)
	}

	err = reg.ImportFile("share/registry.yaml")
	if err != nil {
		os.Exit(1)
	}

	registry = reg

	ret := m.Run()

	os.Exit(ret)
}

func TestReg_FindResourceByOIR(t *testing.T) {
	t.Parallel()

	res, err := registry.FindResourceByOIR("3303/0/5700")
	require.NoError(t, err)

	require.Equal(t, int32(5700), res.ID)
}

func TestReg_FindResourceByOIRParseError(t *testing.T) {
	t.Parallel()

	_, err := registry.FindResourceByOIR("3303/05700")
	require.Error(t, err)

	require.ErrorIs(t, err, ipso.ErrInvalidOIR)
}

func TestReg_FindResourceByOIRObjNotFound(t *testing.T) {
	t.Parallel()

	_, err := registry.FindResourceByOIR("9999/0/0")
	require.Error(t, err)

	require.ErrorIs(t, err, ipso.ErrObjectNotFound)
}

func TestReg_FindResourceByOIRLeadingSlash(t *testing.T) {
	t.Parallel()

	res, err := registry.FindResourceByOIR("/3303/0/5700")
	require.NoError(t, err)

	require.Equal(t, int32(5700), res.ID)
}

func TestReg_FindResourceByOIRAndVer(t *testing.T) {
	t.Parallel()

	// Object 4 gained resource 11 in version 1.1 and resource 12 in version 1.2.
	tests := []struct {
		name    string
		oir     string
		ver     string
		wantID  int32
		wantErr error
	}{
		{name: "resource of the requested version", oir: "4/0/11", ver: "1.1", wantID: 11},
		{name: "leading slash", oir: "/4/0/12", ver: "1.2", wantID: 12},
		{name: "resource absent in the version", oir: "4/0/11", ver: "1.0", wantErr: ipso.ErrResourceNotFound},
		{name: "version absent", oir: "4/0/11", ver: "9.9", wantErr: ipso.ErrObjectNotFound},
		{name: "object absent", oir: "9999/0/0", ver: "1.0", wantErr: ipso.ErrObjectNotFound},
		{name: "invalid OIR", oir: "4/11", ver: "1.1", wantErr: ipso.ErrInvalidOIR},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := registry.FindResourceByOIRAndVer(tt.oir, tt.ver)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantID, res.ID)
		})
	}
}

func TestReg_FindResourceByOIRUsesDefaultVersion(t *testing.T) {
	t.Parallel()

	// Resource 11 of object 4 exists in version 1.1 but not in 1.0, the version FindResourceByOIR uses.
	_, err := registry.FindResourceByOIR("4/0/11")
	require.ErrorIs(t, err, ipso.ErrResourceNotFound)

	got, err := registry.FindResourceByOIR("4/0/0")
	require.NoError(t, err)

	want, err := registry.FindResourceByOIRAndVer("4/0/0", ipso.DefaultObjectVersion)
	require.NoError(t, err)

	require.Equal(t, want, got)
}

func Test_parseOIRString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		oir     string
		obj     int32
		inst    int32
		res     int32
		wantErr bool
	}{
		{name: "plain", oir: "3303/0/5700", obj: 3303, inst: 0, res: 5700},
		{name: "leading slash", oir: "/3303/1/5700", obj: 3303, inst: 1, res: 5700},
		{name: "max int32", oir: "2147483647/2147483647/2147483647", obj: 2147483647, inst: 2147483647, res: 2147483647},
		{name: "empty", oir: "", wantErr: true},
		{name: "too few segments", oir: "3303/5700", wantErr: true},
		{name: "too many segments", oir: "3303/0/5700/1", wantErr: true},
		{name: "trailing slash", oir: "3303/0/5700/", wantErr: true},
		{name: "double leading slash", oir: "//3303/0/5700", wantErr: true},
		{name: "empty segment", oir: "3303//5700", wantErr: true},
		{name: "not a number", oir: "3303/0/abc", wantErr: true},
		{name: "negative object", oir: "-1/0/0", wantErr: true},
		{name: "negative instance", oir: "3303/-1/0", wantErr: true},
		{name: "negative resource", oir: "3303/0/-1", wantErr: true},
		{name: "plus sign", oir: "+3303/0/0", wantErr: true},
		{name: "overflows int32", oir: "2147483648/0/0", wantErr: true},
		{name: "overflows int64", oir: "99999999999999999999/0/0", wantErr: true},
		{name: "wraps to small int32", oir: "4294967299/0/0", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			obj, inst, res, err := ipso.ParseOIRString(tt.oir)
			if tt.wantErr {
				require.Error(t, err)
				require.ErrorIs(t, err, ipso.ErrInvalidOIR)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.obj, obj)
			require.Equal(t, tt.inst, inst)
			require.Equal(t, tt.res, res)
		})
	}
}

func FuzzParseOIRString(f *testing.F) {
	for _, seed := range []string{
		"3303/0/5700", "/3303/0/5700", "2147483647/2147483647/2147483647", "2147483648/0/0",
		"", "/", "//", "1//2", "-1/0/0", "+1/0/0", "0x10/0/0", "1/2/3/4", "007/0/0", "١/٢/٣",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, oir string) {
		obj, inst, res, err := ipso.ParseOIRString(oir)
		if err != nil {
			require.ErrorIs(t, err, ipso.ErrInvalidOIR)

			return
		}

		require.GreaterOrEqual(t, obj, int32(0))
		require.GreaterOrEqual(t, inst, int32(0))
		require.GreaterOrEqual(t, res, int32(0))

		obj2, inst2, res2, err := ipso.ParseOIRString(fmt.Sprintf("%d/%d/%d", obj, inst, res))
		require.NoError(t, err)
		require.Equal(t, [3]int32{obj, inst, res}, [3]int32{obj2, inst2, res2})
	})
}

func TestReg_FindResourceByOIRNegativeID(t *testing.T) {
	t.Parallel()

	_, err := registry.FindResourceByOIR("-1/0/0")
	require.Error(t, err)

	require.ErrorIs(t, err, ipso.ErrInvalidOIR)
}

func TestReg_FindResourceByOIRResNotFound(t *testing.T) {
	t.Parallel()

	_, err := registry.FindResourceByOIR("3303/0/0")
	require.Error(t, err)

	require.ErrorIs(t, err, ipso.ErrResourceNotFound)
}

func TestReg_ExportImportLeadingWhitespace(t *testing.T) {
	t.Parallel()

	objects := []ipso.Object{{
		Name:          "\tObject",
		Description1:  "\t\tLeading tabs\nsecond line  ",
		Description2:  "  leading spaces",
		ObjectID:      1,
		ObjectVersion: ipso.DefaultObjectVersion,
		Resources: ipso.ObjectResources{Item: []ipso.Resource{{
			ID:               0,
			Name:             " Resource ",
			Description:      "\t\tResource description\nsecond line",
			RangeEnumeration: "\t0-255",
			Units:            " m ",
		}}},
	}}
	reg := &ipso.Reg{Objects: objects}
	filename := filepath.Join(t.TempDir(), "registry.yaml")

	require.NoError(t, reg.ExportFile(filename))

	// Export must not modify the registry.
	require.Equal(t, "\tObject", reg.Objects[0].Name)
	require.Equal(t, " Resource ", reg.Objects[0].Resources.Item[0].Name)

	imported := &ipso.Reg{}
	require.NoError(t, imported.ImportFile(filename))
	require.Len(t, imported.Objects, 1)

	got := imported.Objects[0]
	require.Equal(t, "Object", got.Name)
	require.Equal(t, "Leading tabs\nsecond line", got.Description1)
	require.Equal(t, "leading spaces", got.Description2)
	require.Equal(t, "Resource", got.Resources.Item[0].Name)
	require.Equal(t, "Resource description\nsecond line", got.Resources.Item[0].Description)
	require.Equal(t, "0-255", got.Resources.Item[0].RangeEnumeration)
	require.Equal(t, "m", got.Resources.Item[0].Units)
}

func TestRegistry_ExportImportRoundTrip(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("parses the full registry")
	}

	for _, file := range []string{"share/registry.yaml", "share/registry_sanitized.yaml"} {
		t.Run(file, func(t *testing.T) {
			t.Parallel()

			first := &ipso.Reg{}
			require.NoError(t, first.ImportFile(file))

			var exported bytes.Buffer
			require.NoError(t, first.Export(&exported))

			second := &ipso.Reg{}
			require.NoError(t, second.Import(bytes.NewReader(exported.Bytes())))
			require.Equal(t, first.Objects, second.Objects, "objects changed by Export and Import")

			var again bytes.Buffer
			require.NoError(t, second.Export(&again))
			require.True(t, bytes.Equal(exported.Bytes(), again.Bytes()), "second Export differs from the first")
		})
	}
}

func TestReg_ImportErrorKeepsState(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"syntax error": "- Name: [unclosed\n",
		"type error":   "- Name: Broken\n  ObjectID: not-a-number\n- Name: Other\n  ObjectID: 2\n",
	}

	for name, content := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			reg := &ipso.Reg{}
			require.NoError(t, reg.ImportFile("share/registry.yaml"))

			before := len(reg.Objects)
			require.NotZero(t, before)

			filename := filepath.Join(t.TempDir(), "broken.yaml")
			require.NoError(t, os.WriteFile(filename, []byte(content), 0o600))

			require.Error(t, reg.ImportFile(filename))

			require.Len(t, reg.Objects, before)

			obj, err := reg.FindObjectByIDAndVer(3303, ipso.DefaultObjectVersion)
			require.NoError(t, err)
			require.Equal(t, int32(3303), obj.ObjectID)
		})
	}
}

func TestReg_ImportExportReaderWriter(t *testing.T) {
	t.Parallel()

	reg := newSmallReg()

	var buf bytes.Buffer

	require.NoError(t, reg.Export(&buf))

	imported := &ipso.Reg{}
	require.NoError(t, imported.Import(&buf))
	require.Len(t, imported.Objects, len(reg.Objects))
	require.Empty(t, imported.Compare(reg))

	_, err := imported.FindResourceByOIR("3303/0/5700")
	require.NoError(t, err, "index is rebuilt by Import")
}

func TestReg_ImportReadError(t *testing.T) {
	t.Parallel()

	reg := newSmallReg()

	require.ErrorIs(t, reg.Import(iotest.ErrReader(errReadWrite)), errReadWrite)
	require.Len(t, reg.Objects, 2)
	require.NotNil(t, reg.Lookup(3303, ipso.DefaultObjectVersion))
}

func TestReg_ExportWriteError(t *testing.T) {
	t.Parallel()

	reg := &ipso.Reg{Objects: []ipso.Object{{ObjectID: 1}}}

	require.ErrorIs(t, reg.Export(failingWriter{}), errReadWrite)
}

func TestReg_ImportExportFileErrors(t *testing.T) {
	t.Parallel()

	reg := &ipso.Reg{Objects: []ipso.Object{{ObjectID: 1}}}
	dir := t.TempDir()

	require.ErrorIs(t, reg.ImportFile(""), ipso.ErrEmptyFilename)
	require.ErrorIs(t, reg.ExportFile(""), ipso.ErrEmptyFilename)

	require.ErrorIs(t, reg.ImportFile(filepath.Join(dir, "missing.yaml")), os.ErrNotExist)
	require.ErrorIs(t, reg.ExportFile(filepath.Join(dir, "missing-dir", "registry.yaml")), os.ErrNotExist)
	require.Len(t, reg.Objects, 1, "failed ImportFile keeps the registry")
}

func TestReg_ExportFileTruncates(t *testing.T) {
	t.Parallel()

	filename := filepath.Join(t.TempDir(), "registry.yaml")
	require.NoError(t, os.WriteFile(filename, bytes.Repeat([]byte("# stale\n"), 100), 0o600))

	reg := &ipso.Reg{Objects: []ipso.Object{{ObjectID: 1, ObjectVersion: ipso.DefaultObjectVersion}}}
	require.NoError(t, reg.ExportFile(filename))

	imported := &ipso.Reg{}
	require.NoError(t, imported.ImportFile(filename))
	require.Len(t, imported.Objects, 1)
}

func TestReg_SanitizeObjectAndResourceFields(t *testing.T) {
	t.Parallel()

	reg := &ipso.Reg{Objects: []ipso.Object{{
		ObjectID:      1,
		ObjectVersion: ipso.DefaultObjectVersion,
		Description1:  "  first\tdescription  ",
		Description2:  "second#description",
		Resources: ipso.ObjectResources{Item: []ipso.Resource{{
			ID:               0,
			Description:      " res\tdesc ",
			RangeEnumeration: " 0-255\t",
			Units:            " m\t",
		}}},
	}}}

	reg.Sanitize([]string{"\t", "#"})

	got := reg.Objects[0]
	require.Equal(t, "first description", got.Description1)
	require.Equal(t, "seconddescription", got.Description2)
	require.Equal(t, "res desc", got.Resources.Item[0].Description)
	require.Equal(t, "0-255", got.Resources.Item[0].RangeEnumeration)
	require.Equal(t, "m", got.Resources.Item[0].Units)
}

func TestReg_SanitizeText(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		text      string
		sanitizer []string
		want      string
	}{
		"newline keeps sentences apart": {
			text: "Server.\nThese Resources", sanitizer: []string{"\n"}, want: "Server. These Resources",
		},
		"runs of spaces collapse": {
			text: "a   b \n\n c", sanitizer: []string{"\n", "  "}, want: "a b c",
		},
		"other strings are removed": {
			text: "• one • two", sanitizer: []string{"•"}, want: "one  two",
		},
		"no whitespace string keeps inner spaces": {
			text: "a   b#c", sanitizer: []string{"#"}, want: "a   bc",
		},
		"edges are trimmed": {
			text: "  \tx\t  ", sanitizer: []string{"\t"}, want: "x",
		},
		"empty entry is ignored": {
			text: "abc", sanitizer: []string{""}, want: "abc",
		},
		"empty sanitizer only trims": {
			text: " a \n b ", sanitizer: nil, want: "a \n b",
		},
		"default sanitizer": {
			text: "Line one.\n\tLine  two \\ • end\n", sanitizer: ipso.DefaultSanitizer(), want: "Line one. Line two end",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, ipso.SanitizeText(tt.text, tt.sanitizer))
		})
	}
}

func TestRegistry_SanitizedFileIsCurrent(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("parses the full registry")
	}

	reg := &ipso.Reg{}
	require.NoError(t, reg.ImportFile("share/registry.yaml"))
	reg.Sanitize(ipso.DefaultSanitizer())

	var buf bytes.Buffer
	require.NoError(t, reg.Export(&buf))

	want, err := os.ReadFile("share/registry_sanitized.yaml")
	require.NoError(t, err)

	// Regenerate with: Import share/registry.yaml, Sanitize(DefaultSanitizer()), ExportFile share/registry_sanitized.yaml,
	// then copy the file to embedded/registry.yaml.
	require.True(t, bytes.Equal(want, buf.Bytes()), "share/registry_sanitized.yaml is not the sanitized share/registry.yaml")
}

func TestReg_Compare(t *testing.T) {
	t.Parallel()

	for _, tt := range compareCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, tt.a.Compare(tt.b))
		})
	}
}

func Test_compareVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		a, b string
		want int
	}{
		{"1.0", "1.0", 0},
		{"1.0", "1.1", -1},
		{"1.2", "1.10", -1},
		{"1.10", "1.2", 1},
		{"2.0", "1.10", 1},
		{"1", "1.0", -1},
		{"1.0.1", "1.0", 1},
		{"", "1.0", -1},
		{"1.0", "", 1},
		{"", "", 0},
		{"1.0", "1.x", -1},
		{"1.x", "1.y", -1},
		{"1.x", "1.x", 0},
	}

	for _, tt := range tests {
		t.Run(tt.a+" vs "+tt.b, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, ipso.CompareVersions(tt.a, tt.b))
		})
	}
}

func TestReg_FindObjectsByIDEqualVersions(t *testing.T) {
	t.Parallel()

	objects := []ipso.Object{
		{ObjectID: 7, ObjectVersion: "1.00"},
		{ObjectID: 7, ObjectVersion: "1.0"},
	}
	reg := ipso.NewRegFromObjects(objects)

	for range 50 {
		got, err := reg.FindObjectsByID(7)
		require.NoError(t, err)
		require.Equal(t, "1.0", got[0].ObjectVersion)
		require.Equal(t, "1.00", got[1].ObjectVersion)
	}
}

func TestReg_FindObjectsByIDOrder(t *testing.T) {
	t.Parallel()

	reg := newVersionedReg()
	want := []string{"1.0", "1.1", "1.2", "1.10", "2.0"}

	// Map iteration order is random, repeat to catch an unsorted result.
	for range 50 {
		objects, err := reg.FindObjectsByID(7)
		require.NoError(t, err)

		got := make([]string, 0, len(objects))
		for _, o := range objects {
			got = append(got, o.ObjectVersion)
		}

		require.Equal(t, want, got)
	}
}

func TestReg_FindResourcesByObjResIDsOrder(t *testing.T) {
	t.Parallel()

	reg := newVersionedReg()
	want := []string{"1.0", "1.1", "1.2", "1.10", "2.0"}

	for range 50 {
		resources, err := reg.FindResourcesByObjResIDs(7, 1)
		require.NoError(t, err)

		got := make([]string, 0, len(resources))
		for _, res := range resources {
			got = append(got, res.Name)
		}

		require.Equal(t, want, got)
	}
}

func TestReg_FindObjectsByName(t *testing.T) {
	t.Parallel()

	reg := newNamedReg()

	// The fixture is unordered, so the result order comes from the explicit sort.
	got, err := reg.FindObjectsByName("temperature")
	require.NoError(t, err)

	type idVer struct {
		id  int32
		ver string
	}

	ids := make([]idVer, 0, len(got))
	for _, o := range got {
		ids = append(ids, idVer{o.ObjectID, o.ObjectVersion})
	}

	require.Equal(t, []idVer{{3, "1.0"}, {3303, "1.0"}, {3303, "1.2"}, {3303, "1.10"}}, ids)

	for _, name := range []string{"Humidity", ""} {
		_, err = reg.FindObjectsByName(name)
		require.ErrorIs(t, err, ipso.ErrObjectNotFound, "name %q", name)
	}
}

func TestReg_FindResourcesByName(t *testing.T) {
	t.Parallel()

	reg := newNamedReg()

	got, err := reg.FindResourcesByName("sensor value")
	require.NoError(t, err)

	ids := make([]int32, 0, len(got))
	for _, res := range got {
		ids = append(ids, res.ID)
	}

	// Ordered by object ID, then object version (3/1.0, 3303/1.0, 3303/1.10).
	require.Equal(t, []int32{5700, 5701, 5700}, ids)

	for _, name := range []string{"Absent", ""} {
		_, err = reg.FindResourcesByName(name)
		require.ErrorIs(t, err, ipso.ErrResourceNotFound, "name %q", name)
	}
}

func TestReg_FindByNameShippedRegistry(t *testing.T) {
	t.Parallel()

	reg := registry

	objects, err := reg.FindObjectsByName("TEMPERATURE")
	require.NoError(t, err)
	require.NotEmpty(t, objects)

	for _, o := range objects {
		require.Equal(t, int32(3303), o.ObjectID)
	}

	resources, err := reg.FindResourcesByName("sensor value")
	require.NoError(t, err)
	require.Greater(t, len(resources), 10)
}

func TestReg_ReturnedObjectsDoNotAliasRegistry(t *testing.T) {
	t.Parallel()

	for name, get := range objectGetters() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			reg := newAliasReg()
			got := get(t, reg)

			require.Len(t, got.Resources.Item, 1)
			got.Resources.Item[0].Name = "changed"
			got.Resources.Item = append(got.Resources.Item[:0], ipso.Resource{ID: 99})

			require.Equal(t, aliasOrigName, reg.Objects[0].Resources.Item[0].Name)
			require.Equal(t, int32(1), reg.Objects[0].Resources.Item[0].ID)
		})
	}
}

func TestObject_clone(t *testing.T) {
	t.Parallel()

	require.Nil(t, ipso.CloneObject(ipso.Object{}).Resources.Item, "nil stays nil")

	empty := ipso.CloneObject(ipso.Object{Resources: ipso.ObjectResources{Item: []ipso.Resource{}}})
	require.NotNil(t, empty.Resources.Item)
	require.Empty(t, empty.Resources.Item)
}

func TestRegistry_ShippedDataInvariants(t *testing.T) {
	t.Parallel()

	type objectKey struct {
		id  int32
		ver string
	}

	known := newKnownValues()
	objects := registry.GetObjects()
	seen := make(map[objectKey]bool, len(objects))
	emptyTypes := 0

	for _, obj := range objects {
		key := objectKey{obj.ObjectID, obj.ObjectVersion}
		where := fmt.Sprintf("object %d version %q", key.id, key.ver)

		require.False(t, seen[key], "duplicate %s", where)
		seen[key] = true

		require.NotEmpty(t, obj.Name, "%s has no name", where)
		require.NotEmpty(t, obj.ObjectVersion, "%s has no version", where)
		require.True(t, known.instances[obj.MultipleInstances], "%s: unknown instance type %q", where, obj.MultipleInstances)
		require.True(t, known.mandatory[obj.Mandatory], "%s: unknown mandatory type %q", where, obj.Mandatory)

		resIDs := make(map[int32]bool, len(obj.Resources.Item))

		for _, res := range obj.Resources.Item {
			require.False(t, resIDs[res.ID], "%s: duplicate resource %d", where, res.ID)
			resIDs[res.ID] = true

			known.assertResource(t, fmt.Sprintf("%s resource %d", where, res.ID), res)

			if res.Type == ipso.ResourceTypeEmpty {
				emptyTypes++
			}
		}
	}

	// Executable resources carry no data type in the OMA registry.
	require.NotZero(t, emptyTypes)
}

func TestReg_ConcurrentReads(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	other := registry.GetRegistry()

	var wg sync.WaitGroup

	for i := range 4 {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			_, err := registry.FindResourceByOIR("3303/0/5700")
			assert.NoError(t, err)

			_, err = registry.FindResourceByObjIDObjVerResID(3303, ipso.DefaultObjectVersion, 5700)
			assert.NoError(t, err)

			_, err = registry.FindObjectsByID(3303)
			assert.NoError(t, err)

			_, err = registry.FindResourcesByID(5700)
			assert.NoError(t, err)

			_, err = registry.FindResourcesByObjResIDs(3303, 5700)
			assert.NoError(t, err)

			assert.NotEmpty(t, registry.GetObjects())
			assert.Empty(t, registry.Compare(other))

			// Export is slow under the race detector, two concurrent exports are enough.
			if i < 2 {
				assert.NoError(t, registry.ExportFile(filepath.Join(dir, fmt.Sprintf("registry-%d.yaml", i))))
			}
		}(i)
	}

	wg.Wait()
}

func newVersionedReg() *ipso.Reg {
	versions := []string{"1.10", "1.2", "2.0", "1.0", "1.1"}
	objects := make([]ipso.Object, 0, len(versions))

	for i, ver := range versions {
		objects = append(objects, ipso.Object{
			Name:          "obj",
			ObjectID:      7,
			ObjectVersion: ver,
			Resources:     ipso.ObjectResources{Item: []ipso.Resource{{ID: 1, Name: ver}, {ID: int32(100 + i)}}},
		})
	}

	return ipso.NewRegFromObjects(objects)
}

var errReadWrite = errors.New("stream failure")

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errReadWrite }

// newNamedReg returns an unordered registry where several objects and resources share a name.
func newNamedReg() *ipso.Reg {
	item := func(id int32, name string) ipso.ObjectResources {
		return ipso.ObjectResources{Item: []ipso.Resource{{ID: id, Name: name}, {ID: 1, Name: "other"}}}
	}

	objects := []ipso.Object{
		{ObjectID: 3303, ObjectVersion: "1.10", Name: "Temperature", Resources: item(5700, "Sensor Value")},
		{ObjectID: 9, ObjectVersion: "1.0", Name: "Other", Resources: item(2, "Missing")},
		{ObjectID: 3303, ObjectVersion: "1.0", Name: "TEMPERATURE", Resources: item(5701, "SENSOR VALUE")},
		{ObjectID: 3, ObjectVersion: "1.0", Name: "temperature", Resources: item(5700, "sensor value")},
		{ObjectID: 3303, ObjectVersion: "1.2", Name: "Temperature"},
	}

	return ipso.NewRegFromObjects(objects)
}

func newSmallReg() *ipso.Reg {
	objects := []ipso.Object{
		{ObjectID: 3303, ObjectVersion: ipso.DefaultObjectVersion, Name: "Temperature", Resources: ipso.ObjectResources{
			Item: []ipso.Resource{{ID: 5700, Name: "Sensor Value", Type: ipso.ResourceTypeFloat}},
		}},
		{ObjectID: 1, ObjectVersion: ipso.DefaultObjectVersion, Name: "Server"},
	}

	return ipso.NewRegFromObjects(objects)
}

type compareCase struct {
	name string
	a, b *ipso.Reg
	want []ipso.ObjectComparison
}

func compareCases() []compareCase {
	newObj := func(id int32, ver string) ipso.Object {
		return ipso.Object{Name: "obj", ObjectID: id, ObjectVersion: ver}
	}

	build := func(objects ...ipso.Object) *ipso.Reg {
		return ipso.NewRegFromObjects(objects)
	}

	return []compareCase{
		{
			name: "both empty",
			a:    build(),
			b:    build(),
		},
		{
			name: "equal",
			a:    build(newObj(1, "1.0"), newObj(2, "1.0")),
			b:    build(newObj(2, "1.0"), newObj(1, "1.0")),
		},
		{
			name: "new object",
			a:    build(newObj(1, "1.0")),
			b:    build(newObj(1, "1.0"), newObj(2, "1.0")),
			want: []ipso.ObjectComparison{
				{Difference: ipso.DifferenceTypeNewObject, ObjectComp: newObj(2, "1.0")},
			},
		},
		{
			name: "removed object",
			a:    build(newObj(1, "1.0"), newObj(2, "1.0")),
			b:    build(newObj(1, "1.0")),
			want: []ipso.ObjectComparison{
				{Difference: ipso.DifferenceTypeObjectRemoved, Object: newObj(2, "1.0")},
			},
		},
		{
			name: "new version of existing object",
			a:    build(newObj(1, "1.0")),
			b:    build(newObj(1, "1.0"), newObj(1, "1.1")),
			want: []ipso.ObjectComparison{
				{Difference: ipso.DifferenceTypeNewObject, ObjectComp: newObj(1, "1.1")},
			},
		},
		{
			name: "new listed before removed",
			a:    build(newObj(1, "1.0")),
			b:    build(newObj(2, "1.0")),
			want: []ipso.ObjectComparison{
				{Difference: ipso.DifferenceTypeNewObject, ObjectComp: newObj(2, "1.0")},
				{Difference: ipso.DifferenceTypeObjectRemoved, Object: newObj(1, "1.0")},
			},
		},
	}
}

const aliasOrigName = "orig"

func newAliasReg() *ipso.Reg {
	objects := []ipso.Object{{
		Name:          "obj",
		ObjectID:      7,
		ObjectVersion: "1.0",
		Resources:     ipso.ObjectResources{Item: []ipso.Resource{{ID: 1, Name: aliasOrigName}}},
	}}

	return ipso.NewRegFromObjects(objects)
}

// objectGetters lists every way to obtain an object from the registry built by newAliasReg.
func objectGetters() map[string]func(t *testing.T, reg *ipso.Reg) ipso.Object {
	return map[string]func(t *testing.T, reg *ipso.Reg) ipso.Object{
		"Find": func(t *testing.T, reg *ipso.Reg) ipso.Object {
			t.Helper()

			obj, err := reg.Find(ipso.Object{ObjectID: 7, ObjectVersion: "1.0"})
			require.NoError(t, err)

			return obj
		},
		"FindObjectByIDAndVer": func(t *testing.T, reg *ipso.Reg) ipso.Object {
			t.Helper()

			obj, err := reg.FindObjectByIDAndVer(7, "1.0")
			require.NoError(t, err)

			return obj
		},
		"FindObjectsByID": func(t *testing.T, reg *ipso.Reg) ipso.Object {
			t.Helper()

			objs, err := reg.FindObjectsByID(7)
			require.NoError(t, err)

			return objs[0]
		},
		"FindObjectsByName": func(t *testing.T, reg *ipso.Reg) ipso.Object {
			t.Helper()

			objs, err := reg.FindObjectsByName("obj")
			require.NoError(t, err)

			return objs[0]
		},
		"GetObjects": func(_ *testing.T, reg *ipso.Reg) ipso.Object {
			return reg.GetObjects()[0]
		},
		"Compare new": func(_ *testing.T, reg *ipso.Reg) ipso.Object {
			return (&ipso.Reg{}).Compare(reg)[0].ObjectComp
		},
		"Compare removed": func(_ *testing.T, reg *ipso.Reg) ipso.Object {
			return reg.Compare(&ipso.Reg{})[0].Object
		},
	}
}

// knownValues holds the enum values accepted in the shipped registry data.
type knownValues struct {
	instances  map[ipso.InstanceType]bool
	mandatory  map[ipso.MandatoryType]bool
	operations map[ipso.OperationType]bool
	types      map[ipso.ResourceType]bool
}

func newKnownValues() knownValues {
	return knownValues{
		instances: map[ipso.InstanceType]bool{ipso.InstanceTypeSingle: true, ipso.InstanceTypeMultiple: true},
		mandatory: map[ipso.MandatoryType]bool{ipso.MandatoryTypeOptional: true, ipso.MandatoryTypeMandatory: true},
		operations: map[ipso.OperationType]bool{
			ipso.OperationTypeRead:      true,
			ipso.OperationTypeWrite:     true,
			ipso.OperationTypeReadWrite: true,
			ipso.OperationTypeExecute:   true,
			ipso.OperationTypeEmpty:     true,
		},
		types: map[ipso.ResourceType]bool{
			ipso.ResourceTypeString:          true,
			ipso.ResourceTypeInteger:         true,
			ipso.ResourceTypeUnsignedInteger: true,
			ipso.ResourceTypeFloat:           true,
			ipso.ResourceTypeBoolean:         true,
			ipso.ResourceTypeOpaque:          true,
			ipso.ResourceTypeTime:            true,
			ipso.ResourceTypeObjLink:         true,
			ipso.ResourceTypeCoreLink:        true,
			ipso.ResourceTypeEmpty:           true,
		},
	}
}

func (k knownValues) assertResource(t *testing.T, where string, res ipso.Resource) {
	t.Helper()

	require.NotEmpty(t, res.Name, "%s has no name", where)
	require.True(t, k.instances[res.MultipleInstances], "%s: unknown instance type %q", where, res.MultipleInstances)
	require.True(t, k.mandatory[res.Mandatory], "%s: unknown mandatory type %q", where, res.Mandatory)
	require.True(t, k.operations[res.Operations], "%s: unknown operations %q", where, res.Operations)
	require.True(t, k.types[res.Type], "%s: unknown type %q", where, res.Type)
}
