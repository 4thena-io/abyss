package request

type CreateTemplate struct {
	Source      string `json:"source"` // "repo" or "blank"
	Name        string `json:"name"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Language    string `json:"language"`
	RepoURL     string `json:"repoUrl"` // only for source=repo
}

type UpdateTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Language    string `json:"language"`
}
