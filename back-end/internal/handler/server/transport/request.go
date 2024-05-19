package transport

// Define the request struct of API get server list
type ViewRequest struct {
	Limit  int    `json:"limit" validate:"required,gte=1" query:"limit"`
	Offset int    `json:"offset" default:"0" validate:"gte=0" query:"offset"`
	Status string `json:"status"  validate:"omitempty,oneof=true false" query:"status"`
	Field  string `json:"field" query:"field"`
	Order  string `json:"order"  validate:"omitempty,oneof=asc desc" query:"order"`
}

// Define the request struct of API create server
type CreateRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status" validate:"required"`
	IP     string `json:"ip" validate:"required"`
}
 
// Define the request struct of API update server
type UpdateRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
	IP     string `json:"ip" validate:"required"`
}