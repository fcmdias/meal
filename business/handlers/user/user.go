package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	db "github.com/fcmdias/meal/business/db/user"
	models "github.com/fcmdias/meal/business/models/user"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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

	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error decoding new user payload: %v", err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(newUser); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error validating new user: %v", err)
		return
	}

	// Encrypt Password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
	if err != nil {
		panic(err)
	}
	newUser.Password = string(passwordHash)

	// userData := models.NewRecipeToRecipe(newUser)
	if err := db.Save(b.DB, newUser); err != nil {
		b.Log.Println(errors.Wrap(err, "failed to save recipe"))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(newUser)
	if err != nil {
		b.Log.Println(errors.Wrap(err, "failed to marshal recipeData JSON"))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

func (b *Base) GetAll(w http.ResponseWriter, r *http.Request) {

	b.Log.Println("user GetAll handler being called")
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := db.GetAllUsersFromDB(b.DB)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error retrieving users: %v", err)
		return
	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error encoding users as JSON: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}
