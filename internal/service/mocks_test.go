package service

import (
	"context"
	"fmt"

	"github.com/4thena-io/abyss/internal/model"
)

// mockAppRepository implements AppRepository. Each method delegates to an
// optional func field; calling a method whose field is nil panics so a test
// that hits an unexpected call fails loudly instead of silently zero-valuing.
type mockAppRepository struct {
	SaveFn          func(ctx context.Context, app *model.App) error
	GetAllFn        func(ctx context.Context) ([]model.App, error)
	GetByIDFn       func(ctx context.Context, id uint) (*model.App, error)
	GetByNameFn     func(ctx context.Context, name string) (*model.App, error)
	GetByProjectFn  func(ctx context.Context, id uint) ([]model.App, error)
	GetByTemplateFn func(ctx context.Context, id uint) ([]model.App, error)
	UpdateFn        func(ctx context.Context, app *model.App) error
	DeleteFn        func(ctx context.Context, app *model.App) error
	CountByTeamFn   func(ctx context.Context, teamID uint) (int64, error)
}

func (m *mockAppRepository) Save(ctx context.Context, app *model.App) error {
	if m.SaveFn == nil {
		panic("mockAppRepository.Save not implemented")
	}
	return m.SaveFn(ctx, app)
}

func (m *mockAppRepository) GetAll(ctx context.Context) ([]model.App, error) {
	if m.GetAllFn == nil {
		panic("mockAppRepository.GetAll not implemented")
	}
	return m.GetAllFn(ctx)
}

func (m *mockAppRepository) GetByID(ctx context.Context, id uint) (*model.App, error) {
	if m.GetByIDFn == nil {
		panic("mockAppRepository.GetByID not implemented")
	}
	return m.GetByIDFn(ctx, id)
}

func (m *mockAppRepository) GetByName(ctx context.Context, name string) (*model.App, error) {
	if m.GetByNameFn == nil {
		panic("mockAppRepository.GetByName not implemented")
	}
	return m.GetByNameFn(ctx, name)
}

func (m *mockAppRepository) GetByProject(ctx context.Context, id uint) ([]model.App, error) {
	if m.GetByProjectFn == nil {
		panic("mockAppRepository.GetByProject not implemented")
	}
	return m.GetByProjectFn(ctx, id)
}

func (m *mockAppRepository) GetByTemplate(ctx context.Context, id uint) ([]model.App, error) {
	if m.GetByTemplateFn == nil {
		panic("mockAppRepository.GetByTemplate not implemented")
	}
	return m.GetByTemplateFn(ctx, id)
}

func (m *mockAppRepository) Update(ctx context.Context, app *model.App) error {
	if m.UpdateFn == nil {
		panic("mockAppRepository.Update not implemented")
	}
	return m.UpdateFn(ctx, app)
}

func (m *mockAppRepository) Delete(ctx context.Context, app *model.App) error {
	if m.DeleteFn == nil {
		panic("mockAppRepository.Delete not implemented")
	}
	return m.DeleteFn(ctx, app)
}

func (m *mockAppRepository) CountByTeam(ctx context.Context, teamID uint) (int64, error) {
	if m.CountByTeamFn == nil {
		panic("mockAppRepository.CountByTeam not implemented")
	}
	return m.CountByTeamFn(ctx, teamID)
}

// mockProjectRepository implements ProjectRepository the same way.
type mockProjectRepository struct {
	SaveFn        func(ctx context.Context, project *model.Project) error
	UpdateFn      func(ctx context.Context, project *model.Project) error
	GetAllFn      func(ctx context.Context) ([]model.Project, error)
	GetByIDFn     func(ctx context.Context, id uint) (*model.Project, error)
	GetByNameFn   func(ctx context.Context, name string) (*model.Project, error)
	GetByTeamFn   func(ctx context.Context, teamID uint) ([]model.Project, error)
	DeleteFn      func(ctx context.Context, project *model.Project) error
	CountByTeamFn func(ctx context.Context, teamID uint) (int64, error)
}

func (m *mockProjectRepository) Save(ctx context.Context, project *model.Project) error {
	if m.SaveFn == nil {
		panic("mockProjectRepository.Save not implemented")
	}
	return m.SaveFn(ctx, project)
}

func (m *mockProjectRepository) Update(ctx context.Context, project *model.Project) error {
	if m.UpdateFn == nil {
		panic("mockProjectRepository.Update not implemented")
	}
	return m.UpdateFn(ctx, project)
}

func (m *mockProjectRepository) GetAll(ctx context.Context) ([]model.Project, error) {
	if m.GetAllFn == nil {
		panic("mockProjectRepository.GetAll not implemented")
	}
	return m.GetAllFn(ctx)
}

func (m *mockProjectRepository) GetByID(ctx context.Context, id uint) (*model.Project, error) {
	if m.GetByIDFn == nil {
		panic("mockProjectRepository.GetByID not implemented")
	}
	return m.GetByIDFn(ctx, id)
}

func (m *mockProjectRepository) GetByName(ctx context.Context, name string) (*model.Project, error) {
	if m.GetByNameFn == nil {
		panic("mockProjectRepository.GetByName not implemented")
	}
	return m.GetByNameFn(ctx, name)
}

func (m *mockProjectRepository) GetByTeam(ctx context.Context, teamID uint) ([]model.Project, error) {
	if m.GetByTeamFn == nil {
		panic("mockProjectRepository.GetByTeam not implemented")
	}
	return m.GetByTeamFn(ctx, teamID)
}

func (m *mockProjectRepository) Delete(ctx context.Context, project *model.Project) error {
	if m.DeleteFn == nil {
		panic("mockProjectRepository.Delete not implemented")
	}
	return m.DeleteFn(ctx, project)
}

func (m *mockProjectRepository) CountByTeam(ctx context.Context, teamID uint) (int64, error) {
	if m.CountByTeamFn == nil {
		panic("mockProjectRepository.CountByTeam not implemented")
	}
	return m.CountByTeamFn(ctx, teamID)
}

// mockTemplateRepository implements TemplateRepository.
type mockTemplateRepository struct {
	SaveFn      func(ctx context.Context, template *model.Template) error
	UpdateFn    func(ctx context.Context, template *model.Template) error
	GetAllFn    func(ctx context.Context) ([]model.Template, error)
	GetByIDFn   func(ctx context.Context, id uint) (*model.Template, error)
	GetByNameFn func(ctx context.Context, name string) (*model.Template, error)
	DeleteFn    func(ctx context.Context, template *model.Template) error
}

func (m *mockTemplateRepository) Save(ctx context.Context, template *model.Template) error {
	if m.SaveFn == nil {
		panic("mockTemplateRepository.Save not implemented")
	}
	return m.SaveFn(ctx, template)
}

func (m *mockTemplateRepository) Update(ctx context.Context, template *model.Template) error {
	if m.UpdateFn == nil {
		panic("mockTemplateRepository.Update not implemented")
	}
	return m.UpdateFn(ctx, template)
}

func (m *mockTemplateRepository) GetAll(ctx context.Context) ([]model.Template, error) {
	if m.GetAllFn == nil {
		panic("mockTemplateRepository.GetAll not implemented")
	}
	return m.GetAllFn(ctx)
}

func (m *mockTemplateRepository) GetByID(ctx context.Context, id uint) (*model.Template, error) {
	if m.GetByIDFn == nil {
		panic("mockTemplateRepository.GetByID not implemented")
	}
	return m.GetByIDFn(ctx, id)
}

func (m *mockTemplateRepository) GetByName(ctx context.Context, name string) (*model.Template, error) {
	if m.GetByNameFn == nil {
		panic("mockTemplateRepository.GetByName not implemented")
	}
	return m.GetByNameFn(ctx, name)
}

func (m *mockTemplateRepository) Delete(ctx context.Context, template *model.Template) error {
	if m.DeleteFn == nil {
		panic("mockTemplateRepository.Delete not implemented")
	}
	return m.DeleteFn(ctx, template)
}

// mockForge implements forge.Forge.
type mockForge struct {
	GetAuthenticatedUserFn func(ctx context.Context) (*model.ForgeUser, error)
	GetUserByTokenFn       func(ctx context.Context, token string) (*model.ForgeUser, error)
	IsMemberOfOwnerFn      func(ctx context.Context, owner, username string) (bool, error)
	CreateRepoFn           func(ctx context.Context, owner, name string) (*model.Repo, error)
	GetRepoFn              func(ctx context.Context, id int64) (*model.Repo, error)
	GetOrgReposFn          func(ctx context.Context, name string) ([]model.Repo, error)
	DeleteRepoFn           func(ctx context.Context, owner, name string) error
	CreateWebhookFn        func(ctx context.Context, owner, repo, callbackURL, secret, branch string) error
	DeleteWebhookFn        func(ctx context.Context, owner, repo, callbackURL string) error
}

func (m *mockForge) GetAuthenticatedUser(ctx context.Context) (*model.ForgeUser, error) {
	if m.GetAuthenticatedUserFn == nil {
		panic("mockForge.GetAuthenticatedUser not implemented")
	}
	return m.GetAuthenticatedUserFn(ctx)
}

func (m *mockForge) GetUserByToken(ctx context.Context, token string) (*model.ForgeUser, error) {
	if m.GetUserByTokenFn == nil {
		panic("mockForge.GetUserByToken not implemented")
	}
	return m.GetUserByTokenFn(ctx, token)
}

func (m *mockForge) IsMemberOfOwner(ctx context.Context, owner, username string) (bool, error) {
	if m.IsMemberOfOwnerFn == nil {
		panic("mockForge.IsMemberOfOwner not implemented")
	}
	return m.IsMemberOfOwnerFn(ctx, owner, username)
}

func (m *mockForge) CreateRepo(ctx context.Context, owner, name string) (*model.Repo, error) {
	if m.CreateRepoFn == nil {
		panic("mockForge.CreateRepo not implemented")
	}
	return m.CreateRepoFn(ctx, owner, name)
}

func (m *mockForge) GetRepo(ctx context.Context, id int64) (*model.Repo, error) {
	if m.GetRepoFn == nil {
		panic("mockForge.GetRepo not implemented")
	}
	return m.GetRepoFn(ctx, id)
}

func (m *mockForge) GetOrgRepos(ctx context.Context, name string) ([]model.Repo, error) {
	if m.GetOrgReposFn == nil {
		panic("mockForge.GetOrgRepos not implemented")
	}
	return m.GetOrgReposFn(ctx, name)
}

func (m *mockForge) DeleteRepo(ctx context.Context, owner, name string) error {
	if m.DeleteRepoFn == nil {
		panic("mockForge.DeleteRepo not implemented")
	}
	return m.DeleteRepoFn(ctx, owner, name)
}

func (m *mockForge) CreateWebhook(ctx context.Context, owner, repo, callbackURL, secret, branch string) error {
	if m.CreateWebhookFn == nil {
		panic("mockForge.CreateWebhook not implemented")
	}
	return m.CreateWebhookFn(ctx, owner, repo, callbackURL, secret, branch)
}

func (m *mockForge) DeleteWebhook(ctx context.Context, owner, repo, callbackURL string) error {
	if m.DeleteWebhookFn == nil {
		panic("mockForge.DeleteWebhook not implemented")
	}
	return m.DeleteWebhookFn(ctx, owner, repo, callbackURL)
}

// mockCI implements ci.CI.
type mockCI struct {
	ActivateRepoFn func(ctx context.Context, forgeRemoteID int64, slug string) (*model.CIRepo, error)
	DeleteRepoFn   func(ctx context.Context, ciRepoID int64, slug string) error
	GetBuildsFn    func(ctx context.Context, ciRepoID int64, slug string) ([]model.Build, error)
}

func (m *mockCI) ActivateRepo(ctx context.Context, forgeRemoteID int64, slug string) (*model.CIRepo, error) {
	if m.ActivateRepoFn == nil {
		panic("mockCI.ActivateRepo not implemented")
	}
	return m.ActivateRepoFn(ctx, forgeRemoteID, slug)
}

func (m *mockCI) DeleteRepo(ctx context.Context, ciRepoID int64, slug string) error {
	if m.DeleteRepoFn == nil {
		panic("mockCI.DeleteRepo not implemented")
	}
	return m.DeleteRepoFn(ctx, ciRepoID, slug)
}

func (m *mockCI) GetBuilds(ctx context.Context, ciRepoID int64, slug string) ([]model.Build, error) {
	if m.GetBuildsFn == nil {
		panic("mockCI.GetBuilds not implemented")
	}
	return m.GetBuildsFn(ctx, ciRepoID, slug)
}

var errBoom = fmt.Errorf("boom")
