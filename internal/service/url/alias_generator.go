package url

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func generateShortID() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, 8)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", fmt.Errorf("сгенерировать alias: %w", err)
		}

		result[i] = letters[n.Int64()]
	}

	return string(result), nil
}
