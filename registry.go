package ipso

import (
	"cmp"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Sentinel errors returned by the registry, wrapped with context where available. Match them with errors.Is.
var (
	// ErrEmptyFilename is returned by ImportFile and ExportFile when the file name is empty.
	ErrEmptyFilename = errors.New("filename is empty")
	// ErrObjectNotFound is returned when no object matches the requested ID and version.
	ErrObjectNotFound = errors.New("object not found")
	// ErrResourceNotFound is returned when no resource matches the requested ID.
	ErrResourceNotFound = errors.New("resource not found")
	// ErrInvalidOIR is returned when an OIR string is not of the form object/instance/resource with
	// non-negative int32 segments.
	ErrInvalidOIR = errors.New("invalid OIR")
	// ErrNoObjects is returned by ImportFromAPI when the OMA API yielded no objects.
	ErrNoObjects = errors.New("no objects fetched")
	// ErrHTTPStatus is returned when the OMA API answers with a non-2xx status.
	ErrHTTPStatus = errors.New("unexpected HTTP status")
	// ErrBodyTooLarge is returned when an OMA API response exceeds the size limit.
	ErrBodyTooLarge = errors.New("response body too large")
)

const (
	// defaultBaseURL is the OMA API root, used when Configuration.BaseURL is empty.
	defaultBaseURL = "https://www.openmobilealliance.org/api/lwm2m/v1"
	// objectsPath is the OMA API path listing all registered objects, relative to the base URL.
	objectsPath = "/Object"
	// oirSegments is the number of segments in an OIR string: object/instance/resource.
	oirSegments = 3
	// requestTimeout limits a single OMA API request.
	requestTimeout = 5 * time.Second
	// maxBodySize limits a single OMA API response body.
	maxBodySize = 16 << 20
	// exportFileMode is the permission of files created by ExportFile.
	exportFileMode = 0o600
)

// defaultClient is shared by all registries that have no Configuration.HTTPClient, http.Client is safe for concurrent use.
//
//nolint:gochecknoglobals // one shared client lets all registries reuse connections
var defaultClient = &http.Client{Timeout: requestTimeout}

// Registry is the interface of the core registry functions, implemented by *Reg.
// See the methods of Reg for the behavior of each function; Reg has further methods that are not part of the interface.
type Registry interface {
	Import(r io.Reader) error
	ImportFile(filename string) error
	ImportFromAPI() ([]Object, error)
	Export(w io.Writer) error
	ExportFile(filename string) error
	Compare(reg *Reg) []ObjectComparison

	Find(o Object) (Object, error)
	FindObjectsByID(id int32) ([]Object, error)
	FindObjectByIDAndVer(id int32, ver string) (Object, error)
	FindResourcesByID(id int32) ([]Resource, error)
	FindResourcesByObjResIDs(objID, resID int32) ([]Resource, error)
	FindResourceByObjIDObjVerResID(objID int32, objVer string, resID int32) (Resource, error)
	FindResourceByOIR(oir string) (Resource, error)

	GetObjects() []Object

	Sanitize(sanitizer []string)

	GetRegistry() *Reg
}

// Reg is a registry of OMA objects and their resources.
// The zero value is an empty registry; New and NewContext create one from a Configuration.
type Reg struct {
	// Config is the configuration the registry was created with.
	Config Configuration
	// Objects holds the registry content. Lookups use an index built by New, NewContext and Import,
	// so replace the slice through Import rather than by assignment.
	Objects []Object
	// objIDVerMap map contains references to the objects stored in Objects slice.
	objIDVerMap map[int32]map[string]*Object // First key - ObjectID, second - Version
}

// Compile-time check that Reg implements Registry.
var _ Registry = (*Reg)(nil)

// New creates a registry from the configuration; use DefaultConfiguration for the default one.
// With Configuration.InitOnNew the objects are fetched from the OMA API, which cannot be canceled, use NewContext for that.
// With Configuration.Sanitize the registry is sanitized after loading.
func New(cfg Configuration) (*Reg, error) {
	return NewContext(context.Background(), cfg)
}

// NewContext creates a registry like New. The context governs fetching objects from the OMA API
// (Configuration.InitOnNew): canceling it aborts the fetch and the context error is returned.
func NewContext(ctx context.Context, cfg Configuration) (*Reg, error) {
	reg := Reg{Config: cfg}

	if reg.Config.InitOnNew {
		var err error

		reg.Objects, err = reg.ImportFromAPIContext(ctx)
		if err != nil {
			return nil, err
		}

		reg.objIDVerMap = objToMap(reg.Objects)
	}

	if reg.Config.Sanitize {
		sanitizer := reg.Config.Sanitizer
		if len(sanitizer) == 0 {
			sanitizer = DefaultSanitizer()
		}

		reg.Sanitize(sanitizer)
	}

	return &reg, nil
}

// Import loads objects and resources from the YAML read from reader, the format written by Export.
// It replaces the current registry content and rebuilds the lookup index.
// The registry is not modified when an error is returned.
func (r *Reg) Import(reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("read registry: %w", err)
	}

	var objects []Object

	err = yaml.Unmarshal(data, &objects)
	if err != nil {
		return fmt.Errorf("parse registry YAML: %w", err)
	}

	r.Objects = objects
	r.objIDVerMap = objToMap(r.Objects)

	return nil
}

// ImportFile loads objects and resources from the YAML file, see Import.
// It returns ErrEmptyFilename when the file name is empty.
func (r *Reg) ImportFile(filename string) error {
	if filename == "" {
		return ErrEmptyFilename
	}

	file, err := os.Open(filename) //nolint:gosec // the caller chooses the file
	if err != nil {
		return fmt.Errorf("import file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	return r.Import(file)
}

// ImportFromAPI fetches objects and resources from the OMA API and returns them, the registry is not modified.
// It is ImportFromAPIContext with a background context.
func (r *Reg) ImportFromAPI() ([]Object, error) {
	return r.ImportFromAPIContext(context.Background())
}

// ImportFromAPIContext fetches objects and resources from the OMA API and returns them, the registry is not modified.
// Empty object and LwM2M versions are filled with the default "1.0".
// Requests use Configuration.HTTPClient and Configuration.BaseURL, when set, and run one after another.
// It returns the context error when the context is canceled or its deadline passes, even with SkipInitErrors set,
// and ErrNoObjects when no objects were fetched, for example when all requests failed and SkipInitErrors is set.
func (r *Reg) ImportFromAPIContext(ctx context.Context) ([]Object, error) {
	objectsMeta, err := r.getObjectsMeta(ctx)
	if err != nil {
		return nil, err
	}

	var objects []Object

	for _, objectMeta := range objectsMeta {
		object, err := r.getObject(ctx, objectMeta)
		if err != nil {
			ctxErr := ctx.Err()
			if ctxErr != nil {
				return nil, fmt.Errorf("import from API: %w", ctxErr)
			}

			if r.Config.SkipInitErrors {
				continue
			}

			return nil, err
		}

		if object.ObjectVersion == "" {
			object.ObjectVersion = DefaultObjectVersion
		}

		if object.LWM2MVersion == "" {
			object.LWM2MVersion = DefaultLwM2MVersion
		}

		objects = append(objects, object)
	}

	if len(objects) == 0 {
		return nil, ErrNoObjects
	}

	return objects, nil
}

// Export writes registry objects and resources to w in YAML format.
// Leading and trailing spaces are trimmed from text fields in the exported copy, the registry is not modified.
// The YAML encoder writes multi-line text starting with a whitespace as a block scalar,
// which cannot be parsed back by Import.
func (r *Reg) Export(w io.Writer) error {
	data, err := yaml.Marshal(trimObjects(r.Objects))
	if err != nil {
		return fmt.Errorf("marshal registry: %w", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("write registry: %w", err)
	}

	return nil
}

// ExportFile writes registry objects and resources to the file in YAML format, see Export.
// An existing file is truncated, a new one is created with mode 0600.
// It returns ErrEmptyFilename when the file name is empty.
func (r *Reg) ExportFile(filename string) error {
	if filename == "" {
		return ErrEmptyFilename
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, exportFileMode) //nolint:gosec // the caller chooses the file
	if err != nil {
		return fmt.Errorf("export file: %w", err)
	}

	err = r.Export(file)
	if err != nil {
		_ = file.Close()

		return err
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("export file: %w", err)
	}

	return nil
}

// Compare compares r with reg, the objects are matched by ID and version.
// It returns the objects present in only one of the registries with a difference explanation:
// objects present only in reg are listed first, followed by objects present only in r.
// The content of objects present in both registries is not compared.
func (r *Reg) Compare(reg *Reg) []ObjectComparison {
	var objComp []ObjectComparison

	for _, regObj := range reg.Objects {
		if r.lookup(regObj.ObjectID, regObj.ObjectVersion) == nil {
			objComp = append(objComp, ObjectComparison{Difference: DifferenceTypeNewObject, Object: Object{}, ObjectComp: regObj.clone()})
		}
	}

	for _, rObj := range r.Objects {
		if reg.lookup(rObj.ObjectID, rObj.ObjectVersion) == nil {
			objComp = append(objComp, ObjectComparison{Difference: DifferenceTypeObjectRemoved, Object: rObj.clone(), ObjectComp: Object{}})
		}
	}

	return objComp
}

// Find looks for the object with the ID and version of o in the registry.
// The returned object is a copy, changing it does not affect the registry.
// It returns an empty object and ErrObjectNotFound when the object is not found.
func (r *Reg) Find(o Object) (Object, error) {
	return r.FindObjectByIDAndVer(o.ObjectID, o.ObjectVersion)
}

// FindObjectsByID finds objects in the registry by ID.
// Several versions of an object can be returned, ordered by version ascending.
// The returned objects are copies, changing them does not affect the registry.
// It returns ErrObjectNotFound when no objects are found.
func (r *Reg) FindObjectsByID(id int32) ([]Object, error) {
	found, err := r.sortedObjectsByID(id)
	if err != nil {
		return nil, err
	}

	objects := make([]Object, 0, len(found))
	for _, object := range found {
		objects = append(objects, object.clone())
	}

	return objects, nil
}

// FindObjectByIDAndVer finds an object in the registry by ID and version.
// The returned object is a copy, changing it does not affect the registry.
// It returns an empty object and ErrObjectNotFound when the object is not found.
func (r *Reg) FindObjectByIDAndVer(id int32, ver string) (Object, error) {
	obj := r.lookup(id, ver)
	if obj == nil {
		return Object{}, ErrObjectNotFound
	}

	return obj.clone(), nil
}

// FindObjectsByName finds objects in the registry by name, ignoring letter case.
// Several objects can share a name: versions of one object, or different objects ("Location" is object 6 and 3336).
// They are ordered by ID, then by version ascending.
// The returned objects are copies, changing them does not affect the registry.
// An empty name matches nothing. It returns ErrObjectNotFound when no objects are found.
func (r *Reg) FindObjectsByName(name string) ([]Object, error) {
	if name == "" {
		return nil, ErrObjectNotFound
	}

	var found []*Object

	for i := range r.Objects {
		if strings.EqualFold(r.Objects[i].Name, name) {
			found = append(found, &r.Objects[i])
		}
	}

	if len(found) == 0 {
		return nil, ErrObjectNotFound
	}

	slices.SortFunc(found, compareObjects)

	objects := make([]Object, 0, len(found))
	for _, object := range found {
		objects = append(objects, object.clone())
	}

	return objects, nil
}

// FindResourcesByID finds resources in the registry by ID.
// It returns the matching resources from all objects of all versions, in the order of Reg.Objects.
// It returns ErrResourceNotFound when no resource is found.
func (r *Reg) FindResourcesByID(id int32) ([]Resource, error) {
	var resources []Resource

	for _, rObj := range r.Objects {
		for _, rRes := range rObj.Resources.Item {
			if rRes.ID == id {
				resources = append(resources, rRes)
			}
		}
	}

	if len(resources) == 0 {
		return nil, ErrResourceNotFound
	}

	return resources, nil
}

// FindResourcesByName finds resources in the registry by name, ignoring letter case.
// It returns the matching resources from all objects of all versions, ordered by object ID, then by object version ascending.
// An empty name matches nothing. It returns ErrResourceNotFound when no resource is found.
func (r *Reg) FindResourcesByName(name string) ([]Resource, error) {
	if name == "" {
		return nil, ErrResourceNotFound
	}

	var matched []*Object

	for i := range r.Objects {
		for _, res := range r.Objects[i].Resources.Item {
			if strings.EqualFold(res.Name, name) {
				matched = append(matched, &r.Objects[i])

				break
			}
		}
	}

	slices.SortFunc(matched, compareObjects)

	var resources []Resource

	for _, obj := range matched {
		for _, res := range obj.Resources.Item {
			if strings.EqualFold(res.Name, name) {
				resources = append(resources, res)
			}
		}
	}

	if len(resources) == 0 {
		return nil, ErrResourceNotFound
	}

	return resources, nil
}

// FindResourcesByObjResIDs finds resources in the registry by object ID and resource ID.
// It returns the matching resources from all versions of the object, ordered by object version ascending.
// It returns ErrObjectNotFound or ErrResourceNotFound when the object or the resource is not found.
func (r *Reg) FindResourcesByObjResIDs(objID, resID int32) ([]Resource, error) {
	var resources []Resource

	objects, err := r.sortedObjectsByID(objID)
	if err != nil {
		return nil, err
	}

	for _, obj := range objects {
		for _, rRes := range obj.Resources.Item {
			if rRes.ID == resID {
				resources = append(resources, rRes)
			}
		}
	}

	if len(resources) == 0 {
		return nil, ErrResourceNotFound
	}

	return resources, nil
}

// FindResourceByObjIDObjVerResID finds a resource in the registry by object ID, object version and resource ID.
// It returns an empty resource and ErrObjectNotFound or ErrResourceNotFound when the object or the resource is not found.
func (r *Reg) FindResourceByObjIDObjVerResID(objID int32, objVer string, resID int32) (Resource, error) {
	obj := r.lookup(objID, objVer)
	if obj == nil {
		return Resource{}, ErrObjectNotFound
	}

	return obj.resource(resID)
}

// FindResourceByOIR finds a resource in the registry by an OIR string, object/instance/resource, i.e. "3303/0/5700".
// One leading slash is accepted. The segments must be non-negative int32 numbers, otherwise the error wraps ErrInvalidOIR.
// The instance segment is ignored and the object version is always DefaultObjectVersion ("1.0"),
// use FindResourceByOIRAndVer for another version.
// It returns an empty resource and ErrObjectNotFound or ErrResourceNotFound when the object or the resource is not found.
func (r *Reg) FindResourceByOIR(oir string) (Resource, error) {
	return r.FindResourceByOIRAndVer(oir, DefaultObjectVersion)
}

// FindResourceByOIRAndVer finds a resource by an OIR string like FindResourceByOIR, in the object version ver.
// It returns an empty resource and ErrObjectNotFound or ErrResourceNotFound when the object version
// or the resource is not found.
func (r *Reg) FindResourceByOIRAndVer(oir, ver string) (Resource, error) {
	objID, _, resID, err := parseOIRString(oir)
	if err != nil {
		return Resource{}, fmt.Errorf("cannot parse OIR (%s): %w", oir, err)
	}

	obj := r.lookup(objID, ver)
	if obj == nil {
		return Resource{}, fmt.Errorf("cannot find object (%d ver.%s) %w", objID, ver, ErrObjectNotFound)
	}

	return obj.resource(resID)
}

// GetObjects returns all objects of the registry.
// The returned objects are copies, changing them does not affect the registry.
func (r *Reg) GetObjects() []Object {
	if r.Objects == nil {
		return nil
	}

	objects := make([]Object, 0, len(r.Objects))
	for _, object := range r.Objects {
		objects = append(objects, object.clone())
	}

	return objects
}

// Sanitize cleans the description fields of objects and resources in place using sanitizer strings.
// A sanitizer string is removed from the text; a whitespace-only one ("\n", "\t", "  ") is replaced by a single space,
// so words and sentences separated by it stay apart. When a whitespace-only string is present,
// runs of spaces are collapsed. Leading and trailing spaces are always removed.
// Description fields in objects and resources do not follow any single format or convention
// in regard to line breaks, lists presentation, special characters escaping etc.
// thus in some cases cannot be used directly in external applications (i.e. properly displayed in browser).
func (r *Reg) Sanitize(sanitizer []string) {
	for i := range r.Objects {
		object := &r.Objects[i] // Modify the object in registry instead of a copy
		object.Description1 = sanitizeText(object.Description1, sanitizer)
		object.Description2 = sanitizeText(object.Description2, sanitizer)

		for j := range object.Resources.Item {
			resource := &object.Resources.Item[j]
			resource.Description = sanitizeText(resource.Description, sanitizer)
			resource.RangeEnumeration = sanitizeText(resource.RangeEnumeration, sanitizer)
			resource.Units = sanitizeText(resource.Units, sanitizer)
		}
	}
}

// GetRegistry returns the registry itself, not a copy. It lets Compare take another registry through the Registry interface.
func (r *Reg) GetRegistry() *Reg {
	return r
}

// sortedObjectsByID return registry objects with the ID ordered by version ascending.
// The pointers refer to registry data and must not be modified or returned to callers.
func (r *Reg) sortedObjectsByID(id int32) ([]*Object, error) {
	objVerMap, ok := r.objIDVerMap[id]
	if !ok {
		return nil, ErrObjectNotFound
	}

	objects := make([]*Object, 0, len(objVerMap))
	for _, object := range objVerMap {
		objects = append(objects, object)
	}

	slices.SortFunc(objects, compareObjects)

	return objects, nil
}

// compareObjects order objects by ID, then by version numerically.
func compareObjects(left, right *Object) int {
	if c := cmp.Compare(left.ObjectID, right.ObjectID); c != 0 {
		return c
	}

	if c := compareVersions(left.ObjectVersion, right.ObjectVersion); c != 0 {
		return c
	}

	// Versions equal by value ("1.0", "1.00") still get a fixed order.
	return strings.Compare(left.ObjectVersion, right.ObjectVersion)
}

// lookup return the registry object with the ID and version, or nil when there is none.
// The pointer refers to registry data and must not be modified or returned to callers.
func (r *Reg) lookup(id int32, ver string) *Object {
	return r.objIDVerMap[id][ver]
}

// getObjectsMeta retrieve all objects metadata.
func (r *Reg) getObjectsMeta(ctx context.Context) ([]ObjectMeta, error) {
	baseURL := r.Config.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	body, err := r.getURL(ctx, strings.TrimSuffix(baseURL, "/")+objectsPath)
	if err != nil {
		return nil, err
	}

	var objectMetas []ObjectMeta

	err = json.Unmarshal(body, &objectMetas)
	if err != nil {
		return nil, fmt.Errorf("parse objects list: %w", err)
	}

	return objectMetas, nil
}

// getObject fetch object details based on metadata.
func (r *Reg) getObject(ctx context.Context, objectMeta ObjectMeta) (Object, error) {
	body, err := r.getURL(ctx, objectMeta.ObjectLink)
	if err != nil {
		return Object{}, err
	}

	var lwm2m Lwm2M

	err = xml.Unmarshal(body, &lwm2m)
	if err != nil {
		return Object{}, fmt.Errorf("parse object %s: %w", objectMeta.ObjectLink, err)
	}

	return lwm2m.Object, nil
}

// trimObjects returns a deep copy of objects with leading and trailing spaces removed from text fields.
func trimObjects(objects []Object) []Object {
	trimmed := make([]Object, len(objects))

	for i, object := range objects {
		object.Name = strings.TrimSpace(object.Name)
		object.Description1 = strings.TrimSpace(object.Description1)
		object.Description2 = strings.TrimSpace(object.Description2)

		if object.Resources.Item != nil {
			items := make([]Resource, len(object.Resources.Item))

			for j, resource := range object.Resources.Item {
				resource.Name = strings.TrimSpace(resource.Name)
				resource.Description = strings.TrimSpace(resource.Description)
				resource.RangeEnumeration = strings.TrimSpace(resource.RangeEnumeration)
				resource.Units = strings.TrimSpace(resource.Units)
				items[j] = resource
			}

			object.Resources.Item = items
		}

		trimmed[i] = object
	}

	return trimmed
}

func objToMap(objects []Object) map[int32]map[string]*Object {
	objMap := make(map[int32]map[string]*Object)

	for i, object := range objects {
		if _, ok := objMap[object.ObjectID]; !ok {
			objMap[object.ObjectID] = make(map[string]*Object)
		}

		objMap[object.ObjectID][object.ObjectVersion] = &objects[i]
	}

	return objMap
}

// getURL fetch the URL body. Non-2xx responses and bodies over maxBodySize are errors.
func (r *Reg) getURL(ctx context.Context, url string) ([]byte, error) {
	client := r.Config.HTTPClient
	if client == nil {
		client = defaultClient
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("GET %s: %w: %s", url, ErrHTTPStatus, resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize+1))
	if err != nil {
		return nil, fmt.Errorf("GET %s: read body: %w", url, err)
	}

	if len(body) > maxBodySize {
		return nil, fmt.Errorf("GET %s: %w", url, ErrBodyTooLarge)
	}

	return body, nil
}

// sanitizeText applies the sanitizer strings to the text in order, see Sanitize.
func sanitizeText(text string, sanitizer []string) string {
	spaced := false

	for _, remove := range sanitizer {
		if remove == "" {
			continue
		}

		if strings.TrimSpace(remove) == "" {
			text = strings.ReplaceAll(text, remove, " ")
			spaced = true

			continue
		}

		text = strings.ReplaceAll(text, remove, "")
	}

	if spaced {
		text = collapseSpaces(text)
	}

	return strings.TrimSpace(text)
}

// collapseSpaces replaces every run of spaces with a single space.
func collapseSpaces(text string) string {
	var builder strings.Builder

	builder.Grow(len(text))

	prevSpace := false

	for i := 0; i < len(text); i++ {
		isSpace := text[i] == ' '
		if isSpace && prevSpace {
			continue
		}

		builder.WriteByte(text[i])

		prevSpace = isSpace
	}

	return builder.String()
}

// compareVersions compare dot-separated versions ("1.2", "1.10") segment by segment and return -1, 0 or 1.
// Numeric segments compare as numbers, so "1.2" < "1.10". An empty segment sorts before a number,
// a number before any other text, other text compares as a string.
// A version that is a prefix of another one sorts first.
func compareVersions(left, right string) int {
	leftParts, rightParts := strings.Split(left, "."), strings.Split(right, ".")

	for i := 0; i < len(leftParts) && i < len(rightParts); i++ {
		if c := compareVersionSegments(leftParts[i], rightParts[i]); c != 0 {
			return c
		}
	}

	return cmp.Compare(len(leftParts), len(rightParts))
}

// compareVersionSegments compare a single version segment, see compareVersions.
func compareVersionSegments(left, right string) int {
	leftNum, leftErr := strconv.ParseUint(left, 10, 64)
	rightNum, rightErr := strconv.ParseUint(right, 10, 64)

	switch {
	case leftErr == nil && rightErr == nil:
		return cmp.Compare(leftNum, rightNum)
	case left == "" || right == "":
		return cmp.Compare(len(left), len(right))
	case leftErr == nil:
		return -1
	case rightErr == nil:
		return 1
	default:
		return strings.Compare(left, right)
	}
}

// parseOIRString parse OIR ("3303/0/5700" or "/3303/0/5700" - object/instance/resource) string.
// Every segment must be a non-negative decimal number that fits int32.
// All errors wrap ErrInvalidOIR. The results are the object ID, instance number and resource ID.
func parseOIRString(oir string) (int32, int32, int32, error) {
	ids := strings.Split(strings.TrimPrefix(oir, "/"), "/")
	if len(ids) != oirSegments {
		return 0, 0, 0, fmt.Errorf("%w: want object/instance/resource", ErrInvalidOIR)
	}

	parsed := [oirSegments]int32{}

	for i, id := range ids {
		value, err := parseOIRSegment(id)
		if err != nil {
			return 0, 0, 0, err
		}

		parsed[i] = value
	}

	return parsed[0], parsed[1], parsed[2], nil
}

// parseOIRSegment parse a single non-negative int32 OIR segment.
func parseOIRSegment(segment string) (int32, error) {
	// Bit size 31 limits the value to math.MaxInt32, so the conversion to int32 is lossless.
	n, err := strconv.ParseUint(segment, 10, 31)
	if err != nil {
		return 0, fmt.Errorf("%w: segment %q: %w", ErrInvalidOIR, segment, err)
	}

	return int32(n), nil
}
