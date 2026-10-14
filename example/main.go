// Command example shows how to load, search, compare and export the IPSO registry.
//
// By default it uses the registry embedded in the module and works without the network and file system.
// With -fetch it downloads the registry from the OMA API, stores it as YAML in the -out directory
// and compares it with the embedded one.
package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/aliakseiz/ipso"
	"github.com/aliakseiz/ipso/embedded"
)

const outDirMode = 0o750

func main() {
	fetch := flag.Bool("fetch", false, "download the registry from the OMA API and write YAML files")
	out := flag.String("out", ".", "directory for the files written with -fetch")

	flag.Parse()

	reg, err := embedded.New()
	if err != nil {
		log.Fatal(err)
	}

	res, err := reg.FindResourceByOIR("3303/0/5700")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("3303/0/5700: %s (%s)", res.Name, res.Type)

	if !*fetch {
		return
	}

	err = syncFromAPI(reg, *out)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("finished successfully")
}

// syncFromAPI downloads the registry from the OMA API, stores it in dir, and logs how it differs from base.
//
//nolint:wrapcheck // main logs the error as is
func syncFromAPI(base *ipso.Reg, dir string) error {
	err := os.MkdirAll(dir, outDirMode)
	if err != nil {
		return err
	}

	regAPI, err := ipso.New(ipso.DefaultConfiguration())
	if err != nil {
		return err
	}

	file := filepath.Join(dir, "registry.yaml")

	err = regAPI.ExportFile(file)
	if err != nil {
		return err
	}

	log.Printf("differences against embedded registry: %d", len(base.Compare(regAPI)))

	regAPI.Sanitize(ipso.DefaultSanitizer())

	return regAPI.ExportFile(filepath.Join(dir, "registry_sanitized.yaml"))
}
