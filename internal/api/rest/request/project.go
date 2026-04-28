package request

type CreateProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TeamID      *uint  `json:"teamId"`
}
