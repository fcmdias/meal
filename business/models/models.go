package models

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID          uuid.UUID
	Name        string
	Description string
	Ingredients []string
	Directions  []string
	DateCreated time.Time
	DateUpdated *time.Time `json:",omitempty"`
}

type NewRecipe struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Ingredients []string `json:"ingredients" validate:"required"`
	Directions  []string `json:"directions" validate:"required"`
}

func NewRecipeToRecipe(n NewRecipe) Recipe {
	return Recipe{
		ID:          uuid.New(),
		Name:        n.Name,
		Description: n.Description,
		Ingredients: n.Ingredients,
		Directions:  n.Directions,
		DateCreated: time.Now(),
	}
}

type Diet struct {
	Name       string
	Value      string
	IsDisabled bool
	IsChecked  bool
	Text       string
}

type DietType int

const (
	Vegan DietType = iota
	Vegetarian
	Omnivore
)
