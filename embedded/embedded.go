// Package embedded provides the sanitized IPSO registry compiled into the binary.
//
// Importing the package adds the registry data (about 3 MB) to the executable.
// Use [github.com/aliakseiz/ipso] directly to load a registry from a file, a stream or the OMA API instead.
package embedded

import (
	"bytes"
	_ "embed" // registry.yaml
	"fmt"

	"github.com/aliakseiz/ipso"
)

//go:embed registry.yaml
var registryYAML []byte

// New returns a registry loaded from the embedded sanitized registry file.
// Every call parses the data into a new, independent registry.
func New() (*ipso.Reg, error) {
	reg, err := ipso.New(ipso.Configuration{})
	if err != nil {
		return nil, fmt.Errorf("cannot create registry: %w", err)
	}

	err = reg.Import(bytes.NewReader(registryYAML))
	if err != nil {
		return nil, fmt.Errorf("cannot import embedded registry: %w", err)
	}

	return reg, nil
}
