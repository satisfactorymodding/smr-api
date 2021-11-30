package main

import (
	"crypto/ed25519"
	"encoding/hex"
)

func main() {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)

	if err != nil {
		panic(err)
	}

	println("Public Key:", hex.EncodeToString(publicKey))
	println("Private Key:", hex.EncodeToString(privateKey))
}
