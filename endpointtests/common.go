package endpointtests

import (
	"io"
	"net/http"
	"testing"
)

func httpGetRequest(t *testing.T, url string) *http.Response {
	get, err := http.Get(url)
	if err != nil {
		t.Fatalf("Request failed with technical error %s", err.Error())
	}
	return get
}

func assertRequestIsOk(t *testing.T, response *http.Response) {
	if response.StatusCode != 200 {
		t.Fatalf("Request failed with exitcode %s", response.Status)
	}
}

func readBody(t *testing.T, response *http.Response) string {
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("Could not read body")
	}
	if len(body) == 0 {
		t.Fatal("Body is empty")
	}
	return string(body)
}
