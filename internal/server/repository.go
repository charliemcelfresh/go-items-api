package server

import (
	"context"

	"github.com/charliemcelfresh/go-items-api/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type repository struct {
	pool *sqlx.DB
}

func NewRepository() repository {
	return repository{pool: config.GetDB()}
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
		    ui.user_id = $1
		LIMIT 10
		OFFSET $2

	`
	err := r.pool.SelectContext(ctx, &itemsToReturn, statement, userID, page)
	return itemsToReturn, err
}
