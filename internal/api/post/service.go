package post

import "github.com/clevergo/demo/internal/models"

type Service interface {
	Get(id int64) (*models.Post, error)
	Query(limit, offset int) ([]models.Post, error)
	Delete(id int64) error
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) Get(id int64) (*models.Post, error) {
	return nil, nil
}

func (s *service) Query(limit, offset int) ([]models.Post, error) {
	return nil, nil
}

func (s *service) Delete(id int64) error {
	return nil
}
