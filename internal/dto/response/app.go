package response

type CreateApp struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Language    string `json:"language"`
	RepoURL     string `json:"repoUrl"`
	CiURL       string `json:"ciUrl"`
}
