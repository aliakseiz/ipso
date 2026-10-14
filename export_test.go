package ipso

// Test seams: unexported symbols the external tests (package ipso_test) need.
const (
	DefaultBaseURLForTest = defaultBaseURL
	ObjectsPathForTest    = objectsPath
)

var (
	ParseOIRString  = parseOIRString
	SanitizeText    = sanitizeText
	CompareVersions = compareVersions
)

// NewRegFromObjects returns a registry holding the objects with a built lookup index.
func NewRegFromObjects(objects []Object) *Reg {
	return &Reg{Objects: objects, objIDVerMap: objToMap(objects)}
}

// Lookup returns the registry object with the ID and version, or nil.
func (r *Reg) Lookup(id int32, ver string) *Object {
	return r.lookup(id, ver)
}

// CloneObject returns a copy of the object that shares no memory with the original.
func CloneObject(o Object) Object {
	return o.clone()
}
