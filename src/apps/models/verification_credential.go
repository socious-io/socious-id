package models

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"socious-id/src/apps/wallet"
	"socious-id/src/config"
	"time"

	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
)

type VerificationCredential struct {
	ID uuid.UUID `db:"id" json:"id"`

	Status VerificationStatusType `db:"status" json:"status"`

	UserID uuid.UUID `db:"user_id" json:"user_id"`
	User   *User     `db:"-" json:"user"`

	ConnectionID    *string `db:"connection_id" json:"connection_id"`
	ConnectionURL   *string `db:"connection_url" json:"connection_url"`
	PresentID       *string `db:"present_id" json:"present_id"`
	Body            *string `db:"body" json:"body"`
	ValidationError *string `db:"validation_error" json:"validation_error"`

	ConnectionAt *time.Time `db:"connection_at" json:"connection_at"`
	VerifiedAt   *time.Time `db:"verified_at" json:"verified_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

func (VerificationCredential) TableName() string {
	return "verification_credentials"
}

func (VerificationCredential) FetchQuery() string {
	return "verification_credentials/fetch"
}

func (v *VerificationCredential) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"verification_credentials/create",
		v.UserID,
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
	if v.Status == VerificationStatusCreated {
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

	challenge, _ := json.Marshal(wallet.H{
		"type": v.Verification.Schema.Name,
	})

	presentID, err := wallet.ProofRequest(*v.ConnectionID, string(challenge))
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

func (v *VerificationCredential) ProofVerify(ctx context.Context) error {
	if v.PresentID == nil {
		return errors.New("need request proof present first")
	}

	vc, err := wallet.ProofVerify(*v.PresentID)
	if err != nil {
		return err
	}
	vcData, _ := json.Marshal(vc)
	if len(v.Verification.Attributes) > 0 {
		if err := validateVC(*v.Verification.Schema, vc, v.Verification.Attributes); err != nil {
			rows, err := database.Query(
				ctx,
				"verifications/update_present_failed",
				v.ID, vcData, err.Error(),
			)
			if err != nil {
				return err
			}
			rows.Close()
			return nil
		}
	}
	rows, err := database.Query(
		ctx,
		"verifications/update_present_verify",
		v.ID, vcData,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return database.Fetch(v, v.ID)
}

func GetVerifications(id uuid.UUID) (*VerificationCredential, error) {
	v := new(VerificationCredential)

	if err := database.Fetch(v, id); err != nil {
		return nil, err
	}
	verification, err := GetVerification(v.VerificationID)
	if err != nil {
		return nil, err
	}
	v.Verification = verification
	return v, nil
}

func GetVerificationsIndividuals(userId, verificationId uuid.UUID, p database.Paginate) ([]VerificationIndividual, int, error) {
	var (
		verifications = []VerificationIndividual{}
		fetchList     []database.FetchList
		ids           []interface{}
	)

	if err := database.QuerySelect("verifications/get_individuals", &fetchList, userId, verificationId, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return verifications, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&verifications, ids...); err != nil {
		return nil, 0, err
	}
	return verifications, fetchList[0].TotalCount, nil
}
