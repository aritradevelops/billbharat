package cryptoutil

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"math/big"
)

func HumanID(prefix string) string {
	b := make([]byte, 4)
	rand.Read(b)
	return prefix + "-" + base32.StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(b)
}

func GeneratOTP(length int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", length, n), nil
}
