package repository

import (
	"MovieVote/domain/movie"
	"MovieVote/domain/vote"
	"errors"
	"io/ioutil"
	"log"
	"strings"
)

type MovieVoteRepository interface {
	ListMovies() []movie.Movie
	FindMovieById(id int) *movie.Movie
	FindVoteByVoterName(name string) *vote.Vote
	CreateVote(v vote.Vote)
	UpdateMovie(updatedMovie movie.Movie)
	VoterExist(name string) bool
}

type MovieVoteRepositoryInMemoryImpl struct {
	movies []movie.Movie
	votes []vote.Vote
}

func (repo MovieVoteRepositoryInMemoryImpl) ListMovies() []movie.Movie{
	movies_copy := make([]movie.Movie, len(repo.movies))
	copy(movies_copy, repo.movies)
	return movies_copy
}

func (repo MovieVoteRepositoryInMemoryImpl) FindMovieById(id int) *movie.Movie {
	var returnMoviePointer *movie.Movie

	for _, movie := range repo.movies {
		if movie.Id == id {
			returnMoviePointer = &movie
			break
		}
	}

	return returnMoviePointer
}

func (repo MovieVoteRepositoryInMemoryImpl) FindVoteByVoterName(name string) *vote.Vote {
	var returnVote *vote.Vote

	for _, vote := range repo.votes {
		if vote.VoterName == name {
			returnVote = &vote
			break
		}
	}

	return returnVote
}

func (repo *MovieVoteRepositoryInMemoryImpl) CreateVote(v vote.Vote) {
	repo.votes = append(repo.votes, vote.New(v.VoterName, v.MovieId))
}

func (repo *MovieVoteRepositoryInMemoryImpl) UpdateMovie(updatedMovie movie.Movie) {
	for i, movie := range repo.movies {
		if movie.Id == updatedMovie.Id {
			repo.movies[i] = updatedMovie
		}
	}
}

func (repo MovieVoteRepositoryInMemoryImpl) VoterExist(name string) bool {
	var voterExist bool
	for _, vote := range repo.votes {
		if vote.VoterName == name {
			voterExist = true
			break
		}
	}
	return voterExist
}

func Empty() MovieVoteRepositoryInMemoryImpl {
	return MovieVoteRepositoryInMemoryImpl{
		movies: nil,
		votes:  nil,
	}
}

func WithMovies() MovieVoteRepositoryInMemoryImpl {
	var movies []movie.Movie
	movies = append(movies, movie.New(1, "Stargate"))
	movies = append(movies, movie.New(2, "Sharknado"))
	movies = append(movies, movie.New(3, "Indiana Jones and Raiders of the Lost Ark"))

	return MovieVoteRepositoryInMemoryImpl{
		movies: movies,
		votes:  nil,
	}
}

func FromFile(filename string) MovieVoteRepositoryInMemoryImpl {
	var movies []movie.Movie

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(errors.New("Cannot read file"))
	}

	fileAsString := string(file)
	fileAsStringLinuxLineEndings := strings.ReplaceAll(fileAsString, "\r\n", "\n")
	fileAsLineArray := strings.Split(fileAsStringLinuxLineEndings, "\n")
	for index, movieTitle := range fileAsLineArray {
		m := movie.New(index, movieTitle)
		log.Printf("movie %d: %v - %v", index, movieTitle, m)
		movies = append(movies, m)
	}

	return MovieVoteRepositoryInMemoryImpl{
		movies: movies,
		votes:  nil,
	}
}

func WithTestData() MovieVoteRepositoryInMemoryImpl {
	var movies []movie.Movie
	movies = append(movies, movie.New(1, "Stargate"))
	movies = append(movies, movie.New(2, "Sharknado"))
	movies = append(movies, movie.New(3, "Indiana Jones and Raiders of the Lost Ark"))

	var votes []vote.Vote
	votes = append(votes, vote.New("Øyvind", 1))
	(&movies[0]).AddVotes(1)
	votes = append(votes, vote.New("Marie", 2))
	votes = append(votes, vote.New("Shayan", 2))
	(&movies[1]).AddVotes(2)
	votes = append(votes, vote.New("Harald", 3))
	votes = append(votes, vote.New("Johannes", 3))
	votes = append(votes, vote.New("Elizabeth", 3))
	(&movies[2]).AddVotes(3)

	return MovieVoteRepositoryInMemoryImpl{
		movies: movies,
		votes:  votes,
	}
}