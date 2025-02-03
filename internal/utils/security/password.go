package security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

type IHash interface {
	MakePassword() (string, error)
	ComparePassword() (bool, error)
}

type PasswordHash struct {
	Stored   string
	Supplied string
}

const (
	// CPU/memory cost parameter
	cost = 32768
	// The block mixing parameter.
	// This controls the amount of memory used by
	// the algorithm and the degree of parallelism.
	mixing = 8
	// The Parallelization parameter.
	// This controls the number of independent memory blocks
	// that are processed in parallel.
	Parallelization = 1
	keyLen          = 32
	maxSplit        = 2
	byteSize        = 32
)

var (
	ErrorPasswordHashNotValid   = fmt.Errorf("password hash not valid")
	ErrorPasswordUnableToVerify = fmt.Errorf("unable to verify password")
)

func (h PasswordHash) MakePassword(par int) (string, error) {
	var scryptHash []byte
	var err error

	salt := make([]byte, byteSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	if scryptHash, err = scrypt.Key(
		[]byte(h.Supplied),
		salt, cost, mixing,
		par, keyLen,
	); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s.%s",
		hex.EncodeToString(scryptHash),
		hex.EncodeToString(salt),
	), nil
}

func (h PasswordHash) ComparePassword(par int) (bool, error) {
	var scryptHash []byte
	var salt []byte
	var err error

	pwsalt := strings.Split(h.Stored, ".")
	if len(pwsalt) < maxSplit {
		return false, ErrorPasswordHashNotValid
	}

	if salt, err = hex.DecodeString(pwsalt[1]); err != nil {
		return false, ErrorPasswordUnableToVerify
	}

	if scryptHash, err = scrypt.Key(
		[]byte(h.Supplied),
		salt, cost, mixing,
		par, keyLen,
	); err != nil {
		return false, ErrorPasswordUnableToVerify
	}

	return hex.EncodeToString(scryptHash) == pwsalt[0], nil
}
