package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fcmdias/meal/business/handlers"
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

	var cfg struct {
		// todo fix the default configuration
		Port string `conf:"default:8080"`
	}
	if cfg.Port == "" {
		log.Println("defult not set")
		cfg.Port = ":8080"
	}
	b := handlers.Base{Log: log}

	// =======================================================================
	// Database Support

	log.Println("startup initializing database support")
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"postgres", 5432, "user", "mypassword", "user")
	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		return err
	}

	start := time.Now()
	for db.Ping() != nil {
		if start.After(start.Add(10 * time.Second)) {
			return fmt.Errorf("failed to connect after 10 secs.")
		}
	}

	fmt.Println("connected:", db.Ping() == nil)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Println("shutting down stopping database support")
		db.Close()
	}()
	b.DB = db

	// =======================================================================
	// App running

	http.HandleFunc("/recipe-of-the-day", b.RecipeOfTheDayHandler)
	http.HandleFunc("/save", b.Save)
	http.HandleFunc("/get", b.Get)
	http.HandleFunc("/create", b.Create)
	log.Fatal(http.ListenAndServe(cfg.Port, nil))

	return nil
}
