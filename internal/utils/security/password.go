package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type IHash interface {
	MakePassword() (string, error)
	ComparePassword() (bool, error)
}

const (
	// CPU/memory cost parameter
	cost = 32768
	// The Parallelization parameter.
	// This controls the number of independent memory blocks
	// that are processed in parallel.
	Parallelization = 1
	keyLen          = 32
	byteSize        = 32
	timeCost        = 1 // Number of iterations (increase for more security)
)

var (
	ErrorPasswordHashNotValid   = fmt.Errorf("password hash not valid")
	ErrorPasswordUnableToVerify = fmt.Errorf("unable to verify password")
)

// why? because scrypt is so fucking slow

type PasswordHashArgon2 struct {
	Stored   string
	Supplied string
}

func (h PasswordHashArgon2) MakePassword() (string, error) {
	salt := make([]byte, byteSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(h.Supplied),
		salt, timeCost, cost, Parallelization, keyLen)
	hashedPassword := fmt.Sprintf("%s.%s",
		base64.StdEncoding.EncodeToString(hash),
		base64.StdEncoding.EncodeToString(salt))
	return hashedPassword, nil
}

func (h PasswordHashArgon2) ComparePassword() (bool, error) {
	parts := strings.Split(h.Stored, ".")
	if len(parts) != 2 {
		return false, ErrorPasswordHashNotValid
	}
	salt, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, ErrorPasswordUnableToVerify
	}
	expectedHash, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, ErrorPasswordUnableToVerify
	}
	computedHash := argon2.IDKey([]byte(h.Supplied), salt, timeCost,
		cost, Parallelization, keyLen)
	return subtleConstantTimeCompare(computedHash, expectedHash), nil
}

// Safe comparison to prevent timing attacks
func subtleConstantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := range a {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
