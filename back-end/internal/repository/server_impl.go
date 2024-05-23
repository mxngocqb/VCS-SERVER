package repository

import (
	"time"

	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ServerRepositoryImpl holds DB connection logic
type ServerRepositoryImpl struct {
	DB *gorm.DB
}

// NewServerRepositoryImpl creates a new instance of ServerRepository
func NewServerRepositoryImpl(db *gorm.DB) ServerRepository {
	return &ServerRepositoryImpl{DB: db}
}

// GetServersFiltered retrieves servers with pagination and a status filter
func (ss *ServerRepositoryImpl) GetServersFiltered(perPage, offset int, status, field, order string) ([]model.Server, int, error) {
	var servers []model.Server
	query := ss.DB.Model(&model.Server{})

	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Validate and apply sorting if field and order are provided
	if field != "" && (order == "asc" || order == "desc") {
		if order == "asc" {
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: field}, Desc: false})
		} else {
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: field}, Desc: true})
		}
	} else {
		query = query.Order("created_at desc") // Default sorting
	}

	// Apply pagination
	query = query.Limit(perPage).Offset(offset)

	// Execute query
	if err := query.Find(&servers).Error; err != nil {
		return nil, 0, err
	}

	
	return servers, len(servers), nil
}

// GetServersByOptionalDateRange fetches servers based on optional created and updated date ranges.
func (ss *ServerRepositoryImpl) GetServersByOptionalDateRange(startCreated, endCreated, startUpdated, endUpdated time.Time, field, order string) ([]model.Server, error) {
	query := ss.DB

	if !startCreated.IsZero() && !endCreated.IsZero() {
		query = query.Where("created_at >= ? AND created_at <= ?", startCreated, endCreated)
	}
	if !startUpdated.IsZero() && !endUpdated.IsZero() {
		query = query.Where("updated_at >= ? AND updated_at <= ?", startUpdated, endUpdated)
	}

	// Handle ordering
	if field != "" && (order == "asc" || order == "desc") {
		if order == "asc" {
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: field}, Desc: false})
		} else {
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: field}, Desc: true})
		}
	} else {
		query = query.Order("created_at desc") // Default sorting
	}

	var servers []model.Server
	err := query.Find(&servers).Error
	if err != nil {
		return nil, err
	}
	return servers, nil
}

// Create creates a new server in the database.
func (ss *ServerRepositoryImpl) Create(s *model.Server) error {
	return ss.DB.Create(s).Error
}

// CreateMany creates multiple servers in the database.
func (ss *ServerRepositoryImpl) CreateMany(servers []model.Server) error {
	return ss.DB.Create(&servers).Error
}

// Update updates a server by id
func (ss *ServerRepositoryImpl) Update(id string, s *model.Server) error {
	return ss.DB.Model(&model.Server{}).Where("id = ?", id).Updates(
		map[string]interface{}{"name": s.Name, "status": s.Status, "ip": s.IP}).Error
}

// Delete deletes a server from the database.
func (ss *ServerRepositoryImpl) Delete(id string) error {
	return ss.DB.Delete(&model.Server{}, id).Error
}

// GetServerByID retrieves a server by ID.
func (ss *ServerRepositoryImpl) GetServerByID(id string) (*model.Server, error) {
	var server model.Server
	err := ss.DB.First(&server, id).Error
	return &server, err
}
