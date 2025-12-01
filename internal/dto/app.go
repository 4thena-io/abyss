package dto

type CreateAppRequestDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Project     string `json:"project"`
}

type CreateAppRecordDto struct {
	Name        string
	Description string
	RepoId      int64
	RepoUrl     string
	CloneUrl    string
	CiId        int64
	CiUrl       string
	Project     string
}
