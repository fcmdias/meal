package models

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID                  uuid.UUID     // Unique identifier for the recipe.
	Name                string        // Name or title of the recipe.
	Description         string        // Brief description or summary of the recipe.
	Ingredients         []string      // List of ingredients required for the recipe.
	Directions          []string      // Step-by-step instructions for preparing the recipe.
	DateCreated         time.Time     // Date and time when the recipe was created.
	DateUpdated         *time.Time    // Date and time when the recipe was last updated.
	CookingTime         time.Duration // Estimated time required to cook the recipe.
	PreparationTime     time.Duration // Estimated time required for the recipe's preparation.
	Difficulty          int           // Difficulty level of the recipe (1-5, with 1 being the easiest).
	CuisineType         []string      // Culinary tradition or regional origin of the recipe (e.g., Italian, Mexican, Indian).
	MealType            string        // Category or type of meal the recipe belongs to (e.g., breakfast, lunch, dinner, starter, dessert, or snack).
	DietaryRestrictions []string      // Dietary restrictions or considerations for the recipe (e.g., vegan, gluten-free, vegetarian).
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
