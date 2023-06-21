package rating

import "github.com/google/uuid"

type NewRating struct {
	RecipeID uuid.UUID `json:"recipeID" validate:"required"` // Unique identifier for the recipe being rated.
	UserID   uuid.UUID `json:"userID" validate:"required"`   // Unique identifier for the user who provided the rating.
	Score    int       `json:"score" validate:"required"`    // Numeric rating score (e.g., 1-5).
	Comments string    `json:"comments"`                     // User comments or feedback accompanying the rating.
}

func NewRatingToRating(n NewRating) Rating {
	return Rating{
		RecipeID: n.RecipeID,
		UserID:   n.UserID,
		Score:    n.Score,
		Comments: n.Comments,
	}
}
