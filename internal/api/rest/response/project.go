package response

type Project struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	TeamID          *uint  `json:"teamId,omitempty"`
	TeamName        string `json:"teamName,omitempty"`
	CreatorID       uint   `json:"creatorId"`
	CreatorUsername string `json:"creatorUsername"`
}
