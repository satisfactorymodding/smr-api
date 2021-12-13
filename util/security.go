package util

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ed25519"
)

var privateKey ed25519.PrivateKey
var publicKey ed25519.PublicKey

var pasetoV2 *paseto.V2

func InitializeSecurity() {
	var err error
	privateKey, err = hex.DecodeString(viper.GetString("paseto.private_key"))
	if err != nil {
		panic(err)
	}

	publicKey, err = hex.DecodeString(viper.GetString("paseto.public_key"))
	if err != nil {
		panic(err)
	}

	pasetoV2 = paseto.NewV2()
}

func GenerateUserToken(userID string) string {
	jsonToken := paseto.JSONToken{
		Expiration: time.Now().Add(time.Hour * 24 * 30),
	}
	jsonToken.Set("userID", userID)

	token, err := pasetoV2.Sign(privateKey, jsonToken, nil)

	if err != nil {
		panic(err)
	}

	return token
}

func VerifyUserToken(token string) (paseto.JSONToken, error) {
	var payload paseto.JSONToken
	err := pasetoV2.Verify(token, publicKey, &payload, nil)
	if err != nil {
		return payload, fmt.Errorf("unable to verify user token: %w", err)
	}
	return payload, nil
}
