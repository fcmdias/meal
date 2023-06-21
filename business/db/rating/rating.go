package rating

import (
	"database/sql"
	"fmt"
	"time"

	models "github.com/fcmdias/meal/business/models/rating"
	"github.com/google/uuid"
)

func Save(db *sql.DB, rating models.Rating) error {

	rating.DateCreated = time.Now()
	rating.ID = uuid.New()

	sqlStatement := `INSERT INTO ratings (id, recipeID, userID, score, comments, created) 
	VALUES ($1, $2, $3, $4, $5, $6)`
	fmt.Println(sqlStatement)
	_, err := db.Exec(sqlStatement,
		rating.ID,
		rating.RecipeID,
		rating.UserID,
		rating.Score,
		rating.Comments,
		rating.DateCreated,
	)
	if err != nil {
		return err
	}

	return nil

}
