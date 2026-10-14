package embedded_test

import (
	"fmt"

	"github.com/aliakseiz/ipso/embedded"
)

func ExampleNew() {
	reg, err := embedded.New() // parses the data on every call, keep and reuse the registry
	if err != nil {
		fmt.Println(err)

		return
	}

	res, err := reg.FindResourceByOIR("3303/0/5700")
	fmt.Println(res.Name, err)
	// Output: Sensor Value <nil>
}
