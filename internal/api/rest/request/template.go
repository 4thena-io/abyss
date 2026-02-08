package request

type CreateTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Language    string `json:"language"`
	RepoURL     string `json:"repoUrl"`
}
