package recipe

import (
	"time"

	"github.com/google/uuid"
)

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
