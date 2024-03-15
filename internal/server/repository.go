package server

import (
	"context"

	"github.com/charliemcelfresh/go-items-api/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	pool *sqlx.DB
}

func NewRepository() repository {
	return repository{pool: config.GetMySQLDB()}
}

// GetItems retrieves the current user's items, offset by the URL "page"
// param, if it exists
func (r repository) GetItems(ctx context.Context, page int) ([]Item, error) {
	itemsToReturn := []Item{}
	userID := getUserIdFromContext(ctx)
	statement := `
		SELECT
			i.id, i.name, i.created_at, i.updated_at
		FROM
			items i
		JOIN
			user_items ui ON ui.item_id = i.id
		WHERE
		    ui.user_id = ?
		LIMIT 10
		OFFSET ?

	`
	err := r.pool.SelectContext(ctx, &itemsToReturn, statement, userID, page)
	return itemsToReturn, err
}
