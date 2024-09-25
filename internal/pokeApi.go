package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetLocations(url string) (Locations, error) {
	res, err := http.Get(url)
	if err != nil {
		return Locations{}, err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return Locations{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	if err != nil {
		return Locations{}, err
	}

	locations, err := unmarshal(body)
	if err != nil {
		return Locations{}, fmt.Errorf("Failed to unmarshal: %v\n", err)
	}

	return locations, nil
}

func unmarshal(data []byte) (Locations, error) {
	locations := Locations{}
	err := json.Unmarshal(data, &locations)
	if err != nil {
		return Locations{}, fmt.Errorf("Failed to unmarshal: %v\n", err)
	}

	return locations, nil
}
