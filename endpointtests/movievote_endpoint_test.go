package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/theknight1509/MovieVote/domain/movie"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestListMovies(t *testing.T) {
	movies := callListMoviesApi(t)
	assertEquals(t, len(movies), 3, "There should be three movies.")
	assertEquals(t, movies[0], movie.New(1,"Stargate"), "First movie should be Stargate.")
	assertEquals(t, movies[1], movie.New(2,"Sharknado"), "Second movie should be Sharknado.")
	assertEquals(t, movies[2], movie.New(3,"Indiana Jones and Raiders of the Lost Ark"), "Third movie should be Indiana Jones.")
}

func TestVoteForMovie(t *testing.T) {
	movieIndexToVoteFor := 0
	moviesBefore := callListMoviesApi(t)
	t.Log("Movies before voting: ", moviesBefore)
	callMovieVoteAssertCreated(t, strconv.Itoa(moviesBefore[movieIndexToVoteFor].Id))
	moviesAfter := callListMoviesApi(t)
	t.Log("Movies after voting: ", moviesAfter)

	assertEquals(t, len(moviesBefore), len(moviesAfter), "Number of movies should not change.")
	for i := 0; i < len(moviesBefore); i++ {
		if i == movieIndexToVoteFor {
			assertEquals(t, moviesBefore[i].NumVotes + 1,
				moviesAfter[i].NumVotes,
				"Number of votes should increase by one.")
		} else {
			assertEquals(t, moviesBefore[i], moviesAfter[i],
				"Movies that were not voted for should not change.")
		}
	}
}

func TestOnlyAllowedToVoteOncePerUser(t *testing.T) {
	movieIndexToVoteFor := 1
	movies := callListMoviesApi(t)
	callMovieVoteAssertCreated(t, strconv.Itoa(movies[movieIndexToVoteFor].Id))
	callMovieVoteApiAssertForbidden(t, strconv.Itoa(movies[movieIndexToVoteFor].Id))
}

func assertEquals(t *testing.T, o1 interface{}, o2 interface{}, description string) {
	if o1 != o2 {
		t.Log(fmt.Sprintf("%s is not equal to %s", o1, o2))
		t.Error(description)
		t.FailNow()
	}
}

func callListMoviesApi(t *testing.T) []movie.Movie {
	t.Log("Calling ListMovies api")
	get, err := http.Get("http://localhost:8080/movies")
	if (err != nil) {
		t.Error(fmt.Sprintf("Failed to get movies with error: %s", err.Error()))
		t.FailNow()
	}

	defer get.Body.Close()
	if get.StatusCode != 200 {
		body, _ := ioutil.ReadAll(get.Body)
		t.Error(fmt.Sprintf("Failed to get movies with error: %s %s", get.Status, body))
		t.FailNow()
	}

	var movies []movie.Movie
	err = json.NewDecoder(get.Body).Decode(&movies)
	if (err != nil) {
		t.Error(fmt.Sprintf("Failed to read responsebody to Movie-struct with error: %s", err.Error()))
		t.FailNow()
	}
	return movies
}
func callMovieVoteAssertCreated(t *testing.T, movieId string) {
	post := callMovieVoteApi(t, movieId)
	if post.StatusCode != 201 {
		body, _ := ioutil.ReadAll(post.Body)
		t.Error(fmt.Sprintf("Failed to post movies with error: %s %s", post.Status, body))
		t.FailNow()
	}
}

func callMovieVoteApiAssertForbidden(t *testing.T, movieId string) {
	post := callMovieVoteApi(t, movieId)
	if post.StatusCode != 500 {
		body, _ := ioutil.ReadAll(post.Body)
		t.Error(fmt.Sprintf("Vote was unexpectedly allowed: %s %s", post.Status, body))
		t.FailNow()
	}
}

func callMovieVoteApi(t *testing.T, movieId string) *http.Response {
	t.Log(fmt.Sprintf("Calling VoteForMovie api with movieId %s and voterName %s", movieId, t.Name()))
	vote := map[string]string{"id": movieId, "voterName": t.Name()}
	json_data, err := json.Marshal(vote)
	if err != nil {
		t.Error(fmt.Sprintf("Failed to parse request to json with error: %s", err.Error()))
		t.FailNow()
	}
	post, err := http.Post(fmt.Sprintf("http://localhost:8080/movies/%s/vote", movieId),
		"application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		t.Error(fmt.Sprintf("Failed to post movies with error: %s", err.Error()))
		t.FailNow()
	}
	return post
}
