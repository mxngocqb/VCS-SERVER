package transport

import (
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
)

// Define the request and response types for the server service
type ServerResponse struct {
	Total int            `json:"total"`
	Data  []model.Server `json:"data"`
}

type ImportServerResponse struct {
	Message 		string 	`json:"message"`
	Total_success 	int 	`json:"total_success"`
	Lists_success 	[]string 	`json:"lists_success"`
	Total_fail 		int 	`json:"total_fail"`
	Lists_fail 		[]string 	`json:"lists_fail"`
}

type ServerStatusResponse struct {
	Online int64 `json:"online"`
	Offline int64 `json:"offline"`
}