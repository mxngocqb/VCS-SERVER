package service

type Server struct {
	ID        int    `json:"ID"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
	DeletedAt string `json:"DeletedAt"`
	Name      string `json:"name"`
	Status    bool   `json:"status"`
	IP        string `json:"ip"`
}

type ServersResponse struct {
	Total int      `json:"total"`
	Data  []Server `json:"data"`
}