package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

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

func main() {

	fmt.Println("service starting")
	http.HandleFunc("/recipe-of-the-day", recipeOfTheDayHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func recipeOfTheDayHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() //r is url.Values which is a map[string][]string

	var svalues []string
	for _, values := range r.Form { // range over map
		for _, value := range values { // range over []string
			svalues = append(svalues, value) // stick each value in a slice I know the name of
		}
	}
	fmt.Println(svalues)

	fmt.Println("recipeOfTheDayHandler handler starting...")
	recipes := []Recipe{
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
	}

	// Define the data context for the template.
	var data struct {
		Recipe
		Diets []Diet
	}

	var diets = []Diet{
		{
			Name:       "diet",
			Value:      "vegan",
			IsDisabled: false,
			IsChecked:  false,
			Text:       "Vegan",
		},
		{
			Name:       "diet",
			Value:      "vegetarian",
			IsDisabled: false,
			IsChecked:  false,
			Text:       "Vegetarian",
		},
		{
			Name:       "diet",
			Value:      "omnivore",
			IsDisabled: false,
			IsChecked:  false,
			Text:       "Omnivore",
		},
	}

	data.Diets = diets

	if len(svalues) > 0 {
		switch svalues[0] {
		case "vegan":
			data.Diets[0].IsChecked = true
			data.Recipe = recipes[1]
		case "vegetarian":
			data.Diets[1].IsChecked = true
			data.Recipe = recipes[2]
		case "omnivore":
			data.Diets[2].IsChecked = true
			data.Recipe = recipes[0]
		default:
			data.Diets[0].IsChecked = true
			data.Recipe = recipes[1]
		}
	}

	fmt.Println(data)

	// Parse the template file
	tmpl, err := template.ParseFiles("templates/recipe_of_the_day.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template with the data context
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
