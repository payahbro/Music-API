package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Plaintext *string
	Hash      []byte
}

// Set method for hashing plain text
func (p *Password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plainTextPassword
	p.Hash = hash

	return nil
}

func (p *Password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plainTextPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
