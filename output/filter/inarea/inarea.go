package inarea

import (
	"errors"

	"chaos.expert/FreifunkBremen/yanic/output/filter"
	"chaos.expert/FreifunkBremen/yanic/runtime"
)

type area struct {
	latitudeMin  float64
	latitudeMax  float64
	longitudeMin float64
	longitudeMax float64
}

func init() {
	filter.Register("in_area", build)
}

func build(config interface{}) (filter.Filter, error) {
	values, ok := config.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid configuration, map expected")
	}

	a := area{}
	if v, ok := values["latitude_min"]; ok {
		a.latitudeMin = v.(float64)
	}
	if v, ok := values["latitude_max"]; ok {
		a.latitudeMax = v.(float64)
	}
	if v, ok := values["longitude_min"]; ok {
		a.longitudeMin = v.(float64)
	}
	if v, ok := values["longitude_max"]; ok {
		a.longitudeMax = v.(float64)
	}

	if a.latitudeMin >= a.latitudeMax {
		return nil, errors.New("invalid latitude: max is bigger then min")
	}
	if a.longitudeMin >= a.longitudeMax {
		return nil, errors.New("invalid longitude: max is bigger then min")
	}

	// TODO bessere Fehlerbehandlung!

	return &a, nil
}

func (a *area) Apply(node *runtime.Node) *runtime.Node {
	if nodeinfo := node.Nodeinfo; nodeinfo != nil {
		location := nodeinfo.Location
		if location == nil {
			return node
		}
		if location.Latitude >= a.latitudeMin && location.Latitude <= a.latitudeMax && location.Longitude >= a.longitudeMin && location.Longitude <= a.longitudeMax {
			return node
		}
		return nil
	}
	return node
}
