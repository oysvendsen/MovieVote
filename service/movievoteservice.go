package service

import (
	"MovieVote/domain/movie"
	"MovieVote/domain/vote"
	"MovieVote/repository"
	"errors"
	"fmt"
)

var repo repository.MovieVoteRepositoryInMemoryImpl

func Init() {
	repo = repository.WithMovies()
}

func ListMovies() []movie.Movie {
	return repo.ListMovies()
}

func VoteForMovie(id int, voterName string) error {
	if repo.VoterExist(voterName) {
		return errors.New(fmt.Sprintf("Person %v has already voted", voterName))
	}

	storedMovie := repo.FindMovieById(id)
	if storedMovie == nil {
		return errors.New(fmt.Sprintf("Movie with id %v does not exist", id))
	}

	storedMovie.AddVotes(1)
	repo.UpdateMovie(*storedMovie)

	repo.CreateVote(vote.New(voterName, -1))
	return nil
}
