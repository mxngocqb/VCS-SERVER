package transport

// ViewRequest defines the data needed to view a list of servers.
type ViewRequest struct {
	Limit  int    `json:"limit" validate:"required,gte=1" query:"limit"`
	Offset int    `json:"offset" default:"0" validate:"gte=0" query:"offset"`
	Status string `json:"status"  validate:"omitempty,oneof=true false" query:"status"`
	Field  string `json:"field" query:"field"`
	Order  string `json:"order"  validate:"omitempty,oneof=asc desc" query:"order"`
}

// CreateRequest defines the data needed to create a new server. .
type CreateRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status" validate:"required"`
	IP     string `json:"ip" validate:"required"`
}

// UpdateRequest defines the data needed to update a server.
type UpdateRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status" validate:"required"`
	IP     string `json:"ip" validate:"required"`
}

type GetServersReportRequest struct {
	Mail string `json:"mail" validate:"required, email"`
}
