package dto

type CreateProjectRequestDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateProjectRecordDto struct {
	Name        string
	Description string
	RepoId      int64
	RepoUrl     string
	CloneUrl    string
	CiId        int64
	CiUrl       string
	Project     string
}
