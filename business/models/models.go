package models

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
	Directions  []string `json:"directions"`
	DateCreated time.Time
	DateUpdated time.Time
}

type NewRecipe struct {
	ID          uuid.UUID
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
	Directions  []string `json:"directions"`
}

type Diet struct {
	Name       string
	Value      string
	IsDisabled bool
	IsChecked  bool
	Text       string
}

type DietType int

const (
	Vegan DietType = iota
	Vegetarian
	Omnivore
)
