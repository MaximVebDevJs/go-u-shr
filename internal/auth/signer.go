package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type Signer struct {
	secret []byte
}

// NewSigner создает signer для симметричной подписи cookie.
func NewSigner(secret string) *Signer {
	return &Signer{
		secret: []byte(secret),
	}
}

// Sign возвращает подписанное значение cookie для userID.
func (s *Signer) Sign(userID string) string {

	mac := hmac.New(
		sha256.New,
		s.secret,
	)

	mac.Write([]byte(userID))

	signature := hex.EncodeToString(
		mac.Sum(nil),
	)

	return userID + "." + signature
}

// Verify проверяет подпись cookie и возвращает исходный userID.
func (s *Signer) Verify(value string) (string, bool) {

	parts := strings.Split(value, ".")

	if len(parts) != 2 {
		return "", false
	}

	userID := parts[0]
	if userID == "" {
		return "", false
	}

	expected := s.Sign(userID)

	if !hmac.Equal(
		[]byte(expected),
		[]byte(value),
	) {
		return "", false
	}

	return userID, true
}
