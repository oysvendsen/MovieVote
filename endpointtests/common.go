package endpointtests

import (
	"bytes"
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

type RestHeader struct {
	key string
	value string
}

func httpPostRequest(t *testing.T, url string, body []byte) *http.Response {
	return httpPostRequestWithHeader(t, url, body, []RestHeader{})
}

func httpPostRequestWithHeader(t *testing.T, url string, body []byte, restHeaders []RestHeader) *http.Response {
	client := http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Request failed with technical error %s", err.Error())
	}
	for _, header := range restHeaders {
		request.Header.Add(header.key, header.value)
	}
	response, err2 := client.Do(request)
	if err2 != nil {
		t.Fatalf("Request failed with technical error %s", err2.Error())
	}
	return response
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
