package rating

import (
	"time"

	"github.com/google/uuid"
)

type Rating struct {
	ID          uuid.UUID // Unique identifier for the rating.
	RecipeID    uuid.UUID // Unique identifier for the recipe being rated.
	UserID      uuid.UUID // Unique identifier for the user who provided the rating.
	Score       int       // Numeric rating score (e.g., 1-5).
	Comments    string    // User comments or feedback accompanying the rating.
	DateCreated time.Time // Date and time when the rating was created.
}
