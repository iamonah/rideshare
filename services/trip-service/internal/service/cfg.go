package service

type tripService struct {
	repo TripRepository
}

func NewService(repo TripRepository) *tripService {
	return &tripService{
		repo: repo,
	}
}
