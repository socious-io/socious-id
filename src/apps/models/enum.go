package models

import (
	"database/sql/driver"
	"fmt"
)

type OTPType string

const (
	ForgetPasswordOTP OTPType = "FORGET_PASSWORD"
	SSOOTP            OTPType = "SSO"
	VerificationOTP   OTPType = "VERIFICATION"
)

func (pt *OTPType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(pt))
}

func (pt OTPType) Value() (driver.Value, error) {
	return string(pt), nil
}

type AuthModeType string

const (
	AuthModeRegister AuthModeType = "register"
	AuthModeLogin    AuthModeType = "login"
)

func (amt *AuthModeType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(amt))
}

func (amt AuthModeType) Value() (driver.Value, error) {
	return string(amt), nil
}

type UserVerificationType string

const (
	UserVerificationTypeEmail   UserVerificationType = "EMAIL"
	UserVerificationTypePhone   UserVerificationType = "PHONE"
	UserVerificationTypeIdenity UserVerificationType = "IDENTITY"
)

func (uvt *UserVerificationType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(uvt))
}

func (uvt UserVerificationType) Value() (driver.Value, error) {
	return string(uvt), nil
}

type OrganizationVerificationType string

const (
	OrganizationVerificationTypeNormal OrganizationVerificationType = "NORMAL"
	OrganizationVerificationTypeImpact OrganizationVerificationType = "IMPACT"
)

func (ovt *OrganizationVerificationType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(ovt))
}

func (ovt OrganizationVerificationType) Value() (driver.Value, error) {
	return string(ovt), nil
}

type UserStatusType string

const (
	UserStatusTypeActive    UserStatusType = "ACTIVE"
	UserStatusTypeInactive  UserStatusType = "INACTIVE"
	UserStatusTypeSuspended UserStatusType = "SUSPENDED"
)

func (ust *UserStatusType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(ust))
}

func (ust UserStatusType) Value() (driver.Value, error) {
	return string(ust), nil
}

type OrganizationStatusType string

const (
	OrganizationStatusTypeActive    OrganizationStatusType = "ACTIVE"
	OrganizationStatusTypeNotActive OrganizationStatusType = "NOT_ACTIVE"
	OrganizationStatusTypeSuspended OrganizationStatusType = "SUSPENDED"
	OrganizationStatusTypePending   OrganizationStatusType = "PENDING"
)

func (ost *OrganizationStatusType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(ost))
}

func (ost OrganizationStatusType) Value() (driver.Value, error) {
	return string(ost), nil
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

type VerificationStatusType string

const (
	VerificationStatusCreated   VerificationStatusType = "CREATED"
	VerificationStatusRequested VerificationStatusType = "REQUESTED"
	VerificationStatusVerified  VerificationStatusType = "VERIFIED"
	VerificationStatusFailed    VerificationStatusType = "FAILED"
)

func (c *VerificationStatusType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*c = VerificationStatusType(string(v))
	case string:
		*c = VerificationStatusType(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (c VerificationStatusType) Value() (driver.Value, error) {
	return string(c), nil
}

type ImpactPointType string

const (
	ImpactPointTypeWorkSubmit ImpactPointType = "WORKSUBMIT"
	ImpactPointTypeDonation   ImpactPointType = "DONATION"
	ImpactPointTypeVolunteer  ImpactPointType = "VOLUNTEER"
	ImpactPointTypeOther      ImpactPointType = "OTHER"
)

func (c *ImpactPointType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*c = ImpactPointType(string(v))
	case string:
		*c = ImpactPointType(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (c ImpactPointType) Value() (driver.Value, error) {
	return string(c), nil
}

type KybVerificationStatusType string

const (
	KYBStatusPending  KybVerificationStatusType = "PENDING"
	KYBStatusApproved KybVerificationStatusType = "APPROVED"
	KYBStatusRejected KybVerificationStatusType = "REJECTED"
)

func (kvst *KybVerificationStatusType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*kvst = KybVerificationStatusType(string(v))
	case string:
		*kvst = KybVerificationStatusType(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (kvst KybVerificationStatusType) Value() (driver.Value, error) {
	return string(kvst), nil
}

type IdentityType string

const (
	IdentityTypeUsers         IdentityType = "users"
	IdentityTypeOrganizations IdentityType = "organizations"
)

func (it *IdentityType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(it))
}

func (it IdentityType) Value() (driver.Value, error) {
	return string(it), nil
}

type OauthConnectedProviders string

const (
	OauthConnectedProvidersStripe   OauthConnectedProviders = "STRIPE"
	OauthConnectedProvidersStripeJp OauthConnectedProviders = "STRIPE_JP"
)

func (ocp *OauthConnectedProviders) Scan(value interface{}) error {
	return scanEnum(value, (*string)(ocp))
}

func (ocp OauthConnectedProviders) Value() (driver.Value, error) {
	return string(ocp), nil
}
