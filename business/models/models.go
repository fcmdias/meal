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
	// The estimated time required to cook the recipe.
	CookingTime time.Duration
	// The estimated time required for the recipe's preparation.
	PreparationTime time.Duration
	// An indication of the recipe's complexity or difficulty.
	Difficulty int
	// The culinary tradition or regional origin of the recipe (e.g., Italian, Mexican, Indian).
	CuisineType []string
	// The category or type of meal the recipe belongs to, (e.g., breakfast, lunch, dinner, starter, dessert or snack).
	MealType string
	// Tags or flags indicating dietary considerations (e.g., vegan, gluten-free, vegetarian).
	DietaryRestrictions []string
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

type RecipeEdit struct {
	Name        *string   `json:"name" validate:"required"`
	Description *string   `json:"description" validate:"required"`
	Ingredients *[]string `json:"ingredients" validate:"required"`
	Directions  *[]string `json:"directions" validate:"required"`
}

func RecipeEditToRecipe(r *Recipe, e RecipeEdit) {

	if e.Name != nil {
		r.Name = *e.Name
	}
	if e.Description != nil {
		r.Description = *e.Description
	}
	if e.Ingredients != nil {
		r.Ingredients = *e.Ingredients
	}
	if e.Directions != nil {
		r.Directions = *e.Directions
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
