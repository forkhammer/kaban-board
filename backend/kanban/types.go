package kanban

type UpdateColumnRequest struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
	TeamId *int     `json:"team_id"`
}

type CreateColumnRequest struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
	TeamId *int     `json:"team_id"`
}

type SetColumnOrderRequest struct {
	Id    int `json:"id"`
	Order int `json:"order"`
}

type UpdateTeamRequest struct {
	Title  string `json:"title"`
	Groups []int  `json:"groups"`
}

type CreateTeamRequest struct {
	Title  string `json:"title"`
	Groups []int  `json:"groups"`
}

type KanbanLabel struct {
	Id      string  `json:"id"`
	Title   string  `json:"title"`
	AltName *string `json:"altName"`
}

type UpdateKanbanLabelRequest struct {
	AltName *string `json:"altName"`
}

type SetUserVisibilityRequest struct {
	Visible bool `json:"visible"`
}

type SetTeamRequest struct {
	TeamId *int `json:"team_id"`
}
