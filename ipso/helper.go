package ipso

import (
	"strconv"
	"strings"
)

// ParseOIRS parse OIR (object/instance/resource) string.
func ParseOIR(oirStr string) (objectID, instanceNumber, resourceID int64, err error) {
	// Parse object/instance/resource string
	elementIds := strings.Split(oirStr, "/")

	objectID, err = strconv.ParseInt(elementIds[0], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	instanceNumber, err = strconv.ParseInt(elementIds[1], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	resourceID, err = strconv.ParseInt(elementIds[2], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	return objectID, instanceNumber, resourceID, nil
}
