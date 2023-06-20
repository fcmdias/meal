package meal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	db "github.com/fcmdias/meal/business/db/recipe"
	models "github.com/fcmdias/meal/business/models/recipe"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Base struct {
	Log *log.Logger
	DB  *sql.DB
}

func (b *Base) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var newRecipeData models.NewRecipe
	if err := json.NewDecoder(r.Body).Decode(&newRecipeData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error decoding recipe payload: %v", err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(newRecipeData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error validating new recipe: %v", err)
		return
	}

	recipeData := models.NewRecipeToRecipe(newRecipeData)
	if err := db.Save(b.DB, recipeData); err != nil {
		b.Log.Println(errors.Wrap(err, "failed to save recipe"))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(recipeData)
	if err != nil {
		b.Log.Println(errors.Wrap(err, "failed to marshal recipeData JSON"))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

func (b *Base) SaveMany(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var newRecipesData []models.NewRecipe
	if err := json.NewDecoder(r.Body).Decode(&newRecipesData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error decoding recipe payload: %v", err)
		return
	}

	validate := validator.New()
	for _, newRecipeData := range newRecipesData {
		if err := validate.Struct(newRecipeData); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Printf("Error validating new recipe: %v", err)
			return
		}
	}

	recipes := make([]models.Recipe, len(newRecipesData))
	for i, newRecipeData := range newRecipesData {
		recipes[i] = models.NewRecipeToRecipe(newRecipeData)
	}

	if err := db.SaveMany(b.DB, recipes); err != nil {
		b.Log.Println(errors.Wrap(err, "failed to save recipes"))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(recipes)
	if err != nil {
		b.Log.Println(errors.Wrap(err, "failed to marshal recipes JSON"))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

func (b *Base) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	recipeIDStr := r.URL.Query().Get("id")
	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Invalid recipe ID: %v", err)
		return
	}

	var RecipeEdit models.RecipeEdit
	if err := json.NewDecoder(r.Body).Decode(&RecipeEdit); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error decoding recipe payload: %v", err)
		return
	}

	// validate := validator.New()
	// if err := validate.Struct(RecipeEdit); err != nil {
	// 	http.Error(w, "Bad Request", http.StatusBadRequest)
	// 	log.Printf("Error validating new recipe: %v", err)
	// 	return
	// }

	recipe, err := db.GetRecipeByIDFromDB(b.DB, recipeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusMethodNotAllowed)
		return
	}
	models.RecipeEditToRecipe(recipe, RecipeEdit)

	if err := db.Update(b.DB, *recipe); err != nil {
		b.Log.Println(errors.Wrap(err, "failed to update recipe"))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(recipe)
	if err != nil {
		b.Log.Println(errors.Wrap(err, "failed to marshal recipe JSON"))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func (b *Base) GetAll(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	recipes, err := db.GetAllRecipesFromDB(b.DB)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error retrieving recipes: %v", err)
		return
	}

	jsonData, err := json.Marshal(recipes)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error encoding recipes as JSON: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

// Get retrieves a recipe by ID and returns it as JSON.
// It expects the recipe ID to be passed as a query parameter "id" in the URL.
// If the recipe is found, it returns a JSON response with the recipe data.
// If the recipe ID is not provided or is invalid, it returns a "Bad Request" response.
// If the recipe is not found, it returns a "Not Found" response.
// If there is an error retrieving or encoding the recipe, it returns an "Internal Server Error" response.
func (b *Base) Get(w http.ResponseWriter, r *http.Request) {
	recipeIDStr := r.URL.Query().Get("id")
	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Invalid recipe ID: %v", err)
		return
	}

	log.Println("id: ", recipeID.String())
	recipe, err := db.GetRecipeByIDFromDB(b.DB, recipeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error retrieving recipe: %v", err)
		return
	}

	if recipe == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(recipe)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error encoding recipe as JSON: %v", err)
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
		var dietname models.DietType
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
		data.Recipe = db.Get(dietname)
	} else {
		data.Diets[0].IsChecked = true
		data.Recipe = db.Get(models.Vegan)
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
	err := db.CreateTable(b.DB)
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
