package geolocation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	MAPBOX_ACCESS_TOKEN = os.Getenv("MAPBOX_ACCESS_TOKEN")
)

type GeoCodingMapBox struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	PlaceName string    `json:"place_name"`
	Context   []Context `json:"context"`
}

type Context struct {
	Id        string `json:"id"`
	ShortCode string `json:"short_code"`
}

func GetAddressFromMapBox(latitude float64, longitude float64) (string, error) {
	req := fmt.Sprintf("https://api.mapbox.com/geocoding/v5/mapbox.places/%f,%f.json?limit=1&types=address&language=fr&access_token=%s", 
		longitude,
		latitude,
		MAPBOX_ACCESS_TOKEN,
	)

	resp, err := http.Get(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var geoCode GeoCodingMapBox
	err = json.Unmarshal(body, &geoCode)
	if err != nil {
		return "", err
	}

	if geoCode.Features == nil || len(geoCode.Features) == 0 {
		return "", fmt.Errorf("[ERROR] No results returned from MapBox API: %s", req)
	}

	return geoCode.Features[0].PlaceName, nil
}


