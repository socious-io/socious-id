package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"socious-id/src/apps/wallet"
	"socious-id/src/config"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"
)

type VerificationCredential struct {
	ID uuid.UUID `db:"id" json:"id"`

	Status VerificationStatusType `db:"status" json:"status"`

	UserID uuid.UUID `db:"user_id" json:"user_id"`
	User   *User     `db:"-" json:"user"`

	Type VerificationType `db:"type" json:"type"`

	ConnectionID    *string         `db:"connection_id" json:"connection_id"`
	ConnectionURL   *string         `db:"connection_url" json:"connection_url"`
	PresentID       *string         `db:"present_id" json:"present_id"`
	Body            *types.JSONText `db:"body" json:"body"`
	ValidationError *string         `db:"validation_error" json:"validation_error"`

	ConnectionAt *time.Time `db:"connection_at" json:"connection_at"`
	VerifiedAt   *time.Time `db:"verified_at" json:"verified_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`

	UserJson types.JSONText `db:"user" json:"-"`
}

func (VerificationCredential) TableName() string {
	return "verification_credentials"
}

func (VerificationCredential) FetchQuery() string {
	return "verification_credentials/fetch"
}

func (v *VerificationCredential) Create(ctx context.Context, vType VerificationType) error {
	rows, err := database.Query(
		ctx,
		"verification_credentials/create",
		v.UserID, vType,
	)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(v); err != nil {
			return err
		}
	}
	return database.Fetch(v, v.ID)
}

func (v *VerificationCredential) NewConnection(ctx context.Context, callback string) error {
	if v.Status == VerificationStatusRequested {
		return nil
	}
	conn, err := wallet.CreateConnection(callback)
	if err != nil {
		return err
	}
	connectURL, _ := url.JoinPath(config.Config.Host, conn.ShortID)
	rows, err := database.Query(
		ctx,
		"verification_credentials/update_connection",
		v.ID, conn.ID, connectURL,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return database.Fetch(v, v.ID)
}

func (v *VerificationCredential) ProofRequest(ctx context.Context) error {
	if v.ConnectionID == nil {
		return errors.New("connection not valid")
	}
	if time.Since(*v.ConnectionAt) > time.Hour {
		return errors.New("connection expired")
	}

	//Challenge is same as socious work
	presentID, err := wallet.ProofRequest(*v.ConnectionID, "A challenge for the holder to sign")
	if err != nil {
		return err
	}
	rows, err := database.Query(
		ctx,
		"verification_credentials/update_present_id",
		v.ID, presentID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (v *VerificationCredential) HandleByType(ctx context.Context) error {
	switch v.Type {
	case VerificationTypeBadges:
		return v.SendBadges(ctx)
	case VerificationTypeKYC:
		return v.ProofVerify(ctx)
	}

	return v.ProofVerify(ctx)
}

func (v *VerificationCredential) ProofVerify(ctx context.Context) error {
	if v.PresentID == nil {
		return errors.New("need request proof present first")
	}

	vc, err := wallet.ProofVerify(*v.PresentID)
	if err != nil {
		return err
	}
	vcData, _ := json.Marshal(vc)
	duplicateVerification, err := GetSimilar(ctx, v, vc)
	if err == nil && duplicateVerification != nil {
		rows, err := database.Query(
			ctx,
			"verification_credentials/update_present_failed",
			v.ID, vcData, fmt.Sprintf("Duplicate Identity: Verification ID: %s", (*duplicateVerification).ID),
		)
		if err != nil {
			return err
		}
		rows.Close()
		return database.Fetch(v, v.ID)
	}

	rows, err := database.Query(
		ctx,
		"verification_credentials/update_present_verify",
		v.ID, vcData,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return database.Fetch(v, v.ID)
}

func (v *VerificationCredential) SendBadges(ctx context.Context) error {
	if v.Status == VerificationStatusVerified {
		return errors.New("credentials is already sent")
	}

	badges, err := GetImpactBadges(v.UserID)
	if err != nil {
		rows, err := database.Query(
			ctx,
			"verification_credentials/update_present_failed",
			v.ID, nil, fmt.Sprintf("Error Fetching Impact Badge: %s", err.Error()),
		)
		if err != nil {
			return err
		}
		rows.Close()
		return database.Fetch(v, v.ID)
	}

	vc, err := wallet.SendCredentials(*v.ConnectionID, config.Config.Wallet.AgentTrustDID, wallet.H{
		"badges": badges,
	})
	if err != nil {
		rows, err := database.Query(
			ctx,
			"verification_credentials/update_present_failed",
			v.ID, nil, fmt.Sprintf("Error Sending Credentials: %s", err.Error()),
		)
		if err != nil {
			return err
		}
		rows.Close()
		return database.Fetch(v, v.ID)
	}
	vcData, _ := json.Marshal(vc)
	rows, err := database.Query(
		ctx,
		"verification_credentials/update_present_verify",
		v.ID, vcData,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return database.Fetch(v, v.ID)

}

func GetSimilar(ctx context.Context, currentVC *VerificationCredential, data wallet.H) (*VerificationCredential, error) {
	v := new(VerificationCredential)
	err := database.Get(
		v,
		"verification_credentials/get_similar",
		currentVC.ID,
		data["document_number"],
		data["country"],
		data["first_name"],
		data["last_name"],
		data["date_of_birth"],
	)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func GetVerificationByUser(userId uuid.UUID) (*VerificationCredential, error) {
	v := new(VerificationCredential)

	if err := database.Get(v, "verification_credentials/fetch_by_user", userId); err != nil {
		return nil, err
	}
	return v, nil
}

func GetVerificationByUserAndType(userId uuid.UUID, vcType VerificationType) (*VerificationCredential, error) {
	v := new(VerificationCredential)

	if err := database.Get(v, "verification_credentials/fetch_by_user_and_type", userId, vcType); err != nil {
		return nil, err
	}
	return v, nil
}

func GetVerification(id uuid.UUID) (*VerificationCredential, error) {
	v := new(VerificationCredential)

	if err := database.Fetch(v, id); err != nil {
		return nil, err
	}
	return v, nil
}
