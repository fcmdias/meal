package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/fcmdias/meal/business/handlers/meal"
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

	// =======================================================================
	// Database Support

	db, err := connectToDatabase(log)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}
	defer db.Close()
	b.DB = db

	// =======================================================================
	// App running

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      setupRoutes(b),
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

func setupRoutes(b meal.Base) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/recipes", b.GetAll)
	mux.HandleFunc("/recipes/create", b.Save)
	mux.HandleFunc("/recipes/savemany", b.SaveMany)
	mux.HandleFunc("/recipes/get", b.Get)
	mux.HandleFunc("/recipes/recipe-of-the-day", b.RecipeOfTheDayHandler)
	mux.HandleFunc("/recipes/createtable", b.Create)
	mux.HandleFunc("/recipes/update", b.Update)
	return mux
}
