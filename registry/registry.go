// Package registry provides methods to create and control OMA IPSO registry.
package registry

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	errEmptyFilename = errors.New("filename is empty")
	errObjNotFound   = errors.New("object not found")
	errResNotFound   = errors.New("resource not found")
)

// Registry an interface defining public library functions.
type Registry interface {
	Import(filename string) error
	ImportFromAPI() ([]Object, error)
	Export(filename string) error
	Compare(reg *Reg) []ObjectComparison

	Find(o Object) (Object, error)
	FindObjectsByID(id int32) ([]Object, error)
	FindObjectByIDAndVer(id int32, ver string) (Object, error)
	FindResourcesByID(id int32) ([]Resource, error)
	FindResourcesByObjResIDs(objID, resID int32) ([]Resource, error)
	FindResourceByObjIDObjVerResID(objID int32, objVer string, resID int32) (Resource, error)

	Sanitize(sanitizer []string)

	GetRegistry() *Reg
}

// Reg contains registry objects and settings.
type Reg struct {
	Config  Configuration
	Objects []Object
	// objIDVerMap map contains references to objects stored in Objects slice.
	objIDVerMap map[int32]map[string]*Object // First key - ObjectID, second - Version
}

// TODO implement tests for registry
// Export
// Import
// ImportFromAPI
// Compare
// Sanitize
// Find

// New creates a new registry, using provided or default configuration.
func New(cfg Configuration) (Registry, error) {
	reg := Reg{
		Config:  cfg,
		Objects: nil,
	}

	if reg.Config.InitOnNew {
		var err error

		reg.Objects, err = reg.ImportFromAPI()
		if err != nil {
			return nil, err
		}

		reg.objIDVerMap = objToMap(reg.Objects)
	}

	if reg.Config.Sanitize {
		reg.Sanitize(DefaultSanitizer())
	}

	return &reg, nil
}

// Import loads objects and resources from file.
// Overwrites current registry objects and resources.
func (r *Reg) Import(filename string) error {
	if filename == "" {
		return errEmptyFilename
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &r.Objects); err != nil {
		return err
	}

	r.objIDVerMap = objToMap(r.Objects)

	return nil
}

// ImportFromAPI initializes the registry from official OMA API.
// Overwrites current registry objects and resources.
// Fills empty object and LwM2M versions with default values "1.0".
func (r *Reg) ImportFromAPI() ([]Object, error) {
	objectsMeta, err := r.getObjectsMeta()
	if err != nil {
		return nil, err
	}

	// nolint:prealloc
	var objects []Object

	for _, objectMeta := range objectsMeta {
		object, err := r.getObject(objectMeta)
		if err != nil {
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

	return objects, nil
}

// Export stores registry objects and resources in a specified file in YAML format.
func (r *Reg) Export(filename string) error {
	if filename == "" {
		return errEmptyFilename
	}

	data, err := yaml.Marshal(&r.Objects)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0o600)
}

// Compare makes comparison of r and reg registries.
// Returns a list of non-equal objects with difference explanation.
func (r *Reg) Compare(reg *Reg) []ObjectComparison {
	var objComp []ObjectComparison

	// Compare r with reg
	for _, regObj := range reg.Objects {
		if rObj, err := r.Find(regObj); err != nil {
			regObjCopy := regObj

			switch err {
			case errObjNotFound:
				objComp = append(objComp, ObjectComparison{
					Difference: DifferenceTypeNewObject,
					Object:     Object{},
					ObjectComp: regObjCopy,
				})
			default:
				objComp = append(objComp, ObjectComparison{
					Difference: DifferenceTypeUnknown,
					Object:     rObj,
					ObjectComp: regObjCopy,
				})
			}
		}
	}

	// Compare reg with r
	for _, rObj := range r.Objects {
		if regObj, err := reg.Find(rObj); err != nil {
			// regObjCopy := regObj
			rObjCopy := rObj

			switch err {
			case errObjNotFound:
				objComp = append(objComp, ObjectComparison{
					Difference: DifferenceTypeObjectRemoved,
					Object:     rObjCopy,
					ObjectComp: Object{},
				})
			default:
				objComp = append(objComp, ObjectComparison{
					Difference: DifferenceTypeUnknown,
					Object:     rObjCopy,
					ObjectComp: regObj,
				})
			}
		}
	}

	return objComp
}

// Find looks for an object in the current registry.
// Returns an empty object and error, when the object not found.
func (r *Reg) Find(o Object) (Object, error) {
	if objVerMap, ok := r.objIDVerMap[o.ObjectID]; ok {
		if obj, ok := objVerMap[o.ObjectVersion]; ok {
			return *obj, nil
		}
	}

	return Object{}, errObjNotFound
}

// FindObjectsByID finds objects in the registry by ID.
// Multiple objects with same ID and different versions could be returned.
// Returns an error, when no objects found.
func (r *Reg) FindObjectsByID(id int32) ([]Object, error) {
	var objects []Object

	if objVerMap, ok := r.objIDVerMap[id]; ok {
		for _, object := range objVerMap {
			objects = append(objects, *object)
		}
	} else {
		return nil, errObjNotFound
	}

	return objects, nil
}

// TODO implement `FindObjectByName`, `FindResourceByName`, `FindObjectByDescription`, `FindResourceByDescription`

// FindObjectByIDAndVer finds an object in the registry by ID and version.
// Returns an empty object and error, when the object not found.
func (r *Reg) FindObjectByIDAndVer(id int32, ver string) (Object, error) {
	if objVerMap, ok := r.objIDVerMap[id]; ok {
		if obj, ok := objVerMap[ver]; ok {
			return *obj, nil
		}
	}

	return Object{}, errObjNotFound
}

// FindResourcesByID finds resources in registry by ID.
// Returns matching resources from all objects of all versions.
// Returns an error, when resource not found.
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
		return nil, errResNotFound
	}

	return resources, nil
}

// FindResourcesByObjResIDs finds resources in registry by object ID and resource ID.
// Returns matching resources from all versions of specific object.
// Returns an error, when resource or object not found.
func (r *Reg) FindResourcesByObjResIDs(objID, resID int32) ([]Resource, error) {
	var resources []Resource

	objVerMap, err := r.FindObjectsByID(objID)
	if err != nil {
		return nil, err
	}

	for _, obj := range objVerMap {
		for _, rRes := range obj.Resources.Item {
			if rRes.ID == resID {
				resources = append(resources, rRes)
			}
		}
	}

	if len(resources) == 0 {
		return nil, errResNotFound
	}

	return resources, nil
}

// FindResourceByObjIDObjVerResID finds specific resource in registry by object ID, object version and resource ID.
// Returns a single matching resource of specific object.
// Returns an empty resource and error, when resource or object not found.
func (r *Reg) FindResourceByObjIDObjVerResID(objID int32, objVer string, resID int32) (Resource, error) {
	obj, err := r.FindObjectByIDAndVer(objID, objVer)
	if err != nil {
		return Resource{}, err
	}

	for _, res := range obj.Resources.Item {
		if res.ID == resID {
			return res, nil
		}
	}

	return Resource{}, errResNotFound
}

// Sanitize removes unwanted strings from objects and resources description fields
// using sanitizer strings from registry configuration. Also removes leading and trailing spaces.
// Description fields in objects and resources do not follow any single format or convention
// with regards to line breaks, lists presentation, special characters escaping etc.
// thus in some cases cannot be used directly in external applications (i.e. properly displayed in browser).
func (r *Reg) Sanitize(sanitizer []string) {
	// TODO implement sanitization using regular expressions
	// TODO run in parallel goroutines to speed it up
	for _, s := range sanitizer {
		for oIndex := 0; oIndex < len(r.Objects); oIndex++ {
			object := r.Objects[oIndex] // Modify the object in registry instead of object's copy
			object.Description1 = strings.ReplaceAll(object.Description1, s, "")
			object.Description1 = strings.TrimSpace(object.Description1)
			object.Description2 = strings.ReplaceAll(object.Description2, s, "")
			object.Description2 = strings.TrimSpace(object.Description2)

			for rIndex := 0; rIndex < len(r.Objects[oIndex].Resources.Item); rIndex++ {
				resource := &r.Objects[oIndex].Resources.Item[rIndex]
				resource.Description = strings.ReplaceAll(resource.Description, s, "")
				resource.Description = strings.TrimSpace(resource.Description)
				resource.RangeEnumeration = strings.ReplaceAll(resource.RangeEnumeration, s, "")
				resource.RangeEnumeration = strings.TrimSpace(resource.RangeEnumeration)
				resource.Units = strings.ReplaceAll(resource.Units, s, "")
				resource.Units = strings.TrimSpace(resource.Units)
			}
		}
	}
}

// GetRegistry returns internal registry structure. Temporary solution, will be removed in next releases.
func (r *Reg) GetRegistry() *Reg {
	return r
}

// getObjectsMeta retrieve all objects metadata.
func (r *Reg) getObjectsMeta() ([]ObjectMeta, error) {
	body, err := getURL("http://www.openmobilealliance.org/api/lwm2m/v1/Object")
	if err != nil {
		return nil, err
	}

	var objectMetas []ObjectMeta

	if err := json.Unmarshal(body, &objectMetas); err != nil {
		return nil, err
	}

	return objectMetas, nil
}

// getObject fetch object details based on metadata.
func (r *Reg) getObject(objectMeta ObjectMeta) (Object, error) {
	body, err := getURL(objectMeta.ObjectLink)
	if err != nil {
		return Object{}, err
	}

	var lwm2m Lwm2M

	if err := xml.Unmarshal(body, &lwm2m); err != nil {
		return Object{}, err
	}

	return lwm2m.Object, nil
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

func getURL(url string) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second, // nolint:gomnd
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return ioutil.ReadAll(resp.Body)
}
