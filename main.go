package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fcmdias/meal/business/handlers"
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
	// App running

	http.HandleFunc("/recipe-of-the-day", b.RecipeOfTheDayHandler)
	log.Fatal(http.ListenAndServe(cfg.Port, nil))

	return nil
}
