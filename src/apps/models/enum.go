package models

import (
	"database/sql/driver"
	"fmt"
)

type (
	OTPType                      string
	AuthModeType                 string
	UserVerificationType         string
	OrganizationVerificationType string
	UserStatusType               string
	OrganizationStatusType       string
	CredentialStatusType         string
	ImpactPointType              string
	KybVerificationStatusType    string
	IdentityType                 string
	OauthConnectedProviders      string
	CredentialType               string
)

const (
	//OTPType
	ForgetPasswordOTP OTPType = "FORGET_PASSWORD"
	SSOOTP            OTPType = "SSO"
	VerificationOTP   OTPType = "VERIFICATION"

	//AuthModeType
	AuthModeRegister AuthModeType = "register"
	AuthModeLogin    AuthModeType = "login"

	//UserVerificationType
	UserVerificationTypeEmail    UserVerificationType = "EMAIL"
	UserVerificationTypePhone    UserVerificationType = "PHONE"
	UserVerificationTypeIdentity UserVerificationType = "IDENTITY"

	//OrganizationVerificationType
	OrganizationVerificationTypeNormal OrganizationVerificationType = "NORMAL"
	OrganizationVerificationTypeImpact OrganizationVerificationType = "IMPACT"

	//UserStatusType
	UserStatusTypeActive    UserStatusType = "ACTIVE"
	UserStatusTypeInactive  UserStatusType = "INACTIVE"
	UserStatusTypeSuspended UserStatusType = "SUSPENDED"

	//OrganizationStatusType
	OrganizationStatusTypeActive    OrganizationStatusType = "ACTIVE"
	OrganizationStatusTypeNotActive OrganizationStatusType = "NOT_ACTIVE"
	OrganizationStatusTypeSuspended OrganizationStatusType = "SUSPENDED"
	OrganizationStatusTypePending   OrganizationStatusType = "PENDING"

	//CredentialStatusType
	CredentialStatusCreated   CredentialStatusType = "CREATED"
	CredentialStatusRequested CredentialStatusType = "REQUESTED"
	CredentialStatusVerified  CredentialStatusType = "VERIFIED"
	CredentialStatusFailed    CredentialStatusType = "FAILED"

	//ImpactPointType
	ImpactPointTypeWorkSubmit ImpactPointType = "WORKSUBMIT"
	ImpactPointTypeDonation   ImpactPointType = "DONATION"
	ImpactPointTypeVoting     ImpactPointType = "VOTING"
	ImpactPointTypeVolunteer  ImpactPointType = "VOLUNTEER"
	ImpactPointTypeOther      ImpactPointType = "OTHER"

	//KybVerificationStatusType
	KYBStatusPending  KybVerificationStatusType = "PENDING"
	KYBStatusApproved KybVerificationStatusType = "APPROVED"
	KYBStatusRejected KybVerificationStatusType = "REJECTED"

	//IdentityType
	IdentityTypeUsers         IdentityType = "users"
	IdentityTypeOrganizations IdentityType = "organizations"

	//OauthConnectedProviders
	OauthConnectedProvidersStripe   OauthConnectedProviders = "STRIPE"
	OauthConnectedProvidersStripeJp OauthConnectedProviders = "STRIPE_JP"

	//CredentialType
	CredentialTypeKYC    CredentialType = "KYC"
	CredentialTypeBadges CredentialType = "BADGES"
)

func (pt *OTPType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(pt))
}

func (pt OTPType) Value() (driver.Value, error) {
	return string(pt), nil
}

func (amt *AuthModeType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(amt))
}

func (amt AuthModeType) Value() (driver.Value, error) {
	return string(amt), nil
}

func (uvt *UserVerificationType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(uvt))
}

func (uvt UserVerificationType) Value() (driver.Value, error) {
	return string(uvt), nil
}

func (ovt *OrganizationVerificationType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(ovt))
}

func (ovt OrganizationVerificationType) Value() (driver.Value, error) {
	return string(ovt), nil
}

func (ust *UserStatusType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(ust))
}

func (ust UserStatusType) Value() (driver.Value, error) {
	return string(ust), nil
}

func (ost *OrganizationStatusType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(ost))
}

func (ost OrganizationStatusType) Value() (driver.Value, error) {
	return string(ost), nil
}

func (cst *CredentialStatusType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(cst))
}

func (cst CredentialStatusType) Value() (driver.Value, error) {
	return string(cst), nil
}

func (c *ImpactPointType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(c))
}

func (c ImpactPointType) Value() (driver.Value, error) {
	return string(c), nil
}

func (kvst *KybVerificationStatusType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(kvst))
}

func (kvst KybVerificationStatusType) Value() (driver.Value, error) {
	return string(kvst), nil
}

func (it *IdentityType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(it))
}

func (it IdentityType) Value() (driver.Value, error) {
	return string(it), nil
}

func (ocp *OauthConnectedProviders) Scan(value interface{}) error {
	return scanEnum(value, (*string)(ocp))
}

func (ocp OauthConnectedProviders) Value() (driver.Value, error) {
	return string(ocp), nil
}

func (vt *CredentialType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(vt))
}

func (vt CredentialType) Value() (driver.Value, error) {
	return string(vt), nil
}

func scanEnum(value interface{}, target interface{}) error {
	switch v := value.(type) {
	case []byte:
		*target.(*string) = string(v) // Convert byte slice to string.
	case string:
		*target.(*string) = v // Assign string value.
	default:
		return fmt.Errorf("failed to scan type: %v", value) // Error on unsupported type.
	}
	return nil
}
