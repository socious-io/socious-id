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

type VerificationType string

const (
	VerificationTypeEmail   VerificationType = "EMAIL"
	VerificationTypePhone   VerificationType = "PHONE"
	VerificationTypeIdenity VerificationType = "IDENTITY"
)

func (vt *VerificationType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(vt))
}

func (vt VerificationType) Value() (driver.Value, error) {
	return string(vt), nil
}

type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusInactive  UserStatus = "INACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
)

func (us *UserStatus) Scan(value interface{}) error {
	return scanEnum(value, (*string)(us))
}

func (us UserStatus) Value() (driver.Value, error) {
	return string(us), nil
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
