package transport
import "github.com/mxngocqb/VCS-SERVER/back-end/internal/model"

type ServerResponse struct {
	Total int64          `json:"total"`
	Data  []model.Server `json:"data"`
}
