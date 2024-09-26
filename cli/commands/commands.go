package commands

import (
	"context"
	"database/sql"
	"errors"

	"github.com/endalk200/termflow-cli/internal/database"
	"github.com/endalk200/termflow-cli/tags"
)

type AddCommandArgs struct {
	ctx         context.Context
	queries     *database.Queries
	command     string
	description string
}

func AddCommands(arg AddCommandArgs) (database.Command, error) {
	newCommand, err := arg.queries.AddCommand(arg.ctx, database.AddCommandParams{
		Command:     sql.NullString{String: arg.command, Valid: true},
		Description: sql.NullString{String: arg.description, Valid: true},
	})
	if err != nil {
		return database.Command{}, err
	}

	return newCommand, nil
}

type AddCommandWithTagsArgs struct {
	Db          *sql.DB
	Ctx         context.Context
	Queries     *database.Queries
	Command     string
	Description string
	Tag         string
}

func AddCommandWithTags(arg AddCommandWithTagsArgs) (database.Command, error) {
	if arg.Command == "" || arg.Tag == "" {
		return database.Command{}, errors.New("Command and tag can not be empty")
	}
	tx, err := arg.Db.Begin()
	if err != nil {
		return database.Command{}, err
	}

	defer tx.Rollback()
	qtx := arg.Queries.WithTx(tx)

	var tagId int64

	tag, err := tags.GetTag(tags.GetTagArgs{
		Ctx:     arg.Ctx,
		Queries: qtx,
		Name:    arg.Tag,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newTag, _ := tags.CreateTag(tags.CreateTagArgs{
				Ctx:         arg.Ctx,
				Queries:     qtx,
				Name:        arg.Tag,
				Description: "",
			})

			tagId = newTag.ID
		}
	} else {
		tagId = tag.ID
	}

	newCommand, err := AddCommands(AddCommandArgs{
		ctx:         arg.Ctx,
		queries:     qtx,
		command:     arg.Description,
		description: arg.Description,
	})
	if err != nil {
		return database.Command{}, err
	}

	err = qtx.AddCommandTag(arg.Ctx, database.AddCommandTagParams{
		Commandid: sql.NullInt64{Int64: newCommand.ID, Valid: true},
		Tagid:     sql.NullInt64{Int64: tagId, Valid: true},
	})
	if err != nil {
		return database.Command{}, err
	}

	err = tx.Commit()
	if err != nil {
		return database.Command{}, err
	}

	return newCommand, nil
}

type GetCommandsForTagArgs struct {
	Ctx     context.Context
	Queries *database.Queries
	Name    string
}

func GetCommandsForTag(args GetCommandsForTagArgs) ([]database.ListCommandsWithTagsByTagNameRow, error) {
	commands, err := args.Queries.ListCommandsWithTagsByTagName(args.Ctx, args.Name)
	if err != nil {
		return []database.ListCommandsWithTagsByTagNameRow{}, err
	}

	return commands, nil
}
