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
