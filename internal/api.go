package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func fetchApi(url string, target interface{}, cache *Cache) error {
	if cachedData, ok := cache.Get(url); ok {
		return json.Unmarshal(cachedData, target)
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	if err != nil {
		return err
	}

	cache.Add(url, body)

	err = json.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal: %v\n", err)
	}

	return nil
}

func GetPokemonsInArea(url, name string, cache *Cache) (LocationAreaResponse, error) {
	var locationArea LocationAreaResponse
	url = url + name
	err := fetchApi(url, &locationArea, cache)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	return locationArea, nil
}

func GetPokemon(url, name string, cache *Cache) (PokemonDetails, error) {
	var pokemonDetails PokemonDetails
	url = url + name
	err := fetchApi(url, &pokemonDetails, cache)
	if err != nil {
		return PokemonDetails{}, err
	}

	return pokemonDetails, nil
}

func GetLocations(url string, cache *Cache) (Locations, error) {
	var locations Locations
	err := fetchApi(url, &locations, cache)
	if err != nil {
		return Locations{}, err
	}

	return locations, nil
}
