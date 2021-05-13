package vote

type Vote struct {
	VoterName string
	MovieId   int
}

func New(voter string, movieId int) Vote {
	return Vote{
		VoterName: voter,
		MovieId:   movieId,
	}
}