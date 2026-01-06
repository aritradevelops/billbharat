package cryptoutil

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"math/big"
)

const (
	OTP_LENGTH              = 6
	REFRESH_TOKEN_LENGTH    = 64
	INVITATION_TOKEN_LENGTH = 64
)

func HumanID(prefix string) string {
	b := make([]byte, 4)
	rand.Read(b)
	return prefix + "-" + base32.StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(b)
}

func GenerateOTP() (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(OTP_LENGTH)), nil)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", OTP_LENGTH, n), nil
}

func GenerateRefreshToken() (string, error) {
	return generateHash(REFRESH_TOKEN_LENGTH)
}

func generateHash(length int) (string, error) {
	secretBytes := make([]byte, length)
	_, err := rand.Read(secretBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(secretBytes), nil
}

func GenerateInvitationHash() (string, error) {
	return generateHash(INVITATION_TOKEN_LENGTH)
}
