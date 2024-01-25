package keys

import (
	"crypto/rsa"
	"crypto/x509"
	"embed"
	"encoding/pem"
	"errors"
)

//go:embed *.key
var keysFs embed.FS
var privateKey *rsa.PrivateKey

func GetPrivateKey() (*rsa.PrivateKey, error) {
	if privateKey != nil {
		return privateKey, nil
	}

	bts, err := keysFs.ReadFile("oauth-private.key")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bts)

	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	prvKey, ok := pk.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid private key")
	}

	privateKey = prvKey

	return privateKey, nil
}
