package user

import (
	"database/sql"
	"fmt"
	"time"

	models "github.com/fcmdias/meal/business/models/user"
	"github.com/google/uuid"
)

func Save(db *sql.DB, user models.User) error {

	user.DateCreated = time.Now()
	user.ID = uuid.New()

	sqlStatement := `INSERT INTO users (id, username, email, password, name, dateOfBirth, created) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	fmt.Println(sqlStatement)
	_, err := db.Exec(sqlStatement,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.Name,
		user.DateOfBirth,
		user.DateCreated,
	)
	if err != nil {
		return err
	}

	return nil

}

// GetUserByUsername retrieves a user record from the database based on the username
func GetUserByUsername(db *sql.DB, username string) (models.User, error) {
	var user models.User

	sqlStatement := `SELECT id, email, password, name, dateOfBirth, created FROM users WHERE username = $1`
	row := db.QueryRow(sqlStatement, username)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.DateOfBirth,
		&user.DateCreated,
	)
	if err != nil {
		return models.User{}, err
	}

	user.Username = username // Set the username separately as it's not retrieved from the query

	return user, nil
}

func GetAllUsersFromDB(db *sql.DB) ([]models.User, error) {
	var users []models.User

	rows, err := db.Query("SELECT id, name, username, email, dateOfBirth, dietaryRestrictions, created FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		var dietaryRestrictions []byte

		err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.DateOfBirth, &dietaryRestrictions, &user.DateCreated)
		if err != nil {
			return nil, err
		}

		// TODO: handle empty array
		// err = json.Unmarshal(dietaryRestrictions, &user.DietaryRestrictions)
		// if err != nil {
		// 	return nil, err
		// }

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
