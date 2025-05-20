package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	database "github.com/socious-io/pkg_database"
)

type Card struct {
	ID         uuid.UUID `db:"id" json:"id"`
	IdentityID uuid.UUID `db:"identity_id" json:"identity_id"`
	Customer   string    `db:"customer" json:"customer"`
	Card       string    `db:"card" json:"card"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

func (Card) TableName() string {
	return "cards"
}

func (Card) FetchQuery() string {
	return "payments/fetch_card"
}

func (c *Card) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(c)
}

func (c *Card) CreateCard(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"payments/create_card",
		c.IdentityID, c.Customer, c.Card,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := c.Scan(rows); err != nil {
			return err
		}
	}
	return database.Fetch(c, c.ID)
}

func (c *Card) UpdateCard(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"payments/update_card",
		c.ID, c.IdentityID, c.Customer, c.Card,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := c.Scan(rows); err != nil {
			return err
		}
	}
	return database.Fetch(c, c.ID)
}

func (c *Card) DeleteCard(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"payments/delete_card",
		c.ID, c.IdentityID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func GetCards(identityID uuid.UUID, p database.Paginate) ([]Card, int, error) {
	var (
		cards     = []Card{}
		fetchList []database.FetchList
		ids       []interface{}
	)

	if err := database.QuerySelect("payments/get_cards", &fetchList, identityID, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return cards, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&cards, ids...); err != nil {
		return nil, 0, err
	}
	return cards, fetchList[0].TotalCount, nil
}

func GetCard(id uuid.UUID, identityID uuid.UUID) (*Card, error) {
	c := new(Card)
	if err := database.Get(c, "get_card", id, identityID); err != nil {
		return nil, err
	}
	return c, nil
}
