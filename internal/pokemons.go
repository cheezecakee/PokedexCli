package internal

type Pokemons struct {
	ID                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	GameIndex            int                    `json:"game_index"`
	EncounterMethodRates []EncounterMethodRates `json:"encounter_method_rates"`
	Location             Location               `json:"location"`
	Names                []Names                `json:"names"`
	PokemonEncounters    []PokemonEncounters    `json:"pokemon_encounters"`
}

type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VersionDetail struct {
	Rate    int     `json:"rate"`
	Version Version `json:"version"`
}

type EncounterMethodRates struct {
	EncounterMethod EncounterMethod `json:"encounter_method"`
	VersionDetails  []VersionDetail `json:"version_details"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Names struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Method struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EncounterDetails struct {
	MinLevel        int    `json:"min_level"`
	MaxLevel        int    `json:"max_level"`
	ConditionValues []any  `json:"condition_values"`
	Chance          int    `json:"chance"`
	Method          Method `json:"method"`
}

type VersionDetails struct {
	Version          Version            `json:"version"`
	MaxChance        int                `json:"max_chance"`
	EncounterDetails []EncounterDetails `json:"encounter_details"`
}

type PokemonEncounters struct {
	Pokemon        Pokemon          `json:"pokemon"`
	VersionDetails []VersionDetails `json:"version_details"`
}
