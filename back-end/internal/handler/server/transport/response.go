package transport
import "github.com/mxngocqb/VCS-SERVER/back-end/internal/model"

// Define the request and response types for the server service
type ServerResponse struct {
	Total int64          `json:"total"`
	Data  []model.Server `json:"data"`
}
