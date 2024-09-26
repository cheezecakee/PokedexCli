package internal

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type LocationAreaResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

// Structs for Detailed Pokémon API
type PokemonType struct {
	Name string `json:"name"`
}

type TypeEntry struct {
	Type PokemonType `json:"type"`
}

type PokemonDetails struct {
	Name           string      `json:"name"`
	ID             int         `json:"id"`
	BaseExperience int         `json:"base_experience"`
	Height         int         `json:"height"`
	Weight         int         `json:"weight"`
	Types          []TypeEntry `json:"types"`
}
