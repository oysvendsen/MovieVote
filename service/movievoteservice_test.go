package service

import (
	"github.com/theknight1509/MovieVote/domain/movie"
	"github.com/theknight1509/MovieVote/repository"
	"strconv"
	"testing"
)

func Test_fromFile(t *testing.T) {
	repo = repository.FromFile("../repository/test_movies.txt")
	movies := ListMovies()
	assertMoviesContainTitles(t, movies, "Avengers Endgame", "Captain Marvel", "Iron Man")
}

func Test_list(t *testing.T) {
	repo = repository.WithTestData()
	movies := ListMovies()
	assertMoviesContainTitles(t, movies, "Stargate", "Sharknado", "Indiana Jones and Raiders of the Lost Ark")
}

func assertMoviesContainTitles(t *testing.T, movies []movie.Movie, expectedTitles ...string) {
	if (len(movies) != len(expectedTitles)) {
		t.Log("Movie ListMovies is not the expected length")
		t.Fail()
	}
	for _, expectedtitle := range expectedTitles {
		var foundExpectedTitle bool
		for _, movie := range movies {
			if (expectedtitle == movie.Title) {
				foundExpectedTitle = true
				break
			}
		}
		if (!foundExpectedTitle) {
			t.Log("Could not find title " + expectedtitle)
			t.Fail()
		}
	}
}

func Test_Vote_FionaVotesForStargateSuccessfully(t *testing.T) {
	movieId := 1
	voterName := "Fiona"
	existingNumVotes := 1

	repo = repository.WithTestData()
	assertMoviesContainVotes(t, repo.ListMovies(), movieId, existingNumVotes)

	if err := VoteForMovie(movieId, voterName); err != nil {
		t.Errorf("Vote for movie returned %v", err)
	}

	assertMoviesContainVotes(t, repo.ListMovies(), movieId, existingNumVotes + 1)
	if repo.FindVoteByVoterName(voterName) == nil {
		t.Errorf("Vote by %v could not be found", voterName)
	}
}

func Test_VoteForMovie_VoterAlreadyVoted(t *testing.T) {
	movieId := 2
	voterName := "Shayan"
	existingNumVotes := 2

	repo = repository.WithTestData()
	assertMoviesContainVotes(t, repo.ListMovies(), movieId, existingNumVotes)

	err := VoteForMovie(movieId, voterName)

	assertMoviesContainVotes(t, repo.ListMovies(), movieId, existingNumVotes)
	expectedErrorMessage := "Person " + voterName + " has already voted"
	if (err == nil) || (err.Error() != expectedErrorMessage) {
		t.Errorf("With an existing voter, the error should not be %v", err)
	}
}

func Test_VoteForMovie_MovieDoesNotExist(t *testing.T) {
	movieId := 40

	repo = repository.WithTestData()
	if movieById := repo.FindMovieById(movieId); movieById != nil {
		t.Errorf("Movie is not supposed to exist, but is %v", movieById)
	}

	err := VoteForMovie(movieId, "Ardit")

	if (err == nil) || (err.Error() != "Movie with id " + strconv.Itoa(movieId) + " does not exist") {
		t.Errorf("With a non-existent movieId, the error should not be %v", err)
	}
}

func assertMoviesContainVotes(t *testing.T, movies []movie.Movie, id int, expectedNumVotes int) {
	var movieIdFound bool
	for _, movie := range movies {
		if (movie.Id == id) {
			if (movie.NumVotes != expectedNumVotes) {
				t.Errorf("Movie %v contains %v votes, not %v as expected.",
					movie.Title, movie.NumVotes, expectedNumVotes)
			}
			movieIdFound = true
		}
	}
	if !movieIdFound {
		t.Errorf("Movie id %v could not be found in list of movies %v", id, movies)
	}
}
