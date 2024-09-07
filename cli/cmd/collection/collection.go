package collection

import (
	"context"
	"database/sql"
	"path/filepath"

	"github.com/endalk200/termflow-cli/internal/database"
	"github.com/spf13/viper"
)

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath.Join(viper.GetString("databaseDir"), "data.db"))
	if err != nil {
		return nil, err
	}

	return db, nil
}

type CreateCollectionArg struct {
	Name        string
	Description string
}

func CreateCollection(data CreateCollectionArg) (database.Collection, error) {
	db, err := openDB()
	if err != nil {
		return database.Collection{}, err
	}
	defer db.Close()

	queries := database.New(db)
	ctx := context.Background()

	collection, err := queries.CreateCollection(ctx, database.CreateCollectionParams{
		Name: data.Name,
		Description: sql.NullString{
			String: data.Description,
			Valid:  true,
		},
	})

	return collection, err
}

func FetchCollections() ([]database.Collection, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	queries := database.New(db)
	ctx := context.Background()

	collections, err := queries.ListCollections(ctx)
	if err != nil {
		return nil, err
	}

	return collections, nil
}

func GetCollection(id int64) (database.Collection, error) {
	db, err := openDB()
	defer db.Close()

	queries := database.New(db)
	ctx := context.Background()

	collection, err := queries.GetCollection(ctx, id)
	if err != nil {
		return database.Collection{}, err
	}

	return collection, nil
}

func DeleteCollection(id int64) error {
	db, err := openDB()
	defer db.Close()

	queries := database.New(db)
	ctx := context.Background()

	err = queries.DeleteCollection(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
