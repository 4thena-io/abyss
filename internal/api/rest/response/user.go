package response

type UserProfile struct {
	ID           uint              `json:"id"`
	Username     string            `json:"username"`
	Email        string            `json:"email"`
	AvatarURL    string            `json:"avatarUrl"`
	IsAdmin      bool              `json:"isAdmin"`
	TeamCount    int64             `json:"teamCount"`
	ProjectCount int64             `json:"projectCount"`
	AppCount     int64             `json:"appCount"`
	Teams        []UserTeamMember  `json:"teams"`
}

type UserTeamMember struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	MemberCount  int64  `json:"memberCount"`
	ProjectCount int64  `json:"projectCount"`
}
