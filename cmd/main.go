package main

import (
	"log"

	"github.com/aliakseiz/lwm2m-registry/registry"
)

// TODO implement unit tests
func main() {
	// reg, err := registry.New(registry.DefaultConfiguration())
	// if err != nil {
	// 	panic(err)
	// }

	// if err = reg.Export("registry.yaml"); err != nil {
	// 	panic(err)
	// }

	cfg := &registry.Configuration{
		InitOnNew:      false,
		SkipInitErrors: false,
	}

	reg, err := registry.New(cfg)
	if err != nil {
		panic(err)
	}

	if err = reg.Import("registry.yaml"); err != nil {
		panic(err)
	}

	reg2, err := registry.New(cfg)
	if err != nil {
		panic(err)
	}

	if err = reg2.Import("registry2.yaml"); err != nil {
		panic(err)
	}

	objComp := reg.Compare(reg2)

	log.Printf("objComp length: %d", len(objComp))
}
