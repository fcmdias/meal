package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "github.com/fcmdias/meal/business/db/user"
	models "github.com/fcmdias/meal/business/models/user"
	"github.com/fcmdias/meal/business/web/auth"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (b *Base) Ping(w http.ResponseWriter, r *http.Request) {
	// Get the token from the Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verify the token
	// claims, err := verifyToken(tokenString)
	// if err != nil {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	// // Pass the username from the token in the request context
	// ctx := r.Context()
	// ctx = context.WithValue(ctx, "username", claims.Username)
	// r = r.WithContext(ctx)

	if auth.ValidateToken(tokenString) != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (b *Base) Register(w http.ResponseWriter, r *http.Request) {
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

	// Check if the username or email already exist in the database

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

	// ========================================================
	// Generate token

	token, err := generateToken(newUser.Email, newUser.Username)
	if err != nil {
		b.Log.Println(errors.Wrap(err, "Failed to generate token"))
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	b.Log.Println("token created", token)

	// ========================================================

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"token": "%s"}`, token)
}

func (b *Base) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData models.User
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error decoding user login payload: %v", err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(userData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error validating user login details: %v", err)
		return
	}

	// check user credentials
	// Check if the username exists in the database
	storedUser, err := db.GetUserByUsername(b.DB, userData.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		log.Printf("Not able to get user from DB: %v", err)
		return
	}

	// ============================================================
	// Compare the provided password with the stored hashed password

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(userData.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// ============================================================
	// create token

	token, err := generateToken(storedUser.Email, storedUser.Password)
	if err != nil {
		b.Log.Println(errors.Wrap(err, "Failed to generate token"))
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	b.Log.Println("token created", token)

	// ============================================================

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"token": "%s"}`, token)
	w.WriteHeader(http.StatusOK)

}
