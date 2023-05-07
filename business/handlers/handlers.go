package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/fcmdias/meal/business/db/recipe"
	"github.com/fcmdias/meal/business/models"
)

type Base struct {
	Log *log.Logger
}

func (b *Base) RecipeOfTheDayHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() //r is url.Values which is a map[string][]string

	var svalues []string
	for _, values := range r.Form { // range over map
		for _, value := range values { // range over []string
			svalues = append(svalues, value) // stick each value in a slice I know the name of
		}
	}
	b.Log.Println(svalues)

	b.Log.Println("recipeOfTheDayHandler handler starting...")

	// Define the data context for the template.
	var data struct {
		models.Recipe
		Diets []models.Diet
	}

	var diets = []models.Diet{
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
		case "vegetarian":
			data.Diets[1].IsChecked = true
		case "omnivore":
			data.Diets[2].IsChecked = true
		default:
			data.Diets[0].IsChecked = true
		}
		data.Recipe = recipe.Get(svalues[0])
	} else {
		data.Diets[0].IsChecked = true
		data.Recipe = recipe.Get("")
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
