// Package registry provides methods to create and control OMA IPSO registry.
package registry

import (
	"context"
	"errors"
	"io/ioutil"
	"net/url"
	"strings"

	openapi "github.com/aliakseiz/ipso-registry/api/client"
	"gopkg.in/yaml.v2"
)

const (
	objectsMetaBasePath = "/api/lwm2m/v1/"
)

var (
	errEmptyObjectURL = errors.New("empty object description URL")
	errEmptyFilename  = errors.New("filename is empty")
	errObjNotFound    = errors.New("object not found")
	errResNotFound    = errors.New("resource not found")
)

// Registry holds objects and settings.
type Registry struct {
	Config *Configuration

	Objects []Object
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
			object := &r.Objects[oIndex] // Modify the object in registry instead of object's copy
			object.Description1 = strings.ReplaceAll(object.Description1, s, "")
			object.Description1 = strings.TrimSpace(object.Description1)
			object.Description2 = strings.ReplaceAll(object.Description2, s, "")
			object.Description2 = strings.TrimSpace(object.Description2)

			for rIndex := 0; rIndex < len(r.Objects[oIndex].Resources); rIndex++ {
				resource := &r.Objects[oIndex].Resources[rIndex]
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
	if filename == "" {
		return errEmptyFilename
	}

	data, err := yaml.Marshal(&r.Objects)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

// Import loads objects and resources from file.
// Overwrites current registry objects and rObesources.
func (r *Registry) Import(filename string) error {
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
func (r *Registry) ImportFromAPI() ([]Object, error) {
	objectsMeta, err := r.getObjectsMeta()
	if err != nil {
		return nil, err
	}

	var objects []Object

	for _, objectMeta := range objectsMeta {
		object, err := r.getObject(objectMeta)
		if err != nil {
			if r.Config.SkipInitErrors {
				continue
			}

			return nil, err
		}

		objects = append(objects, mapObject(object))
	}

	return objects, nil
}

// Compare makes comparison of r and reg registries.
// Returns a list of non-equal objects with difference explanation.
func (r *Registry) Compare(reg *Registry) []ObjectComparison {
	// TODO store objects in registry in maps to improve lookup performance
	var objComp []ObjectComparison

	// Compare r with reg
	for _, regObj := range reg.Objects {
		if rObj, err := r.Find(&regObj); err != nil {
			regObjCopy := regObj
			// rObjCopy := rObj

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
		if regObj, err := reg.Find(&rObj); err != nil {
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

// TODO implement `FindObjectByName`, `FindResourceByName`, `FindObjectByDescription`, `FindResourceByDescription`

// Find looks for an object in current registry.
// Returns an empty object and error, when object not found.
func (r *Registry) Find(o *Object) (Object, error) {
	for _, rObj := range r.Objects {
		if rObj.ObjectID == o.ObjectID && rObj.ObjectVersion == o.ObjectVersion {
			return rObj, nil
		}
	}

	return Object{}, errObjNotFound
}

// FindObjectsByID finds objects in registry by ID.
// Multiple objects with same ID and different versions could be returned.
// Returns an error, when object not found.
func (r *Registry) FindObjectsByID(id int32) ([]Object, error) {
	var objects []Object

	for _, rObj := range r.Objects {
		if rObj.ObjectID == id {
			objects = append(objects, rObj)
		}
	}

	if len(objects) == 0 {
		return nil, errObjNotFound
	}

	return objects, nil
}

// FindResourcesByID finds resources in registry by ID.
// Returns matching resources from all objects of all versions.
// Returns an error, when resource not found.
func (r *Registry) FindResourcesByID(id int32) ([]Resource, error) {
	var resources []Resource

	for _, rObj := range r.Objects {
		for _, rRes := range rObj.Resources {
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

// FindResourcesByObjResIDs finds specific resource in registry by object ID and resource ID.
// Returns matching resources from all versions of specific object.
// Returns an error, when resource or object not found.
func (r *Registry) FindResourcesByObjResIDs(objID, resID int32) ([]Resource, error) {
	var resources []Resource

	for _, rObj := range r.Objects {
		if rObj.ObjectID == objID {
			for _, rRes := range rObj.Resources {
				if rRes.ID == resID {
					resources = append(resources, rRes)
				}
			}
		}
	}

	if len(resources) == 0 {
		return nil, errResNotFound
	}

	return resources, nil
}

// getObjectsMeta retrieve all objects metadata.
func (r *Registry) getObjectsMeta() ([]openapi.ObjectMeta, error) {
	cfg := openapi.NewConfiguration()
	cfg.BasePath += objectsMetaBasePath
	client := openapi.NewAPIClient(cfg)

	// TODO allow to choose object version
	// ObjectVersion: optional.NewString(objectVersionLatest)

	objects, _, err := client.ObjectsApi.FindObjects(context.Background(), &openapi.FindObjectsOpts{})
	if err != nil {
		return nil, err
	}

	return objects, nil
}

// getObject fetch object details based on metadata.
func (r *Registry) getObject(objectMeta openapi.ObjectMeta) (*openapi.Object, error) {
	cfg := openapi.NewConfiguration()
	client := openapi.NewAPIClient(cfg)

	objURL, err := url.Parse(objectMeta.ObjectLink)
	if err != nil {
		return nil, err
	}

	if objURL.Path == "" {
		return nil, errEmptyObjectURL
	}

	object, _, err := client.ObjectApi.FindObject(context.Background(), objURL.Path[1:]) // 1: - skip leading /
	if err != nil {
		return nil, err
	}

	return &object.Object, nil
}
