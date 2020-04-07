package main

import (
	"log"

	"github.com/aliakseiz/lwm2m-registry/registry"
)

// TODO implement unit tests
func main() {
	r := registry.New(registry.DefaultConfiguration())
	err := r.Save("registry.yaml")
	if err != nil {
		panic(err)
	}
	log.Printf("loaded")
}
