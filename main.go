package main

import (
	"encoding/json"
	"fmt"
	"github.com/theknight1509/MovieVote/api/encrypt"
	"github.com/theknight1509/MovieVote/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

/*
DEFINITION OF EACH SINGLE ENDPOINT
*/

type Endpoint struct {
	method      string
	uri         regexp.Regexp
	handlerFunc http.HandlerFunc
}

func listMoviesEndpoint() Endpoint {
	return Endpoint{
		method: "GET",
		uri:    panicWhenError(regexp.Compile("^/movies$")),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			movies := service.ListMovies()
			marshal, err := json.Marshal(movies)
			if err == nil {
				w.Write(marshal)
			} else {
				w.WriteHeader(500)
				log.Println(err.Error())
				w.Write([]byte(err.Error()))
			}
		},
	}
}

func readMoviesEndpoint() Endpoint {
	regexp := panicWhenError(regexp.Compile("^/movies/([0-9a-zA-Z]+)$"))
	return Endpoint{
		method: "GET",
		uri:    regexp,
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			movieId := regexp.FindStringSubmatch(r.URL.Path)[1]
			w.Write([]byte(movieId))
		},
	}
}

func postVoteEndpoint() Endpoint {
	regexp := panicWhenError(regexp.Compile("^/movies/(.+)/vote$"))
	return Endpoint{
		method: "POST",
		uri:    regexp,
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			movieId := regexp.FindStringSubmatch(r.URL.Path)[1]
			movieIdAsInt, err := strconv.Atoi(movieId)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte("InputError - invalid movieId " + movieId))
				return
			}

			var body map[string]interface{}
			err = json.NewDecoder(r.Body).Decode(&body)
			voterName, exists := body["voterName"].(string)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte(fmt.Sprintf("InputError - invalid request %v", body)))
				return
			}
			if !exists {
				w.WriteHeader(400)
				w.Write([]byte("InputError - invalid voterName"))
				return
			}

			err = service.VoteForMovie(movieIdAsInt, voterName)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
			} else {
				w.WriteHeader(201)
			}
		},
	}
}

func getPublicKeyEndpoint() Endpoint {
	return Endpoint{
		method: "GET",
		uri:    panicWhenError(regexp.Compile("api/encryption/public")),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte(encrypt.GlobalInstance.EncodePub()))
		},
	}
}

func postEncryptionValidationEndpoint() Endpoint {
	return Endpoint{
		method: "POST",
		uri:    panicWhenError(regexp.Compile("api/encryption/validation")),
		handlerFunc: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error())) //Are we exposing internals??
				return
			}
			decrypt, err := encrypt.GlobalInstance.Decrypt(string(body))
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error())) //Are we exposing internals??
				return
			}

			if clientPublicKey := r.Header.Get("client-public-key"); clientPublicKey != "" {
				encrypted, err := encrypt.GlobalInstance.Encrypt(decrypt, clientPublicKey)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(err.Error())) //Are we exposing internals??
					return
				}
				w.Write([]byte(encrypted))
				w.WriteHeader(200)
			} else {
				w.Write([]byte(decrypt))
				w.WriteHeader(200)
			}
		},
	}
}

/*
SERVER HANDLER FOR ALL ENDPOINTS
*/

type RestHandler struct {
	endpoints []Endpoint
}

func (h RestHandler) addEndpoints() RestHandler {
	h.endpoints = append(h.endpoints, getPublicKeyEndpoint())
	h.endpoints = append(h.endpoints, postEncryptionValidationEndpoint())
	h.endpoints = append(h.endpoints, listMoviesEndpoint())
	h.endpoints = append(h.endpoints, readMoviesEndpoint())
	h.endpoints = append(h.endpoints, postVoteEndpoint())
	return h
}

func (rh RestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	for _, endpoint := range rh.endpoints {
		if endpoint.method == r.Method && endpoint.uri.MatchString(r.URL.Path) {
			endpoint.handlerFunc.ServeHTTP(w,r)
			return
		}
	}
	w.WriteHeader(404)
}

/*
UTILITIES
*/

func panicWhenError(compile *regexp.Regexp, err error) regexp.Regexp {
	if err != nil {
		panic(err)
	}
	return *compile
}

/*
MAIN METHOD
*/
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %v", port)
	service.Init("movies.txt")
	http.Handle("/", http.FileServer(http.Dir("./ws-client")))
	log.Fatal(http.ListenAndServe(":"+port, new(RestHandler).addEndpoints()))
}
