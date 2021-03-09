package genkey

import (
	"crypto/rsa"

	"io/ioutil"
	"log"

	"github.com/dgrijalva/jwt-go"
)

func GenerateRsaKey(prvPath, pubPath string) (*rsa.PrivateKey, *rsa.PublicKey) {
	key, err := ioutil.ReadFile(prvPath)
	if err != nil {
		log.Fatal(err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		log.Fatal(err)
	}
	key, err = ioutil.ReadFile(pubPath)
	if err != nil {
		log.Fatal(err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		log.Fatal(err)
	}
	return privKey, pubKey
}
