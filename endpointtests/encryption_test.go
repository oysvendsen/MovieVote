package main

import (
	"io"
	"net/http"
	"testing"
)

func TestGetPublicKey(t *testing.T) {
	get, err := http.Get("http://localhost:8080/api/encryption/public")
	if err != nil {
		t.Fatal("Request failed;" + err.Error())
	}
	if get.StatusCode != 200 {
		t.Fatalf("Request failed with exitcode %s", get.Status)
	}
	defer get.Body.Close()
	bytes, err := io.ReadAll(get.Body)
	if err != nil {
		t.Fatal("Could not read body")
	}
	t.Logf("Successfully recieved public key; '%s'", string(bytes))

}
