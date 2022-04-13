package main

import (
	"log"

	"github.com/aliakseiz/ipso"
)

func main() {
	// Initialize registry from OMA API using default configuration
	regAPI, err := ipso.New(ipso.DefaultConfiguration())
	if err != nil {
		panic(err)
	}
	// Store imported registry in the file
	if err = regAPI.Export("registry.yaml"); err != nil {
		panic(err)
	}

	// Initialize another registry from file
	cfg := ipso.Configuration{
		InitOnNew:      false,
		SkipInitErrors: false,
		Sanitize:       false,
		Sanitizer:      nil,
	}

	regFile, err := ipso.New(cfg)
	if err != nil {
		panic(err)
	}

	if err = regFile.Import("registry.yaml"); err != nil {
		panic(err)
	}

	objComp := regFile.Compare(regAPI.GetRegistry())

	log.Printf("objComp length: %d", len(objComp))

	// Sanitize registry imported from API
	regFile.Sanitize(ipso.DefaultSanitizer())
	// Store sanitized registry in file
	if err = regFile.Export("registry_sanitized.yaml"); err != nil {
		panic(err)
	}

	log.Printf("finished successfully")
}
