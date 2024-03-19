package report

import (
	"fmt"
	"worker-report-matrix/internal/geolocation"

	"github.com/bradfitz/latlong"
)

type GeoLoc struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
	TimeZone  string  `json:"timezone"`
}

func (g *GeoLoc) Initialize(lat float64, long float64) *GeoLoc {
	timeZone := latlong.LookupZoneName(lat, long)

	g.Latitude = lat
	g.Longitude = long
	g.TimeZone = timeZone

	return g
}

func (g *GeoLoc) GetAddressFromMapBox() *GeoLoc {
	address, err := geolocation.GetAddressFromMapBox(g.Latitude, g.Longitude)
	if err != nil {
		fmt.Println("[ERROR] geolocation.GetAddressFromMapBox:", err)
		return g
	}
	g.Address = address
	return g
}
