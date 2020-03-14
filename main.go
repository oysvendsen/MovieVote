package main

import (
	"MovieVote/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const QUERY_PARAM_VOTERNAME = "voterName"
const QUERY_PARAM_ID = "id"

func main() {
	log.Print("Starting server")
	service.Init()
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) { writer.Write([]byte("Hello World")) })
	//http.Handle("/", http.FileServer(http.Dir("./ws-client")))
	http.HandleFunc("/movies/list", listMoviesHttpWrapper)
	http.HandleFunc("/movies/vote", voteMoviesHttpWrapper)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listMoviesHttpWrapper(writer http.ResponseWriter, request *http.Request) {
	log.Println("Recieving list-request")
	log.Println(request.URL)
	log.Println(request.Header)
	log.Println(request.Body)
	movies := service.ListMovies()
	marshal, err := json.MarshalIndent(movies, "", "  ")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	if err == nil {
		writer.Write(marshal)
	} else {
		writer.Write([]byte(fmt.Sprintf("marshalling failed: %v", err)))
	}

	log.Println("Sending list-response")
	log.Println(writer.Header())
	log.Println(movies)
}

func voteMoviesHttpWrapper(writer http.ResponseWriter, request *http.Request) {
	if !validVoteQuery(*request) {
		writer.Write([]byte(fmt.Sprintf("Well this went to hell didn't it! Required params are '%v' and '%v'! Query is %v",
			QUERY_PARAM_ID, QUERY_PARAM_VOTERNAME, request.URL.Query().Encode())))
		return
	}
	queryId, _ := strconv.Atoi(request.URL.Query()[QUERY_PARAM_ID][0])
	queryVoterName := request.URL.Query()[QUERY_PARAM_VOTERNAME][0]
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	err := service.VoteForMovie(queryId, queryVoterName)
	if err == nil {
		writer.Write([]byte(fmt.Sprintf("You just voted! :-)")))
	} else {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
	}
}

func validVoteQuery(request http.Request) bool {
	query := request.URL.Query()
	if query.Get(QUERY_PARAM_ID) == "" {
		return false
	} else {
		if len(query[QUERY_PARAM_ID]) != 1 {
			return false
		}
	}
	if query.Get(QUERY_PARAM_VOTERNAME) == "" {
		return false
	} else {
		if len(query[QUERY_PARAM_VOTERNAME]) != 1 {
			return false
		}
	}
	return true
}
