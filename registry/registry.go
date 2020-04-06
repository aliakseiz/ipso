package registry

import (
	"context"
	"errors"
	"log"
	"net/url"

	openapi "github.com/aliakseiz/lwm2m-registry/api/client"

	"github.com/antihax/optional"
)

const (
	objectsMetaBasePath = "/api/lwm2m/v1/"
	objectVersionLatest = "latest"
)

// Registry holds objects and settings
type Registry struct {
	// TODO add field to represent current state, if initializing from the web
	Objects []Object
}

// TODO implement `FindObjectByID`, `FindResourceByID`,`FindObjectByName`, `FindResourceByName`,
//  `FindObjectByDescription`, `FindResourceByDescription`

// TODO implement `Throttle` setting to decrease load on OMAs API

// TODO implement registry export/import to/from file

// New creates a new registry, using custom configuration
func New(cfg *Configuration) *Registry {
	if cfg == nil {
		cfg = DefaultConfiguration()
	}

	return nil // TODO initialize registry
}

// Init initialize the registry
func (r Registry) Init() error {
	objectsMeta, err := r.GetObjectsMeta()
	if err != nil {
		return err
	}

	var objects []*Object

	for _, objectMeta := range objectsMeta {
		object, err := r.GetObject(objectMeta)
		if err != nil { // TODO add flag to skip init error
			return err
		}

		objects = append(objects, mapObject(object))
	}
	return nil
}

// GetObjectsMeta retrieve all objects metadata
func (r Registry) GetObjectsMeta() ([]openapi.ObjectMeta, error) {
	cfg := openapi.NewConfiguration()
	cfg.BasePath += objectsMetaBasePath
	client := openapi.NewAPIClient(cfg)

	objects, _, err := client.ObjectsApi.FindObjects(context.Background(), &openapi.FindObjectsOpts{
		ObjectVersion: optional.NewString(objectVersionLatest),
	})
	if err != nil {
		return nil, err
	}

	return objects, nil
}

// GetObject fetch object details based on metadata
func (r Registry) GetObject(objectMeta openapi.ObjectMeta) (*openapi.Object, error) {
	cfg := openapi.NewConfiguration()
	client := openapi.NewAPIClient(cfg)

	objURL, err := url.Parse(objectMeta.ObjectLink)
	if err != nil {
		log.Printf("failed to parse URL: %s", objectMeta.ObjectLink)
		return nil, err
	}

	if objURL.Path == "" {
		log.Printf("empty URL: %s", objectMeta.ObjectLink)
		return nil, errors.New("empty object description URL")
	}

	log.Printf("%d\t%s\n", objectMeta.ObjectID, objectMeta.Name)

	object, _, err := client.ObjectApi.FindObject(context.Background(), objURL.Path[1:]) // 1: - skip leading /
	if err != nil {
		return nil, err
	}

	for _, res := range object.Object.Resources.Item {
		log.Printf("\t%d\t%s", res.ID, res.Name)
	}

	return &object.Object, nil
}
