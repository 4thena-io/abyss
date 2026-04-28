package request

type CreateTeam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AddTeamMember struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
