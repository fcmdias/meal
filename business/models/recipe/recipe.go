package recipe

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
	DateUpdated         *time.Time    `json:",omitempty"` // Date and time when the recipe was last updated.
	CookingTime         time.Duration `json:",omitempty"` // Estimated time required to cook the recipe.
	PreparationTime     time.Duration `json:",omitempty"` // Estimated time required for the recipe's preparation.
	Difficulty          int           `json:",omitempty"` // Difficulty level of the recipe (1-5, with 1 being the easiest).
	CuisineType         []string      `json:",omitempty"` // Culinary tradition or regional origin of the recipe (e.g., Italian, Mexican, Indian).
	MealType            string        `json:",omitempty"` // Category or type of meal the recipe belongs to (e.g., breakfast, lunch, dinner, starter, dessert, or snack).
	DietaryRestrictions []string      `json:",omitempty"` // Dietary restrictions or considerations for the recipe (e.g., vegan, gluten-free, vegetarian).
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
