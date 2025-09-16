package schema

// AccessControlServerRequest defines the request body for creating/updating a server.
type AccessControlServerRequest struct {
	Name        string  `json:"name" validate:"required"`
	HostAddress string  `json:"hostAddress" validate:"required"`
	Username    *string `json:"username"`
	Password    *string `json:"password"`
	AccessToken *string `json:"accessToken"`
	ApiToken    *string `json:"apiToken"`
	Status      string  `json:"status"`
}

// AccessControlServerSearchQuery defines the search parameters for servers.
type AccessControlServerSearchQuery struct {
	Name        string `form:"name"`
	HostAddress string `json:"hostAddress"`
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
}

type AccessControlServerInfoResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	HostAddress string `json:"hostAddress"`
}

// AccessControlServerResponse defines the response structure for a server.
type AccessControlServerResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	HostAddress string `json:"hostAddress"`
	Status      string `json:"status"`
}
