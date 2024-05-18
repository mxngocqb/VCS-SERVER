package service

import (
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/elastic"
)

type Service struct {
	repository *ServerRepository
	elastic    elastic.ElasticService
}

func NewServerService(repository *ServerRepository, elastic elastic.ElasticService) *Service {
	return &Service{
		repository: repository,
		elastic:    elastic,
	}
}

// Update updates a server.
func (s *Service) Update(id string, status bool) error {
	err := s.repository.Update(id, status)
	if err != nil {
		return err
	}

	// Retrieve updated server
	updatedServer, err := s.repository.GetServerByID(id)
	if err != nil {
		return err
	}

	err = s.elastic.LogStatusChange(*updatedServer, status)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTotalServer() (*[]model.Server, error) {
	servers, err := s.repository.GetTotalServer()
	if err != nil {
		return nil, err
	}
	return &servers, nil
}
