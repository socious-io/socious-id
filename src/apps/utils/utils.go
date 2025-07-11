package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/jmoiron/sqlx/types"
)

const RandCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Copy(src interface{}, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}

func DecodeJWT(token string) (header, payload []byte, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("invalid token format")
	}

	// Decode Header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding header: %v", err)
	}

	// Decode Payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding payload: %v", err)
	}

	return headerBytes, payloadBytes, nil
}

func ArrayContains[T comparable](arr []T, x T) bool {
	for _, item := range arr {
		if item == x {
			return true
		}
	}
	return false
}

func AppendIfNotExists[T comparable](arr []T, x T) []T {
	if !ArrayContains(arr, x) {
		arr = append(arr, x)
	}
	return arr
}

func GenerateRandomDigits(digits int) int {
	if digits < 1 {
		return 0
	}

	lowerBound := int64(1)
	upperBound := int64(10)

	for i := 1; i < digits; i++ {
		lowerBound *= 10
		upperBound *= 10
	}

	n, err := rand.Int(rand.Reader, big.NewInt(upperBound-lowerBound))
	if err != nil {
		// Fallback: return a fixed mid-range number to ensure function always works
		return int((upperBound + lowerBound) / 2)
	}

	return int(n.Int64() + lowerBound)
}

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(RandCharset))))
		if err != nil {
			return RandomString(length)
		}
		b[i] = RandCharset[num.Int64()]
	}
	return string(b)
}

func GenerateChecksum(file io.Reader) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	checksum := hash.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}

func MapToJSONText(m map[string]interface{}) (*types.JSONText, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	jt := types.JSONText(b)

	return &jt, nil
}

func NullableString(s *string) string {
	if s != nil {
		return *s
	}
	return "N/A"
}

func DedentString(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimLeft(line, " \t")
	}
	return strings.Join(lines, "\n")
}
