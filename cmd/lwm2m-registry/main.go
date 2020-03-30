package main

import (
	"context"
	"errors"
	"log"
	"net/url"

	openapi "github.com/aliakseiz/lwm2m-registry-importer/api/client"
	"github.com/antihax/optional"
)

const (
	objectsMetaBasePath = "/api/lwm2m/v1/"

	objectVersion = "latest"
)

// TODO implement unit tests
func main() {
	objects, err := GetObjectsMeta()
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range objects {
		if err := GetObject(object); err != nil {
			log.Print(err)
		}
	}
}

// GetObjectsMeta retrieve all objects metadata
func GetObjectsMeta() ([]openapi.ObjectMeta, error) {
	cfg := openapi.NewConfiguration()
	cfg.BasePath += objectsMetaBasePath
	client := openapi.NewAPIClient(cfg)

	objects, _, err := client.ObjectsApi.FindObjects(context.Background(), &openapi.FindObjectsOpts{
		ObjectVersion: optional.NewString(objectVersion),
	})
	if err != nil {
		return nil, err
	}

	return objects, nil
}

// GetObject fetch object details based on metadata
func GetObject(objectMeta openapi.ObjectMeta) error {
	cfg := openapi.NewConfiguration()
	client := openapi.NewAPIClient(cfg)

	objURL, err := url.Parse(objectMeta.ObjectLink)
	if err != nil {
		log.Printf("failed to parse URL: %s", objectMeta.ObjectLink)
		return err
	}

	if objURL.Path == "" {
		log.Printf("empty URL: %s", objectMeta.ObjectLink)
		return errors.New("empty object description URL")
	}

	log.Printf("%d\t%s\n", objectMeta.ObjectID, objectMeta.Name)

	object, _, err := client.ObjectApi.FindObject(context.Background(), objURL.Path[1:]) // 1: - skip leading /
	if err != nil {
		return err
	}

	for _, res := range object.Object.Resources.Item {
		log.Printf("\t%d\t%s", res.ID, res.Name)
	}

	return nil
}
