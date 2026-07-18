package response

type App struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Kind         string `json:"kind"`
	Language     string `json:"language"`
	RepoFullName string `json:"repoFullName"`
	RepoURL      string `json:"repoUrl"`
	CiURL        string `json:"ciUrl"`
	ProjectID       *uint  `json:"projectId,omitempty"`
	ProjectName     string `json:"projectName,omitempty"`
	TemplateID      *uint  `json:"templateId,omitempty"`
	TemplateName    string `json:"templateName,omitempty"`
	CreatorID       uint   `json:"creatorId"`
	CreatorUsername string `json:"creatorUsername"`
}
