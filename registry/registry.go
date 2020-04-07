package registry

import (
	"context"
	"errors"
	"io/ioutil"
	"net/url"

	openapi "github.com/aliakseiz/lwm2m-registry/api/client"

	"github.com/antihax/optional"

	"gopkg.in/yaml.v2"
)

const (
	objectsMetaBasePath = "/api/lwm2m/v1/"
	objectVersionLatest = "latest"
)

var (
	errEmptyObjectURL = errors.New("empty object description URL")
	errEmptyFilename  = errors.New("filename is empty")
)

// Registry holds objects and settings.
type Registry struct {
	Config *Configuration

	Objects []Object
}

// TODO implement `FindObjectByID`, `FindResourceByID`,`FindObjectByName`, `FindResourceByName`,
//  `FindObjectByDescription`, `FindResourceByDescription`

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

	return reg, nil
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
// Overwrites current registry Objects and Resources.
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
// Overwrites current registry Objects and Resources.
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

// getObjectsMeta retrieve all objects metadata.
func (r *Registry) getObjectsMeta() ([]openapi.ObjectMeta, error) {
	cfg := openapi.NewConfiguration()
	cfg.BasePath += objectsMetaBasePath
	client := openapi.NewAPIClient(cfg)

	objects, _, err := client.ObjectsApi.FindObjects(context.Background(), &openapi.FindObjectsOpts{
		ObjectVersion: optional.NewString(objectVersionLatest), // TODO allow to choose object version
	})
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
