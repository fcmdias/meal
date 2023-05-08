package models

type Recipe struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
	Directions  []string `json:"directions"`
}

type Diet struct {
	Name       string
	Value      string
	IsDisabled bool
	IsChecked  bool
	Text       string
}

type DietName int

const (
	Vegan = iota
	Vegetarian
	Omnivore
)
