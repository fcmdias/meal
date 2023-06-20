package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID // Unique identifier for the user.
	Username            string    // Username of the user.
	Email               string    // Email address of the user.
	Password            string    // Password for the user's account.
	Name                string    // Full name of the user.
	DateOfBirth         time.Time // Date of birth of the user.
	DietaryRestrictions []string  // Dietary restrictions or considerations for the user (e.g., vegan, gluten-free, vegetarian).
}
