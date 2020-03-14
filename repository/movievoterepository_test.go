package repository_test

import (
	"MovieVote/domain/movie"
	"MovieVote/domain/vote"
	"MovieVote/repository"
	"testing"
)

func Test_EmptyRepository(t *testing.T) {
	filled := repository.Empty()
	if len(filled.ListMovies()) != 0 {
		t.Errorf("Repository is supposed to be empty!")
	}
}

func Test_RepositoryWithMovies(t *testing.T) {
	filled := repository.WithMovies()
	assertMovies(t, filled.ListMovies(),
		movie.Movie{
			Id:       1,
			Title:    "Stargate",
			NumVotes: 0,
		},
		movie.Movie{
			Id:       2,
			Title:    "Sharknado",
			NumVotes: 0,
		},
		movie.Movie{
			Id:       3,
			Title:    "Indiana Jones and Raiders of the Lost Ark",
			NumVotes: 0,
		})
}


func Test_RepositoryFromFile(t *testing.T) {
	filled := repository.FromFile("test_movies.txt")
	assertMovies(t, filled.ListMovies(),
		movie.Movie{
			Id:       0,
			Title:    "Avengers Endgame",
			NumVotes: 0,
		},
		movie.Movie{
			Id:       1,
			Title:    "Captain Marvel",
			NumVotes: 0,
		},
		movie.Movie{
			Id:       2,
			Title:    "Iron Man",
			NumVotes: 0,
		})
}

func Test_RepositoryWithTestData(t *testing.T) {
	filled := repository.WithTestData()
	assertMovies(t, filled.ListMovies(),
		movie.Movie{
			Id:       1,
			Title:    "Stargate",
			NumVotes: 1,
		},
		movie.Movie{
			Id:       2,
			Title:    "Sharknado",
			NumVotes: 2,
		},
		movie.Movie{
			Id:       3,
			Title:    "Indiana Jones and Raiders of the Lost Ark",
			NumVotes: 3,
		})
	//TODO! assert votes
}

func Test_FindFirstMovieById_ReturnsExpectedVoter(t *testing.T) {
	repo := repository.WithTestData()
	movieById := repo.FindMovieById(1)
	expectedMovie := movie.Movie{
		Id:       1,
		Title:    "Stargate",
		NumVotes: 1}
	if (movieById == nil) || (*movieById != expectedMovie) {
		t.Errorf("Did not find expected movie %v, the return value is: %v", expectedMovie, movieById)
	}
}

func Test_FindMovieById_ReturnsNil(t *testing.T) {
	repo := repository.WithTestData()
	movieById := repo.FindMovieById(4)
	if movieById != nil {
		t.Errorf("Movie is supposed to be nil, but is %v", movieById)
	}
}

func Test_FindVoteByVoter_ReturnsExpectedVoter(t *testing.T) {
	repo := repository.WithTestData()
	voteByVoter := repo.FindVoteByVoterName("Elizabeth")
	expectedVoter := vote.Vote{
		VoterName: "Elizabeth",
		MovieId:   3,
	}
	if (voteByVoter == nil) || (*voteByVoter != expectedVoter) {
		t.Errorf("Did not find expected vote %v, the return value is %v", expectedVoter, voteByVoter)
	}
}

func Test_FindVoteByVoter_ReturnsNil(t *testing.T) {
	repo := repository.WithTestData()
	voteByVoter := repo.FindVoteByVoterName("Shrek")
	if voteByVoter != nil {
		t.Errorf("Vote is supposed to be nil, but is %v", *voteByVoter)
	}
}

func Test_VoterExist_ReturnsTrueForExistentVoter(t *testing.T) {
	repo := repository.WithTestData()
	if !repo.VoterExist("Shayan") {
		t.Errorf("Could not find voter 'Shayan'")
	}
}

func Test_VoterExist_ReturnsFalseForNonexistentVoter(t *testing.T) {
	repo := repository.WithTestData()
	if repo.VoterExist("Donkey") {
		t.Errorf("Somehow found voter 'Donkey'")
	}
}

func Test_CreateVote_AddsNewVote(t *testing.T) {
	repo := repository.WithTestData()
	newVote := vote.Vote{
		VoterName: "Shrek",
		MovieId:   3,
	}

	if !repo.VoterExist(newVote.VoterName) {
		repo.CreateVote(newVote)
		if !repo.VoterExist(newVote.VoterName) {
			t.Errorf("Vote was not added")
		}
	} else {
		t.Errorf("Vote already existed")
	}
}

func Test_UpdateMovie_SuccessfullyModifyMovie(t *testing.T) {
	movieId := 2
	newMovieTitle := "Shrek 2"
	newNumVotes := 0

	repo := repository.WithTestData()
	movieById := repo.FindMovieById(movieId)
	if (movieById == nil) ||
		((*movieById).Title == newMovieTitle) ||
		((*movieById).NumVotes == newNumVotes) {
		t.Errorf("The movie already has the expecting values %v", movieById)
	}

	movieById.Title = newMovieTitle
	movieById.NumVotes = newNumVotes
	repo.UpdateMovie(*movieById)

	movieByIdReRead := repo.FindMovieById(movieId)
	if movieByIdReRead == nil {
		t.Errorf("Movie is nil")
	} else if ((*movieByIdReRead).Title != newMovieTitle) &&
		((*movieByIdReRead).NumVotes != newNumVotes) {
		t.Errorf("The movie was not updated, but is %v", movieByIdReRead)
	}
}

func Test_UpdateMovie_GivenMovieWithNonExistentId_UpdatesNothing(t *testing.T) {
	newMovie := movie.Movie{
		Id:       49,
		Title:    "Shrek 3",
		NumVotes: 0,
	}
	repo := repository.WithTestData()
	if movie := repo.FindMovieById(newMovie.Id); movie != nil {
		t.Errorf("Movie %v is not supposed to be present", movie)
	}

	repo.UpdateMovie(newMovie)

	if movie := repo.FindMovieById(newMovie.Id); movie != nil {
		t.Errorf("Movie %v is not supposed to be present", movie)
	}
}
func assertMovies(t *testing.T, movies []movie.Movie, expectedMovies ...movie.Movie) {
	if len(movies) != len(expectedMovies) {
		t.Errorf("Movie list is not the expected length")
	}
	for _, expectedMovie := range expectedMovies {
		var foundExpectedMovie bool
		for _, movie := range movies {
			if expectedMovie == movie {
				foundExpectedMovie = true
				break
			}
		}
		if !foundExpectedMovie {
			t.Errorf("Can not find expected movie %v, in the list of movies \n %v", expectedMovie, movies)
		}
	}
}

func assertVotes(t *testing.T, votes []vote.Vote, expectedVotes ...vote.Vote) {
	if len(votes) != len(expectedVotes) {
		t.Errorf("Movie list is not the expected length")
	}
	for _, expectedVote := range expectedVotes {
		var foundExpectedVote bool
		for _, voteObject := range votes {
			if expectedVote == voteObject {
				foundExpectedVote = true
				break
			}
		}
		if !foundExpectedVote {
			t.Errorf("Can not find expected vote %v, in the list of votes \n %v", expectedVote, votes)
		}
	}
}
