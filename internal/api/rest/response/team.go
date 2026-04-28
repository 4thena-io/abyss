package response

type Team struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	MemberCount  int64  `json:"memberCount"`
	ProjectCount int64  `json:"projectCount"`
	AppCount     int64  `json:"appCount"`
}

type TeamMember struct {
	ID       uint   `json:"id"`
	TeamID   uint   `json:"teamId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	JoinedAt string `json:"joinedAt"`
}
