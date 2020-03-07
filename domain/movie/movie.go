package movie

type Movie struct {
	Id int
	Title string
	NumVotes int
}

func (movie *Movie) AddVotes(additionalVotes int) {
	movie.NumVotes += additionalVotes;
}

func New(id int, title string) Movie {
	return Movie{
		Id:       id,
		Title:    title,
		NumVotes: 0,
	}
}
