package recipe

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
