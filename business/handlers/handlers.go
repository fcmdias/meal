package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/fcmdias/meal/business/db/recipe"
	"github.com/fcmdias/meal/business/models"
)

type Base struct {
	Log *log.Logger
	DB  *sql.DB
}

func (b *Base) Save(w http.ResponseWriter, r *http.Request) {
	err := recipe.Save(b.DB)
	if err != nil {
		b.Log.Println(err, "printing")
	}
	b.Log.Println("no errors")
}
func (b *Base) Get(w http.ResponseWriter, r *http.Request) {

	recipe, err := recipe.GetFirstRecipeFromDB(b.DB)
	if err != nil {
		b.Log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(recipe)
	if err != nil {
		log.Printf("Error encoding recipe as JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

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
		var dietname models.DietName
		switch svalues[0] {
		case "vegan":
			data.Diets[0].IsChecked = true
			dietname = models.Vegan
		case "vegetarian":
			data.Diets[1].IsChecked = true
			dietname = models.Vegetarian
		case "omnivore":
			data.Diets[2].IsChecked = true
			dietname = models.Omnivore
		default:
			data.Diets[0].IsChecked = true
			dietname = models.Vegan
		}
		data.Recipe = recipe.Get(dietname)
	} else {
		data.Diets[0].IsChecked = true
		data.Recipe = recipe.Get(models.Vegan)
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

func (b *Base) Create(w http.ResponseWriter, r *http.Request) {
	err := recipe.CreateTable(b.DB)
	if err != nil {
		log.Println(err)
		return
	}
	recipe := models.Recipe{Name: "something"}
	jsonData, err := json.Marshal(recipe)
	if err != nil {
		log.Printf("Error encoding recipe as JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
