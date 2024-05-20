package service

type Server struct {
	ID        uint    `json:"ID"`
	Name      string `json:"name"`
	Status    bool   `json:"status"`
	IP        string `json:"ip"`
}

type DropServer struct {
	ID uint `json:"ID"`
	Message string `json:"message"`
}

type ServersResponse struct {
	Total int      `json:"total"`
	Data  []Server `json:"data"`
}