package request

type CreateApp struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Language    string `json:"language"`
	ProjectID   uint   `json:"projectId"`
	TemplateID  uint   `json:"templateId,omitempty"`
	RepoID      int64  `json:"repoId,omitempty"`
}
