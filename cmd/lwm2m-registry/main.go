package main

import (
	"log"
	"net/url"
	"time"

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

func main() {
	// Get all objects meta
	objects, err := GetObjectsMeta()
	if err != nil {
		log.Fatal(err)
	}
	_ = objects

	// Get objects details one-by-one
	err = GetObject()
	if err != nil {
		log.Fatal(err)
	}
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

// TODO implement structure with New() which would accept flag `InitializeRegistry`
//  to cache meta and/or details about objects in order to provide search interface
// TODO implement `Refresh` to update the cached registry
// TODO implement `FindObjectByID`, `FindResourceByID`,`FindObjectByName`, `FindObjectByNameLike`,
//  `FindResourceByName`,`FindResourceByNameLike`,`FindObjectByDescriptionLike`, `FindResourceByDescriptionLike`
// TODO implement `Throttle` setting to decrease load on OMAs API
// TODO implement registry export and import to/from file

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
			ID:                 "getObject",
			Method:             "GET",
			PathPattern:        objURL.Path,
			// ProducesMediaTypes: []string{"text/xml"},
			// ConsumesMediaTypes: []string{"text/xml"},
			Schemes:            []string{"https"},
			Params:             objectParams,
			Reader:             &object.GetObjectReader{},
			Context:            objectParams.Context,
			Client:             objectParams.HTTPClient,
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
		time.Sleep(1)
	}

	return nil
}
