package request

type CreateTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Language    string `json:"language"`
	RepoUrl     string `json:"repo_url"`
	CloneUrl    string `json:"clone_url"`
}
