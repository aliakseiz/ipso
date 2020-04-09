package main

import (
	"log"

	"github.com/aliakseiz/lwm2m-registry/registry"
)

func main() {
	regAPI, err := registry.New(registry.DefaultConfiguration())
	if err != nil {
		panic(err)
	}

	if err = regAPI.Export("registry.yaml"); err != nil {
		panic(err)
	}

	cfg := &registry.Configuration{
		InitOnNew:      false,
		SkipInitErrors: false,
	}

	reg1, err := registry.New(cfg)
	if err != nil {
		panic(err)
	}

	if err = reg1.Import("registry.yaml"); err != nil {
		panic(err)
	}

	reg2, err := registry.New(cfg)
	if err != nil {
		panic(err)
	}

	if err = reg2.Import("registry2.yaml"); err != nil {
		panic(err)
	}

	objComp := reg1.Compare(reg2)

	log.Printf("objComp length: %d", len(objComp))
}
