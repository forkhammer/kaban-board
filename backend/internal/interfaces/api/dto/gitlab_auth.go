package dto

type GitLabAuthConfigResponse struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url,omitempty"`
}

type GitLabAuthCallbackResponse struct {
	Token string     `json:"token"`
	User  AccountDto `json:"user"`
}
