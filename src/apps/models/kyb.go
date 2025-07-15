package models

import (
	"context"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type KYBDocuments struct {
	Url      string `db:"url" json:"url"`
	Filename string `db:"filename" json:"filename"`
}

type KYBVerification struct {
	ID        uuid.UUID                 `db:"id" json:"id"`
	UserID    uuid.UUID                 `db:"user_id" json:"user_id"`
	OrgID     uuid.UUID                 `db:"organization_id" json:"organization_id"`
	Status    KybVerificationStatusType `db:"status" json:"status"`
	Documents []KYBDocuments            `db:"-" json:"documents"`
	CreatedAt time.Time                 `db:"created_at" json:"created_at"`
	UpdatedAt time.Time                 `db:"updated_at" json:"updated_at"`

	//Json temp fields
	DocumentsJson types.JSONText `db:"documents" json:"-"`
}

func (KYBVerification) TableName() string {
	return "kyb_verifications"
}

func (KYBVerification) FetchQuery() string {
	return "kyb/fetch"
}

func (k *KYBVerification) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(k)
}

func (k *KYBVerification) Create(ctx context.Context, documents []string) error {
	tx, err := database.GetDB().Beginx()
	if err != nil {
		return err
	}
	rows, err := database.TxQuery(ctx, tx, "kyb/create",
		k.UserID, k.OrgID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for rows.Next() {
		if err := rows.StructScan(k); err != nil {
			tx.Rollback()
			return err
		}
	}
	rows.Close()

	rows, err = database.TxQuery(ctx, tx, "kyb/delete_documents",
		k.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()

	for _, document := range documents {
		rows, err = database.TxQuery(ctx, tx, "kyb/create_document",
			k.ID, document,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
		rows.Close()
	}
	tx.Commit()

	return database.Fetch(k, k.ID)
}

func (k *KYBVerification) ChangeStatus(ctx context.Context, status KybVerificationStatusType) error {
	rows, err := database.Query(ctx, "kyb/change_status", k.ID, status)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := k.Scan(rows); err != nil {
			return err
		}
	}
	return database.Fetch(k, k.ID)
}

func GetKyb(id uuid.UUID) (*KYBVerification, error) {
	k := new(KYBVerification)
	if err := database.Fetch(k, id); err != nil {
		return nil, err
	}
	return k, nil
}

func GetKybByOrganization(id uuid.UUID) (*KYBVerification, error) {
	k := new(KYBVerification)
	if err := database.Get(k, "kyb/fetch_by_org", id); err != nil {
		return nil, err
	}
	return k, nil
}
