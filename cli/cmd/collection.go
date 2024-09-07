package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"

	coll "github.com/endalk200/termflow-cli/cmd/collection"
	"github.com/endalk200/termflow-cli/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var collectionCmd = &cobra.Command{
	Use:   "coll",
	Short: "Collection",
	Run: func(cmd *cobra.Command, args []string) {
		collections, _ := coll.FetchCollections()
		for _, collection := range collections {
			fmt.Println(collection.Name)
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add [name] [description]",
	Short: "Add new collection",
	Args:  cobra.ExactArgs(2), // Ensures 2 arguments are passed (name and description)
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		description := args[1]

		collection, err := coll.CreateCollection(coll.CreateCollectionArg{
			Name:        name,
			Description: description,
		})
		if err != nil {
			fmt.Printf("Something went wrong while trying to create a new collection %v\n", err)
		}

		fmt.Println(collection)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your collections",
	Run: func(cmd *cobra.Command, args []string) {
		// number, _ := cmd.Flags().GetInt("number")

		collections, _ := coll.FetchCollections()
		for _, collection := range collections {
			fmt.Println(collection.Name)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete collection with the given id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			fmt.Printf("Invalid id format: %v\n", err)
			return
		}

		_ = coll.DeleteCollection(id)
	},
}

func init() {
	// listCmd.Flags().IntP("number", "n", 10, "Number of collections to list")

	collectionCmd.AddCommand(addCmd)
	collectionCmd.AddCommand(listCmd)
	collectionCmd.AddCommand(deleteCmd)

	rootCmd.AddCommand(collectionCmd)
}
