package recipe

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/fcmdias/meal/business/models"

	_ "github.com/lib/pq"
)

// Get this package if it's missing.
// go get -u github.com/lib/pq
// If you encounter problems like I did with a newer version of Go. Run the following:
// GO111MODULE="off" go get -u github.com/lib/pq
// Restart IDE

var recipes = []models.Recipe{
	{
		Name:        "Sweet Potato and Chickpea Curry",
		Description: "A vegan curry recipe made with sweet potatoes, chickpeas, and a variety of spices.",
		Ingredients: []string{
			"1 large sweet potato, peeled and diced",
			"1 can chickpeas, drained and rinsed",
			"1 onion, diced", "2 cloves garlic, minced",
			"1 tablespoon ginger, minced",
			"1 can coconut milk",
			"2 tablespoons tomato paste",
			"1 tablespoon curry powder",
			"1 teaspoon cumin",
			"1 teaspoon coriander",
			"Salt and pepper to taste",
		},
		Directions: []string{
			"Heat a large pot over medium heat. Add the sweet potato and onion and cook until the onion is translucent, about 5 minutes.",
			"Add the garlic and ginger and cook for an additional minute.",
			"Stir in the chickpeas, coconut milk, tomato paste, curry powder, cumin, and coriander. Season with salt and pepper to taste.",
			"Bring the mixture to a boil, then reduce heat to low and let simmer until the sweet potato is tender, about 20 minutes.",
			"Serve hot, garnished with chopped fresh cilantro or parsley if desired.",
		},
	},
	{
		Name:        "Spinach and Feta Stuffed Portobello Mushrooms",
		Description: "A vegetarian dish featuring portobello mushrooms stuffed with spinach, feta cheese, and breadcrumbs.",
		Ingredients: []string{
			"4 large portobello mushrooms",
			"4 cups baby spinach",
			"4 oz crumbled feta cheese",
			"1/2 cup breadcrumbs",
			"1/4 cup grated Parmesan cheese",
			"2 cloves garlic, minced",
			"2 tablespoons olive oil",
			"Salt and pepper to taste",
		},
		Directions: []string{
			"Preheat the oven to 375°F.",
			"Remove the stems from the mushrooms and scrape out the gills with a spoon. Brush the mushrooms with olive oil and season with salt and pepper.",
			"In a large skillet, heat the olive oil over medium heat. Add the garlic and cook until fragrant, about 1 minute.",
			"Add the spinach to the skillet and cook until wilted, about 3 minutes. Remove from heat and stir in the feta cheese and breadcrumbs.",
			"Divide the spinach mixture among the mushrooms, pressing down to fill them completely. Sprinkle the Parmesan cheese over the top.",
			"Bake for 20-25 minutes, or until the mushrooms are tender and the filling is golden brown.",
			"Serve hot, garnished with additional Parmesan cheese and chopped fresh herbs if desired.",
		},
	},
	{
		Name:        "Spaghetti Carbonara",
		Description: "A classic pasta dish made with spaghetti, pancetta, eggs, and Parmesan cheese.",
		Ingredients: []string{"8 oz spaghetti", "4 oz pancetta, diced", "2 eggs", "1/2 cup grated Parmesan cheese", "Salt and pepper to taste"},
		Directions: []string{
			"Bring a large pot of salted water to a boil. Add the spaghetti and cook until al dente, according to package instructions.",
			"While the spaghetti is cooking, heat a large skillet over medium heat. Add the pancetta and cook until crispy, about 8-10 minutes.",
			"In a small bowl, whisk together the eggs and Parmesan cheese.",
			"When the spaghetti is done cooking, reserve 1/2 cup of the pasta water and drain the rest. Add the spaghetti to the skillet with the pancetta and toss to combine. Remove the skillet from heat.",
			"Pour the egg mixture over the spaghetti and toss quickly to coat. Add pasta water as needed to thin the sauce and create a creamy texture.",
			"Season with salt and pepper to taste. Serve hot, garnished with additional Parmesan cheese if desired.",
		},
	},
}

func Get(diet models.DietName) models.Recipe {
	fmt.Println("getting recipe for ", diet)
	switch diet {
	case models.Vegan:
		return recipes[0]
	case models.Vegetarian:
		return recipes[1]
	case models.Omnivore:
		return recipes[2]
	}
	return recipes[0]
}

func Save(db *sql.DB) error {

	fmt.Println("saving new recipe")

	recipe := models.Recipe{
		Name:        "Sweet Potato and Chickpea Curry",
		Description: "A vegan curry recipe made with sweet potatoes, chickpeas, and a variety of spices.",
		Ingredients: []string{
			"1 large sweet potato, peeled and diced",
			"1 can chickpeas, drained and rinsed",
			"1 onion, diced", "2 cloves garlic, minced",
			"1 tablespoon ginger, minced",
			"1 can coconut milk",
			"2 tablespoons tomato paste",
			"1 tablespoon curry powder",
			"1 teaspoon cumin",
			"1 teaspoon coriander",
			"Salt and pepper to taste",
		},
		Directions: []string{
			"Heat a large pot over medium heat. Add the sweet potato and onion and cook until the onion is translucent, about 5 minutes.",
			"Add the garlic and ginger and cook for an additional minute.",
			"Stir in the chickpeas, coconut milk, tomato paste, curry powder, cumin, and coriander. Season with salt and pepper to taste.",
			"Bring the mixture to a boil, then reduce heat to low and let simmer until the sweet potato is tender, about 20 minutes.",
			"Serve hot, garnished with chopped fresh cilantro or parsley if desired.",
		},
	}

	ingredients, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return err
	}
	directions, err := json.Marshal(recipe.Directions)
	if err != nil {
		return err
	}

	sqlStatement := `INSERT INTO recipes (name, description, ingredients, directions) 
	VALUES ($1, $2, $3, $4)`
	fmt.Println(sqlStatement)
	_, err = db.Exec(sqlStatement, recipe.Name,
		recipe.Description,
		ingredients,
		directions)
	if err != nil {
		return err
	}

	return nil

}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS recipes;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE recipes (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		ingredients JSONB,
		directions JSONB
	)`)
	if err != nil {
		return err
	}
	fmt.Println("table recipes is created")
	return nil

}

func GetFirstRecipeFromDB(db *sql.DB) (models.Recipe, error) {
	var recipe models.Recipe
	var ingredientsJson, directionsJson []byte

	row := db.QueryRow("SELECT name, description, ingredients, directions FROM recipes LIMIT 1")
	err := row.Scan(&recipe.Name, &recipe.Description, &ingredientsJson, &directionsJson)
	if err != nil {
		return recipe, err
	}

	err = json.Unmarshal(ingredientsJson, &recipe.Ingredients)
	if err != nil {
		return recipe, err
	}

	err = json.Unmarshal(directionsJson, &recipe.Directions)
	if err != nil {
		return recipe, err
	}

	return recipe, nil
}
