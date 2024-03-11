package secret

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateSecret() (string, error) {
	const op = "lib.secret.GenerateSecret"

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return hex.EncodeToString(randomBytes), nil
}
