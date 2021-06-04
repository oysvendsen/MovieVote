package endpointtests

import (
	"github.com/theknight1509/MovieVote/api/encrypt"
	"testing"
)

const publicKeyUrl = "http://localhost:8080/api/encryption/public"

func TestGetPublicKey(t *testing.T) {
	body := successfullyGetAndReadBody(t, publicKeyUrl)
	t.Logf("Successfully recieved public key; '%s'", body)
}

func TestGetPublicKeyTwiceIsSame(t *testing.T) {
	body1 := successfullyGetAndReadBody(t, publicKeyUrl)
	body2 := successfullyGetAndReadBody(t, publicKeyUrl)

	if body1 == body2 {
		t.Logf("Public key is same twice in a row")
	} else {
		t.Fatalf("Public key is not the same twice in a row")
	}
}

func TestDecryptValidation(t *testing.T) {
	secret := "hello world"
	publicKeyBody := successfullyGetAndReadBody(t, publicKeyUrl)
	encryption, err := encrypt.New()
	if err != nil {
		t.Fatalf("Failed to create Encryption struct; %s", err.Error())
	}
	encryptedSecret, err := encryption.Encrypt(secret, publicKeyBody)
	if err != nil {
		t.Fatalf("Failed to encrypt secret; %s", err.Error())
	}
	request := httpPostRequest(t, "http://localhost:8080/api/encryption/validation", []byte(encryptedSecret))
	assertRequestIsOk(t, request)
	body := readBody(t, request)

	if body == secret {
		t.Logf("Response from validation is same as encrypted payload")
	} else {
		t.Fatalf("Response from validation is not the original secret")
	}
}

func TestDecryptValidationWithPubKeyInHeader(t *testing.T) {
	secret := "hello world"
	publicKeyBody := successfullyGetAndReadBody(t, publicKeyUrl)
	encryption, err := encrypt.New()
	if err != nil {
		t.Fatalf("Failed to create Encryption struct; %s", err.Error())
	}
	encryptedSecret, err := encryption.Encrypt(secret, publicKeyBody)
	if err != nil {
		t.Fatalf("Failed to encrypt secret; %s", err.Error())
	}
	request := httpPostRequestWithHeader(t,
		"http://localhost:8080/api/encryption/validation",
		[]byte(encryptedSecret),
		[]RestHeader{{
			key:   "client-public-key",
			value: encryption.EncodePub(),
		}})
	assertRequestIsOk(t, request)
	body := readBody(t, request)
	decryptedSecret, err := encryption.Decrypt(body)
	if err != nil {
		t.Fatalf("Failed to decrypt responsebody; %s", err.Error())
	}

	if decryptedSecret == secret {
		t.Logf("Response from validation is same as encrypted payload")
	} else {
		t.Fatalf("Response from validation is not the original secret")
	}
}

func successfullyGetAndReadBody(t *testing.T, url string) string {
	request := httpGetRequest(t, publicKeyUrl)
	assertRequestIsOk(t, request)
	return readBody(t, request)
}
