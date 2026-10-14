package ipso_test

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aliakseiz/ipso"
)

// temperatureYAML is a one-object registry in the format written by Export.
const temperatureYAML = `- Name: Temperature
  Description1: "Temperature sensor.\n"
  ObjectID: 3303
  ObjectURN: urn:oma:lwm2m:ext:3303
  LWM2MVersion: "1.0"
  ObjectVersion: "1.0"
  MultipleInstances: Multiple
  Mandatory: Optional
  Resources:
    Item:
      - id: 5700
        name: Sensor Value
        operations: R
        multipleinstances: Single
        mandatory: Mandatory
        type: Float
        units: Cel
`

func ExampleReg_Import() {
	reg, err := ipso.New(ipso.Configuration{})
	if err != nil {
		fmt.Println(err)

		return
	}

	err = reg.Import(strings.NewReader(temperatureYAML))
	if err != nil {
		fmt.Println(err)

		return
	}

	obj, err := reg.FindObjectByIDAndVer(3303, "1.0")
	fmt.Println(obj.Name, len(obj.Resources.Item), err)
	// Output: Temperature 1 <nil>
}

func ExampleReg_FindResourceByOIR() {
	reg := newExampleReg()

	res, err := reg.FindResourceByOIR("/3303/0/5700")
	fmt.Println(res.Name, res.Type, res.Units, err)

	_, err = reg.FindResourceByOIR("3303/5700")
	fmt.Println(errors.Is(err, ipso.ErrInvalidOIR))
	// Output:
	// Sensor Value Float Cel <nil>
	// true
}

func ExampleReg_FindObjectsByName() {
	reg := newExampleReg()

	objects, err := reg.FindObjectsByName("temperature") // letter case is ignored
	fmt.Println(len(objects), objects[0].ObjectID, err)

	_, err = reg.FindObjectsByName("Humidity")
	fmt.Println(errors.Is(err, ipso.ErrObjectNotFound))
	// Output:
	// 1 3303 <nil>
	// true
}

func ExampleReg_FindResourcesByName() {
	reg := newExampleReg()

	resources, err := reg.FindResourcesByName("Sensor Value")
	fmt.Println(len(resources), resources[0].ID, err)
	// Output: 1 5700 <nil>
}

func ExampleReg_Sanitize() {
	reg := newExampleReg()
	fmt.Printf("%q\n", reg.GetObjects()[0].Description1)

	reg.Sanitize(ipso.DefaultSanitizer())
	fmt.Printf("%q\n", reg.GetObjects()[0].Description1)
	// Output:
	// "Temperature sensor.\n"
	// "Temperature sensor."
}

func ExampleReg_Compare() {
	older := newExampleReg()

	newer, err := ipso.New(ipso.Configuration{})
	if err != nil {
		fmt.Println(err)

		return
	}

	err = newer.Import(strings.NewReader(temperatureYAML + strings.ReplaceAll(temperatureYAML, "3303", "3304")))
	if err != nil {
		fmt.Println(err)

		return
	}

	for _, diff := range older.Compare(newer) {
		fmt.Println(diff.Difference, diff.ObjectComp.ObjectID)
	}
	// Output: New object added 3304
}

func ExampleReg_Export() {
	reg := newExampleReg()

	err := reg.Export(os.Stdout)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// - Name: Temperature
	//   Description1: Temperature sensor.
	//   Description2: ""
	//   ObjectID: 3303
	//   ObjectURN: urn:oma:lwm2m:ext:3303
	//   LWM2MVersion: "1.0"
	//   ObjectVersion: "1.0"
	//   MultipleInstances: Multiple
	//   Mandatory: Optional
	//   Resources:
	//     Item:
	//         - id: 5700
	//           name: Sensor Value
	//           operations: R
	//           multipleinstances: Single
	//           mandatory: Mandatory
	//           type: Float
	//           rangeenumeration: ""
	//           units: Cel
	//           description: ""
}

// newExampleReg returns a registry holding the object from temperatureYAML.
func newExampleReg() *ipso.Reg {
	reg, err := ipso.New(ipso.Configuration{})
	if err != nil {
		panic(err)
	}

	err = reg.Import(strings.NewReader(temperatureYAML))
	if err != nil {
		panic(err)
	}

	return reg
}
