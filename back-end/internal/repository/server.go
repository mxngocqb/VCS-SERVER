package repository

import (
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"time"
)

type ServerRepository interface {
	GetServersFiltered(perPage, offset int, status, field, order string) ([]model.Server, int, error)
	GetServersByOptionalDateRange(startCreated, endCreated, startUpdated, endUpdated time.Time, field, order string) ([]model.Server, error)
	Create(s *model.Server) error
	CreateMany(servers []model.Server) error
	Update(id string, s *model.Server) error
	Delete(id string) error
	GetServerByID(id string) (*model.Server, error)
}