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

const QUERY_PARAM_VOTERNAME = "voterName"
const QUERY_PARAM_ID = "id"

type MovieVoteRequestHandler struct{}

func (handler MovieVoteRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("Serving %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr))
	listMoviesRegexp, _ := regexp.Compile("^/movies$")
	readMovieRegexp, _ := regexp.Compile("^/movies/([0-9a-zA-Z]+)$")
	voteMovieRegexp, _ := regexp.Compile("^/movies/(.+)/vote$")
	switch {
	case (r.Method == "GET") && listMoviesRegexp.Match([]byte(r.URL.Path)):
		movies := service.ListMovies()
		marshal, err := json.Marshal(movies)
		if err == nil {
			w.Write(marshal)
		} else {
			w.WriteHeader(500)
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
		}
	case (r.Method == "GET") && readMovieRegexp.Match([]byte(r.URL.Path)):
		movieId := readMovieRegexp.FindStringSubmatch(r.URL.Path)[1]
		w.Write([]byte(movieId))
	case (r.Method == "POST") && voteMovieRegexp.Match([]byte(r.URL.Path)):
		movieId := voteMovieRegexp.FindStringSubmatch(r.URL.Path)[1]
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
	default:
		w.WriteHeader(404)
	}
}

func regexpOrPanic(expr string) *regexp.Regexp {
	compile, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}
	return compile
}

/*
func ifErrWriteError(w http.ResponseWriter, err error, optionalErrorMessage string) bool{
	if err != nil {
		w.WriteHeader(500)
		if len(optionalErrorMessage) > 0 {
			w.Write([]byte(optionalErrorMessage))
		} else {
			w.Write([]byte(err.Error()))
		}
		return true
	}
	return false
}
*/

type Endpoint struct {
	method string
	uri regexp.Regexp
	handler http.HandlerFunc
}

func GetPublicKeyEndpoint() Endpoint {
	return Endpoint{
		method: "GET",
		uri:     *regexpOrPanic("api/encryption/public"),
		handler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte(encrypt.GlobalInstance.EncodePub()))
		},
	}
}

func GetEncryptionValidationEndpoint() Endpoint {
	return Endpoint{
		method: "POST",
		uri:     *regexpOrPanic("api/encryption/validation"),
		handler: func(w http.ResponseWriter, r *http.Request) {
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
			w.Write([]byte(decrypt))
			w.WriteHeader(200)
		},
	}
}

type RestHandler struct {
	endpoints []Endpoint
}

func (h RestHandler) addEndpoints() RestHandler {
	h.endpoints = append(h.endpoints, GetPublicKeyEndpoint())
	h.endpoints = append(h.endpoints, GetEncryptionValidationEndpoint())
	return h
}

func (rh RestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	for _, endpoint := range rh.endpoints {
		if endpoint.method == r.Method && endpoint.uri.MatchString(r.URL.Path) {
			endpoint.handler.ServeHTTP(w,r)
			return
		}
	}
	w.WriteHeader(404)
}

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
