package rating

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	db "github.com/fcmdias/meal/business/db/rating"
	models "github.com/fcmdias/meal/business/models/rating"
	"github.com/go-playground/validator"
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

	var newRating models.NewRating
	if err := json.NewDecoder(r.Body).Decode(&newRating); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error decoding new rating payload: %v", err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(newRating); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error validating new rating: %v", err)
		return
	}

	rating := models.NewRatingToRating(newRating)
	if err := db.Save(b.DB, rating); err != nil {
		b.Log.Println(errors.Wrap(err, "failed to save rating"))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(rating)
	if err != nil {
		b.Log.Println(errors.Wrap(err, "failed to marshal rating JSON"))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}
