package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	db "github.com/fcmdias/meal/business/db/user"
)

type Base struct {
	Log *log.Logger
	DB  *sql.DB
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
