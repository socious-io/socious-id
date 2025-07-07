package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"
)

type OauthConnect struct {
	ID             uuid.UUID               `db:"id" json:"id"`
	IdentityID     uuid.UUID               `db:"identity_id" json:"identity_id"`
	Provider       OauthConnectedProviders `db:"provider" json:"provider"`
	Status         UserStatusType          `db:"status" json:"status"`
	MatrixUniqueID string                  `db:"matrix_unique_id" json:"matrix_unique_id"`
	AccessToken    string                  `db:"access_token" json:"access_token"`
	RefreshToken   *string                 `db:"refresh_token" json:"refresh_token"`
	RedirectURL    *string                 `db:"redirect_url" json:"redirect_url"`
	IsConfirmed    bool                    `db:"is_confirmed" json:"is_confirmed"`
	Meta           *types.JSONText         `db:"meta" json:"meta"`
	ExpiredAt      *time.Time              `db:"expired_at" json:"expired_at"`
	CreatedAt      time.Time               `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time               `db:"updated_at" json:"updated_at"`
}

func (OauthConnect) TableName() string {
	return "oauth_connects"
}

func (OauthConnect) FetchQuery() string {
	return "oauth_connects/fetch"
}

func (oc *OauthConnect) Upsert(ctx context.Context) error {
	rows, err := database.Query(
		ctx, "oauth_connects/upsert",
		oc.IdentityID,
		oc.Provider,
		oc.MatrixUniqueID,
		oc.AccessToken,
		oc.RefreshToken,
		oc.Meta,
		oc.ExpiredAt,
		oc.IsConfirmed,
		oc.RedirectURL,
		oc.Status,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(oc); err != nil {
			return err
		}
	}
	return nil
}

func GetOauthConnectByIdentityId(identityID uuid.UUID, provider OauthConnectedProviders) (*OauthConnect, error) {
	oc := new(OauthConnect)
	if err := database.Get(oc, "oauth_connects/get_by_identityid", identityID, provider); err != nil {
		return nil, err
	}
	return oc, nil
}

func GetOauthConnectByMUI(mui string, provider OauthConnectedProviders) (*OauthConnect, error) {
	oc := new(OauthConnect)
	if err := database.Get(oc, "oauth_connects/get_by_mui", mui, provider); err != nil {
		return nil, err
	}
	return oc, nil
}
