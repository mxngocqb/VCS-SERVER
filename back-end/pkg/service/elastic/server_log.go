package elastic

import (
	"time"

	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
)

type ElasticService interface {
	IndexServer(server model.Server) error
	DeleteServerFromIndex(id string) error
	LogStatusChange(server model.Server, status bool) error
	CalculateServerUptime(serverID string, date time.Time) (time.Duration, error)
	CalculateServerUptimeFromStartToEnd(serverID string, startDate, endDate time.Time) (time.Duration, error) 
	CreateStatusLogIndex() error
	DeleteServerLogs(serverID string) error
	FetchServersInfo(start, end time.Time) (float64, int, int, int, error)
}
