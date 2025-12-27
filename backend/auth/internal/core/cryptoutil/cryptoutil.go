package cryptoutil

import (
	"crypto/rand"
	"encoding/base32"
)

func HumanID(prefix string) string {
	b := make([]byte, 4)
	rand.Read(b)
	return prefix + "-" + base32.StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(b)
}
