package main

import (
	"log"

	"github.com/aliakseiz/ipso-registry/registry"
)

func main() {
	// Initialize registry from OMA API using default configuration
	regAPI, err := registry.New(registry.DefaultConfiguration())
	if err != nil {
		panic(err)
	}
	// Store imported registry in file
	if err = regAPI.Export("registry.yaml"); err != nil {
		panic(err)
	}

	// Initialize another registry from file
	cfg := &registry.Configuration{
		InitOnNew:      false,
		SkipInitErrors: false,
		Sanitize:       false,
		Sanitizer:      nil,
	}

	regFile, err := registry.New(cfg)
	if err != nil {
		panic(err)
	}

	if err = regFile.Import("registry.yaml"); err != nil {
		panic(err)
	}

	objComp := regFile.Compare(regAPI)

	log.Printf("objComp length: %d", len(objComp))

	// Sanitize registry imported from API
	regFile.Config.Sanitizer = registry.DefaultSanitizer()
	regFile.Sanitize()
	// Store sanitized registry in file
	if err = regFile.Export("registry_sanitized.yaml"); err != nil {
		panic(err)
	}
}
