package registry

import (
	"context"
	"errors"
	"net/url"

	openapi "github.com/aliakseiz/lwm2m-registry/api/client"

	"github.com/antihax/optional"
)

const (
	objectsMetaBasePath = "/api/lwm2m/v1/"
	objectVersionLatest = "latest"
)

var (
	errNotSupportedSourceType = errors.New("not supported source type")
	errEmptyObjectURL         = errors.New("empty object description URL")
)

// Registry holds objects and settings
type Registry struct {
	Config *Configuration

	Objects []Object

	state            State
	stateDescription string
}

// TODO implement `FindObjectByID`, `FindResourceByID`,`FindObjectByName`, `FindResourceByName`,
//  `FindObjectByDescription`, `FindResourceByDescription`

// TODO implement registry export/import to/from file

// New creates a new registry, using provided or default configuration
func New(cfg *Configuration) *Registry {
	var err error

	reg := &Registry{
		Config:           cfg,
		Objects:          nil,
		state:            StateNotInitialized,
		stateDescription: "",
	}

	if reg.Config == nil {
		reg.Config = DefaultConfiguration()
	}

	if cfg.InitOnNew {
		// TODO initialize registry
		err = reg.Load()
		if err != nil {
			reg.state = StateNotInitialized
			reg.stateDescription = err.Error()
		} else {
			reg.state = StateInitialized
		}
	}

	return reg
}

// Load loads objects and resources from configured source
func (r *Registry) Load() (err error) {
	switch r.Config.Source {
	case SourceTypeAPI:
		r.Objects, err = r.loadFromAPI()
	case SourceTypeFile:
		// TODO
	default:
		return errNotSupportedSourceType
	}

	return
}

// State returns registry State
// Mandatory wrapper to protect State field from external changes
func (r *Registry) State() State {
	return r.state
}

// StateDescription returns registry StateDescription
func (r *Registry) StateDescription() string {
	return r.stateDescription
}

// LoadFromURL initialize the registry from official OMA API
func (r *Registry) loadFromAPI() ([]Object, error) {
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

// getObjectsMeta retrieve all objects metadata
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

// getObject fetch object details based on metadata
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
