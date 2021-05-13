package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
)

type Encryption struct {
	privateKey rsa.PrivateKey
	encodedPublicKey string
}

func New() (Encryption, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return Encryption{}, err
	}
	encodedPublicKey, err := encodePublicKey(key.Public().(*rsa.PublicKey))
	if err != nil {
		return Encryption{}, err
	}
	return Encryption{
		privateKey: *key,
		encodedPublicKey: encodedPublicKey,
	}, nil
}

func (e Encryption) encodePub() string {
	return e.encodedPublicKey
}

func encodePublicKey(pub *rsa.PublicKey) (string,error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(pub)
	if (err != nil) {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(publicKeyBytes), nil
}

func decodePublicKey(encodedPub string) (*rsa.PublicKey,error) {
	decodeString, err := base64.StdEncoding.DecodeString(encodedPub)
	if err != nil {
		return nil, err
	}
	key, err2 := x509.ParsePKIXPublicKey(decodeString)
	if err2 != nil {
		return nil, err2
	}
	return key.(*rsa.PublicKey), nil
}

func (e Encryption) Encrypt(secret string, encodedPublicKey string) (string, error) {
	publicKey, err := decodePublicKey(encodedPublicKey)
	if err != nil {
		return "",nil
	}
	encrypted, err2 := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(secret))
	return string(encrypted), err2
}

func (e Encryption) Decrypt(encrypted string) (string, error) {
	v15, err := rsa.DecryptPKCS1v15(rand.Reader, &e.privateKey, []byte(encrypted))
	return string(v15), err
}
