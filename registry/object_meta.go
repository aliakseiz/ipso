package registry

type LwM2M struct {
	Object Object
}

type ObjectMeta struct {
	ObjectID          int64  `json:"ObjectID"`
	Ver               string `json:"Ver"`
	URN               string `json:"URN"`
	Name              string `json:"Name"`
	Description       string `json:"Description"`
	Owner             string `json:"Owner"`
	Label             string `json:"Label"`
	ObjectLink        string `json:"ObjectLink"`
	ObjectLinkVisible string `json:"ObjectLinkVisible"`
	SpecLink          string `json:"SpecLink"`
	SpecLinkVisible   string `json:"SpecLinkVisible"`
	VortoLink         string `json:"VortoLink"`
}
