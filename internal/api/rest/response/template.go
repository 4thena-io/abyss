package response

type Template struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Kind            string `json:"kind"`
	Language        string `json:"language"`
	RepoURL         string `json:"repoUrl"`
	CreatedAt       string `json:"createdAt"`
	CreatorID       uint   `json:"creatorId"`
	CreatorUsername string `json:"creatorUsername"`
}
