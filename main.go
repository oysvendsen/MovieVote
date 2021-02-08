package main

import (
	"encoding/json"
	"fmt"
	"github.com/theknight1509/MovieVote/service"
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %v", port)
	service.Init("movies.txt")
	http.Handle("/", http.FileServer(http.Dir("./ws-client")))
	handler := MovieVoteRequestHandler{}
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
