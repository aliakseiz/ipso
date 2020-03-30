package main

import (
	"context"
	"errors"
	"log"
	"net/url"
	"time"

	openapi "github.com/aliakseiz/lwm2m-registry-importer/api/client"
	"github.com/aliakseiz/lwm2m-registry-importer/rest/client/object"
	"github.com/aliakseiz/lwm2m-registry-importer/rest/models"
	"github.com/go-openapi/runtime"

	"github.com/aliakseiz/lwm2m-registry-importer/rest/client"
	httpTransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

var (
	objectVersion = "latest"
)

// TODO implement unit tests
func main() {
	objects, err := GetObjectsMetaOpenAPI()
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range objects {
		if err := GetObjectOpenAPI(object); err != nil {
			log.Print(err)
		}
	}

	// Get all objects meta
	/*objects, err := GetObjectsMeta()
	if err != nil {
		log.Fatal(err)
	}
	_ = objects

	// Get objects details one-by-one
	err = GetObject()
	if err != nil {
		log.Fatal(err)
	}*/
}

// GetObjectsMeta retrieve all objects meta
func GetObjectsMeta() (models.Objects, error) {
	httpClient := client.NewHTTPClient(strfmt.Default)

	objectsParams := object.NewGetObjectsParams()
	objectsParams.ObjectVersion = &objectVersion

	objects, err := httpClient.Object.GetObjects(objectsParams)
	if err != nil {
		return nil, err
	}
	return objects.Payload, nil
}

func GetObjectsMetaOpenAPI() ([]openapi.ObjectMeta, error) {
	cfg := openapi.NewConfiguration()
	cfg.BasePath += "/api/lwm2m/v1/"
	client := openapi.NewAPIClient(cfg)

	objects, _, err := client.ObjectsApi.FindObjects(context.Background(), &openapi.FindObjectsOpts{})
	if err != nil {
		return nil, err
	}

	return objects, nil
}

// GetObject retrieve single object
func GetObject() error {
	objects, err := GetObjectsMeta()
	if err != nil {
		return err
	}

	transport := httpTransport.New(client.DefaultHost, "", []string{"https"})
	transport.Producers["text/xml"] = runtime.XMLProducer()
	transport.Consumers["text/xml"] = runtime.XMLConsumer()
	httpClient := client.New(transport, strfmt.Default)
	objectParams := object.NewGetObjectParams()

	for _, obj := range objects {
		log.Printf("%d\t%s\n", obj.ObjectID, obj.Name)

		objURL, err := url.Parse(obj.ObjectLink)
		if err != nil {
			log.Printf("failed to parse URL: %s", obj.ObjectLink)
			continue
		}

		if len(objURL.Path) == 0 {
			log.Printf("no URL for object: %d", obj.ObjectID)
			continue
		}

		result, err := httpClient.Transport.Submit(&runtime.ClientOperation{
			ID:          "getObject",
			Method:      "GET",
			PathPattern: objURL.Path,
			Schemes:     []string{"https"},
			Params:      objectParams,
			Reader:      &object.GetObjectReader{},
			Context:     objectParams.Context,
			Client:      objectParams.HTTPClient,
		})
		if err != nil {
			log.Fatal(err)
		}

		objResp, ok := result.(*object.GetObjectOK)
		if !ok {
			log.Printf("failed to unmarshal response for object: %d", obj.ObjectID)
			continue
		}

		for _, res := range objResp.Payload.Object.Resources.Item {
			log.Printf("\t%d\t%s", res.ID, res.Name)
		}

		time.Sleep(1) // TODO take Throttle parameter from registry
	}

	return nil
}

func GetObjectOpenAPI(objectMeta openapi.ObjectMeta) error {
	cfg := openapi.NewConfiguration()
	// cfg.BasePath += "/tech/profiles/"
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

	path := objURL.Path[1:]
	object, resp, err := client.ObjectApi.FindObject(context.Background(), path)
	if err != nil {
		return err
	}
	_ = resp
	for _, res := range object.Object.Resources.Item {
		log.Printf("\t%d\t%s", res.ID, res.Name)
	}

	return nil
}
