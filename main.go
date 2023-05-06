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

func main() {

	fmt.Println("service starting")
	http.HandleFunc("/recipe-of-the-day", recipeOfTheDayHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func recipeOfTheDayHandler(w http.ResponseWriter, r *http.Request) {
	recipe := Recipe{
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
	}

	// Define the data context for the template
	data := struct {
		Recipe Recipe
	}{
		Recipe: recipe,
	}

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
