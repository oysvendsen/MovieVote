package endpointtests

import (
	"testing"
)
const publicKeyUrl = "http://localhost:8080/api/encryption/public"

func TestGetPublicKey(t *testing.T) {
	request := httpGetRequest(t, publicKeyUrl)
	assertRequestIsOk(t, request)
	body := readBody(t, request)
	t.Logf("Successfully recieved public key; '%s'", body)
}

func TestGetPublicKeyTwiceIsSame(t *testing.T) {
	request1 := httpGetRequest(t, publicKeyUrl)
	assertRequestIsOk(t, request1)
	body1 := readBody(t, request1)

	request2 := httpGetRequest(t, publicKeyUrl)
	assertRequestIsOk(t, request2)
	body2 := readBody(t, request2)

	if body1 == body2 {
		t.Logf("Public key is same twice in a row")
	} else {
		t.Fatalf("Public key is not the same twice in a row")
	}
}
