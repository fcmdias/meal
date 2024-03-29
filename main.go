package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v6"

	rating "github.com/fcmdias/meal/business/handlers/rating"
	meal "github.com/fcmdias/meal/business/handlers/recipe"
	user "github.com/fcmdias/meal/business/handlers/user"
	_ "github.com/lib/pq"
)

func main() {
	log := log.New(os.Stdout, "ADMIN : ", log.LstdFlags|log.Lmicroseconds)

	if err := run(log); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	// =======================================================================
	// Configuration

	cfg := loadConfig()
	b := meal.Base{Log: log}
	ratingB := rating.Base{Log: log}
	userB := user.Base{Log: log}

	// =======================================================================
	// Database Support

	db, err := connectToDatabase(log)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}
	defer db.Close()
	b.DB = db
	ratingB.DB = db
	userB.DB = db

	// =======================================================================
	// App running

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      setupRoutes(b, ratingB, userB),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server started listening on port %s", cfg.Port)
	return server.ListenAndServe()

}

func loadConfig() struct {
	Port string `env:"PORT" envDefault:":8080"`
} {
	cfg := struct {
		Port string `env:"PORT" envDefault:":8080"`
	}{}
	if err := env.Parse(&cfg); err != nil {
		// Handle error if configuration parsing fails
		// For example, log the error and set default values manually
		log.Println("Failed to load configuration:", err)
		cfg.Port = ":8080" // Set default value manually
	}
	return cfg
}

func connectToDatabase(log *log.Logger) (*sql.DB, error) {
	log.Println("Initializing database support")
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"postgres", 5432, "user", "mypassword", "user")
	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	log.Println("Connected to the database")
	return db, nil
}

func setupRoutes(b meal.Base, ratingB rating.Base, userB user.Base) *http.ServeMux {
	mux := http.NewServeMux()

	// =============================================
	// users
	mux.HandleFunc("/users", userB.GetAll)
	mux.HandleFunc("/users/register", userB.Register)
	mux.HandleFunc("/users/login", userB.Login)
	mux.HandleFunc("/users/ping", userB.Ping)

	// =============================================
	// ratings
	mux.HandleFunc("/ratings/save", ratingB.Save)

	// =============================================
	// recipes
	mux.HandleFunc("/recipes", b.GetAll)
	mux.HandleFunc("/recipes/save", b.Save)
	mux.HandleFunc("/recipes/savemany", b.SaveMany)
	mux.HandleFunc("/recipes/get", b.Get)
	mux.HandleFunc("/recipes/recipe-of-the-day", b.RecipeOfTheDayHandler)
	mux.HandleFunc("/recipes/createtable", b.Create)
	mux.HandleFunc("/recipes/update", b.Update)

	return mux
}
