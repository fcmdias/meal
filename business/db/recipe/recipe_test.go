package recipe

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fcmdias/meal/business/models"
	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

var DateUpdatedTestData = time.Now()
var recipesTestData = []models.Recipe{
	{
		ID:          uuid.New(),
		Name:        "Recipe 1",
		Description: "Description 1",
		Ingredients: []string{"Ingredient 1", "Ingredient 2"},
		Directions:  []string{"Direction 1", "Direction 2"},
		DateCreated: time.Now().Add(-10 * time.Hour),
		DateUpdated: nil,
	},
	{
		ID:          uuid.New(),
		Name:        "Recipe 2",
		Description: "Description 2",
		Ingredients: []string{"Ingredient 3", "Ingredient 4"},
		Directions:  []string{"Direction 3", "Direction 4"},
		DateCreated: time.Now().Add(-10 * time.Hour),
		DateUpdated: &DateUpdatedTestData,
	},
}

func TestGetAllRecipesFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Prepare the mock rows with sample data
	rows := sqlmock.NewRows([]string{"id", "name", "description", "ingredients", "directions", "created", "updated"}).
		AddRow(
			recipesTestData[0].ID,
			"Recipe 1",
			"Description 1",
			`["Ingredient 1", "Ingredient 2"]`,
			`["Direction 1", "Direction 2"]`,
			recipesTestData[0].DateCreated,
			nil,
		).
		AddRow(
			recipesTestData[1].ID,
			"Recipe 2",
			"Description 2",
			`["Ingredient 3", "Ingredient 4"]`,
			`["Direction 3", "Direction 4"]`,
			recipesTestData[1].DateCreated,
			DateUpdatedTestData,
		)

	mock.ExpectQuery("SELECT id, name, description, ingredients, directions, created, updated FROM recipes").
		WillReturnRows(rows)

	recipes, err := GetAllRecipesFromDB(db)
	if err != nil {
		t.Fatalf("error fetching recipes: %s", err)
	}

	assert.Equal(t, recipesTestData, recipes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllRecipesFromDB_EmptyResultSet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, description, ingredients, directions, created, updated FROM recipes").
		WillReturnRows(sqlmock.NewRows([]string{}))

	recipes, err := GetAllRecipesFromDB(db)
	if err != nil {
		t.Fatalf("error fetching recipes: %s", err)
	}

	assert.Empty(t, recipes)
	assert.NoError(t, mock.ExpectationsWereMet())
}
