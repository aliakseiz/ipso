// Package registry provides methods to create and control OMA IPSO registry.
package registry

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliakseiz/ipso-registry/ipso"
	"github.com/cenkalti/backoff/v4"
	"gopkg.in/yaml.v3"
)

var (
	errEmptyFilename  = errors.New("filename is empty")
	errObjNotFound    = errors.New("object not found")
	errObjVerNotFound = errors.New("object version not found")
	errResNotFound    = errors.New("resource not found")
)

// Registry holds objects and settings.
type Registry struct {
	Config *Configuration

	Objects []*Object
	objByID map[int64]map[float64]*Object // primary key - object ID, secondary - object version
}

// TODO implement tests for registry
// Registry interface
// Export
// Import
// ImportFromAPI
// Compare
// Sanitize
// Find

// New creates a new registry, using provided or default configuration.
func New(cfg *Configuration) (*Registry, error) {
	var err error

	reg := &Registry{
		Config:  cfg,
		Objects: nil,
	}

	if reg.Config == nil {
		reg.Config = DefaultConfiguration()
	}

	if reg.Config.InitOnNew {
		reg.Objects, err = reg.ImportFromAPI()
		if err != nil {
			return nil, err
		}
	}

	if reg.Config.Sanitize {
		reg.Sanitize()
	}

	return reg, nil
}

// TODO implement sanitization using regular expressions

// Sanitize removes unwanted strings from objects and resources description fields
// using sanitizer strings from registry configuration. Also removes leading and trailing spaces.
// Description fields in objects and resources do not follow any single format or convention
// with regards to line breaks, lists presentation, special characters escaping etc.
// thus in some cases cannot be used directly in external applications (i.e. properly displayed in browser).
func (r *Registry) Sanitize() {
	// TODO run in parallel goroutines to speed it up
	for _, s := range r.Config.Sanitizer {
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

// Export stores registry objects and resources in a specified file in YAML format.
func (r *Registry) Export(filename string) error {
	filename = filepath.Clean(filename)
	if filename == "" {
		return errEmptyFilename
	}

	data, err := yaml.Marshal(&r.Objects)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0o644)
}

// Import loads objects and resources from file.
// Overwrites current registry objects and resources.
func (r *Registry) Import(filename string) error {
	filename = filepath.Clean(filename)
	if filename == "" {
		return errEmptyFilename
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &r.Objects)
}

// ImportFromAPI initializes the registry from official OMA API.
// Overwrites current registry objects and resources.
// TODO make import asynchronous, run it in separate go routine
// TODO block Find and Export operations while importing to avoid inconsistent state
func (r *Registry) ImportFromAPI() ([]*Object, error) {
	objectsMeta, err := r.getObjectsMeta()
	if err != nil {
		return nil, err
	}

	var objects []*Object

	for _, objectMeta := range objectsMeta {
		object, err := r.getObject(objectMeta)
		if err != nil {
			if r.Config.SkipInitErrors {
				continue
			}

			return nil, err
		}

		objects = append(objects, object)
	}

	return objects, nil
}

// Compare makes comparison of r and reg registries.
// Returns a list of non-equal objects with difference explanation.
func (r *Registry) Compare(reg *Registry) []*ObjectComparison {
	// TODO store objects in registry in maps to improve lookup performance
	var objComp []*ObjectComparison

	// Compare r with reg
	for _, regObj := range reg.Objects {
		if rObj, err := r.Find(regObj); err != nil {
			regObjCopy := regObj
			// rObjCopy := rObj

			switch err {
			case errObjNotFound:
				objComp = append(objComp, &ObjectComparison{
					Difference: DifferenceTypeNewObject,
					Object:     nil,
					ObjectComp: regObjCopy,
				})
			default:
				objComp = append(objComp, &ObjectComparison{
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
				objComp = append(objComp, &ObjectComparison{
					Difference: DifferenceTypeObjectRemoved,
					Object:     rObjCopy,
					ObjectComp: nil,
				})
			default:
				objComp = append(objComp, &ObjectComparison{
					Difference: DifferenceTypeUnknown,
					Object:     rObjCopy,
					ObjectComp: regObj,
				})
			}
		}
	}

	return objComp
}

// TODO implement `FindObjectByName`, `FindResourceByName`, `FindObjectByDescription`, `FindResourceByDescription`

// Find looks for an object in current registry.
func (r *Registry) Find(o *Object) (*Object, error) {
	for _, rObj := range r.Objects {
		if rObj.ObjectID == o.ObjectID && rObj.ObjectVersion == o.ObjectVersion {
			return rObj, nil
		}
	}

	return nil, errObjNotFound
}

// FindObject finds objects in registry by ID and version.
// Uses latest object version, when `ver` is 0.
func (r *Registry) FindObject(id int64, ver float64) (*Object, error) {
	objByMap, ok := r.objByID[id]
	if !ok {
		return nil, errObjNotFound
	}

	// Find the latest version of an object
	if ver == 0 {
		max := float64(-1)

		for v := range objByMap {
			if v > max || max == -1 {
				max = v
			}
		}
	}

	// Get the object by version
	obj, ok := objByMap[ver]
	if !ok {
		return nil, errObjVerNotFound
	}

	return obj, nil
}

// FindObjectByURN finds an object in registry by URN.
func (r *Registry) FindObjectByURN(urn string) (*Object, error) {
	u, err := parseURN(urn)
	if err != nil {
		return nil, err
	}

	if objByVer, ok := r.objByID[u.ObjID]; ok {
		if obj, ok := objByVer[u.Version]; ok {
			return obj, nil
		}
	}

	return nil, errObjNotFound
}

// FindObjectsByID finds objects in registry by ID.
// Multiple objects with same ID and different versions could be returned.
func (r *Registry) FindObjectsByID(id int64) ([]*Object, error) {
	if objByVer, ok := r.objByID[id]; ok {
		var objects []*Object
		// Convert objByVer map to a slice
		for _, obj := range objByVer {
			objects = append(objects, obj)
		}

		return objects, nil
	}

	return nil, errObjNotFound
}

// FindResource returns specific resource from registry by object ID, object version and resource ID.
// Uses latest object version, when `objVer` is 0.
// Returns an error, when resource or object not found.
func (r *Registry) FindResource(objID, resID int64, objVer float64) (*Resource, error) {
	_, res, err := r.findResource(objID, resID, objVer)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// findResource returns specific resource and object containing it.
// Finds in registry by object ID, object version and resource ID.
// Uses latest object version, when `objVer` is 0.
// Returns an error, when resource or object not found.
func (r *Registry) findResource(objID, resID int64, objVer float64) (*Object, *Resource, error) {
	obj, err := r.FindObject(objID, objVer)
	if err != nil {
		return nil, nil, err
	}

	// Find a resource by ID
	for _, res := range obj.Resources.Item {
		if res.ID == resID {
			return obj, &res, nil
		}
	}

	return nil, nil, errResNotFound
}

// FindResourcesByID finds resources in registry by ID.
// Returns matching resources from all objects of all versions.
// Returns an error, when resource not found.
func (r *Registry) FindResourcesByID(id int64) ([]Resource, error) {
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

// FindByOIR returns Object and Resource corresponding to OIR.
// Uses latest object version, when `ver` is 0.
func (r *Registry) FindByOIR(oir string, ver float64) (*Object, *Resource, error) {
	objID, _, resID, err := ipso.ParseOIR(oir)
	if err != nil {
		return nil, nil, err
	}

	obj, res, err := r.findResource(objID, resID, ver)
	if err != nil {
		return nil, nil, err
	}

	return obj, res, nil
}

// getObjectsMeta retrieve all objects metadata.
func (r *Registry) getObjectsMeta() ([]ObjectMeta, error) {
	body, err := httpGet("http://www.openmobilealliance.org/api/lwm2m/v1/Object")
	if err != nil {
		return nil, err
	}

	var objectsMeta []ObjectMeta

	if err := json.Unmarshal(body, &objectsMeta); err != nil {
		return nil, err
	}

	return objectsMeta, nil
}

// getObject fetch object details based on metadata.
func (r *Registry) getObject(objectMeta ObjectMeta) (*Object, error) {
	body, err := httpGet(objectMeta.ObjectLink)
	if err != nil {
		return nil, err
	}

	var lwm2m Lwm2M

	if err := xml.Unmarshal(body, &lwm2m); err != nil {
		return nil, err
	}

	return lwm2m.Object, nil
}

// httpGet performs HTTP GET of specified URL.
// Retries several times before returning an error.
func httpGet(url string) ([]byte, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 30 * time.Second            // timeout for all retries
	bo.MaxInterval = time.Second * 5                // max wait time until next retry
	client := http.Client{Timeout: 3 * time.Second} // timeout for single request

	var resp *http.Response

	err := backoff.Retry(func() error {
		var err error

		resp, err = client.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		return nil
	}, bo)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}
