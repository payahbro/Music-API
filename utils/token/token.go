package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"music-echo/api/domain/dao"
	"time"
)

const ScopeActivation = "activation"

func GenerateToken(userId int64, ttl time.Duration, scope string) (*dao.Token, string, error) {

	// Token
	token := &dao.Token{
		UserId: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Plaintext
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, "", err
	}

	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Hash
	hash := sha256.Sum256([]byte(plaintext))
	token.Hash = hash[:]

	hexString := fmt.Sprintf("%x", token.Hash)
	hexHash := hex.EncodeToString(token.Hash)
	fmt.Println(hexHash, hexString)

	return token, plaintext, nil
}
