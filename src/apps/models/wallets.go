package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	database "github.com/socious-io/pkg_database"
)

type Wallet struct {
	ID         uuid.UUID `db:"id" json:"id"`
	IdentityID uuid.UUID `db:"identity_id" json:"identity_id"`
	Chain      string    `db:"chain" json:"chain"`
	ChainID    *string   `db:"chain_id" json:"chain_id"`
	Address    string    `db:"address" json:"address"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

func (Wallet) TableName() string {
	return "wallets"
}

func (Wallet) FetchQuery() string {
	return "payments/fetch_wallet"
}

func (w *Wallet) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(w)
}

func (w *Wallet) Upsert(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"payments/upsert_wallet",
		w.IdentityID, w.Chain, w.ChainID, w.Address,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := w.Scan(rows); err != nil {
			return err
		}
	}
	return database.Fetch(w, w.ID)
}

func GetWallets(identityID uuid.UUID) ([]Wallet, error) {
	var (
		wallets   = []Wallet{}
		fetchList []database.FetchList
		ids       []interface{}
	)

	if err := database.QuerySelect("payments/get_wallets", &fetchList, identityID); err != nil {
		return nil, err
	}

	if len(fetchList) < 1 {
		return wallets, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&wallets, ids...); err != nil {
		return nil, err
	}

	return wallets, nil
}

func GetWallet(id uuid.UUID, identityID uuid.UUID) (*Wallet, error) {
	w := new(Wallet)
	if err := database.Get(w, "payments/get_wallet", id, identityID); err != nil {
		return nil, err
	}
	return w, nil
}
