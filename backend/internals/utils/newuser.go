package utils

import (
	"context"

	"github.com/arian-nj/chibazi/backend/database"
)

func CreateBrandNewPerson(queries *database.Queries, tgId int, name string) (database.Person, error) {

	tgUserRow, err := queries.CreateTgUser(context.Background(), database.CreateTgUserParams{
		TgID: tgId,
		Name: name,
	})

	if err != nil {
		return tgUserRow, err
	}
	// _, err = app.Queries.InsertUserStatistic(context.Background(), personRow.ID)
	return tgUserRow, err
}
