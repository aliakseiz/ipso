package registry

import (
	"errors"
	"strconv"
	"strings"
)

// URN model according to OMA-TS-LightweightM2M_Core-V1_1_1-20190617-A, chapter 7.2.2 "Object version format".
type URN struct {
	Prefix  string  // urn
	Vendor  string  // oma
	Spec    string  // lwm2m
	Type    string  // oma/ext/x
	ObjID   int64   // object ID
	Version float64 // object version, optional
}

var (
	errInvalidURN    = errors.New("invalid URN")
	errInvalidObjID  = errors.New("invalid object ID")
	errInvalidObjVer = errors.New("invalid object version")
)

// parseURN parses OMA URN to a structure.
// Valid URN example: urn:oma:lwm2m:oma:0:1.1
func parseURN(s string) (*URN, error) {
	var (
		ver float64
		err error
	)

	el := strings.Split(s, ":")

	switch len(el) {
	case 6:
		// URN of six elements contains optional object version
		ver, err = strconv.ParseFloat(el[5], 64)
		if err != nil {
			return nil, errInvalidObjVer
		}

		fallthrough
	case 5:
		id, err := strconv.ParseInt(el[4], 10, 64)
		if err != nil {
			return nil, errInvalidObjID
		}

		return &URN{
			Prefix:  el[0], // urn
			Vendor:  el[1], // oma
			Spec:    el[2], // lwm2m
			Type:    el[3], // oma
			ObjID:   id,    // 0
			Version: ver,   // 1.1
		}, nil
	default:
		return nil, errInvalidURN
	}
}
