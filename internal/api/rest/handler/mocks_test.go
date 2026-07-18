package handler

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
)

// Minimal test doubles for the repository/forge/ci interfaces the services
// depend on. Only the methods exercised by these handler tests are wired up;
// anything else panics so an unexpected call fails loudly.

type fakeAppRepository struct {
	SaveFn         func(ctx context.Context, app *model.App) error
	GetAllFn       func(ctx context.Context) ([]model.App, error)
	GetByIDFn      func(ctx context.Context, id uint) (*model.App, error)
	GetByNameFn    func(ctx context.Context, name string) (*model.App, error)
	GetByProjectFn func(ctx context.Context, id uint) ([]model.App, error)
	UpdateFn       func(ctx context.Context, app *model.App) error
	DeleteFn       func(ctx context.Context, app *model.App) error
}

func (f *fakeAppRepository) Save(ctx context.Context, app *model.App) error {
	if f.SaveFn == nil {
		panic("fakeAppRepository.Save not implemented")
	}
	return f.SaveFn(ctx, app)
}
func (f *fakeAppRepository) GetAll(ctx context.Context) ([]model.App, error) {
	if f.GetAllFn == nil {
		panic("fakeAppRepository.GetAll not implemented")
	}
	return f.GetAllFn(ctx)
}
func (f *fakeAppRepository) GetByID(ctx context.Context, id uint) (*model.App, error) {
	if f.GetByIDFn == nil {
		panic("fakeAppRepository.GetByID not implemented")
	}
	return f.GetByIDFn(ctx, id)
}
func (f *fakeAppRepository) GetByName(ctx context.Context, name string) (*model.App, error) {
	if f.GetByNameFn == nil {
		panic("fakeAppRepository.GetByName not implemented")
	}
	return f.GetByNameFn(ctx, name)
}
func (f *fakeAppRepository) GetByProject(ctx context.Context, id uint) ([]model.App, error) {
	if f.GetByProjectFn == nil {
		panic("fakeAppRepository.GetByProject not implemented")
	}
	return f.GetByProjectFn(ctx, id)
}
func (f *fakeAppRepository) GetByTemplate(ctx context.Context, id uint) ([]model.App, error) {
	panic("fakeAppRepository.GetByTemplate not implemented")
}
func (f *fakeAppRepository) Update(ctx context.Context, app *model.App) error {
	if f.UpdateFn == nil {
		panic("fakeAppRepository.Update not implemented")
	}
	return f.UpdateFn(ctx, app)
}
func (f *fakeAppRepository) Delete(ctx context.Context, app *model.App) error {
	if f.DeleteFn == nil {
		panic("fakeAppRepository.Delete not implemented")
	}
	return f.DeleteFn(ctx, app)
}
func (f *fakeAppRepository) CountByTeam(ctx context.Context, teamID uint) (int64, error) {
	panic("fakeAppRepository.CountByTeam not implemented")
}

type fakeProjectRepository struct {
	SaveFn      func(ctx context.Context, project *model.Project) error
	UpdateFn    func(ctx context.Context, project *model.Project) error
	GetAllFn    func(ctx context.Context) ([]model.Project, error)
	GetByIDFn   func(ctx context.Context, id uint) (*model.Project, error)
	GetByNameFn func(ctx context.Context, name string) (*model.Project, error)
	DeleteFn    func(ctx context.Context, project *model.Project) error
}

func (f *fakeProjectRepository) Save(ctx context.Context, project *model.Project) error {
	if f.SaveFn == nil {
		panic("fakeProjectRepository.Save not implemented")
	}
	return f.SaveFn(ctx, project)
}
func (f *fakeProjectRepository) Update(ctx context.Context, project *model.Project) error {
	if f.UpdateFn == nil {
		panic("fakeProjectRepository.Update not implemented")
	}
	return f.UpdateFn(ctx, project)
}
func (f *fakeProjectRepository) GetAll(ctx context.Context) ([]model.Project, error) {
	if f.GetAllFn == nil {
		panic("fakeProjectRepository.GetAll not implemented")
	}
	return f.GetAllFn(ctx)
}
func (f *fakeProjectRepository) GetByID(ctx context.Context, id uint) (*model.Project, error) {
	if f.GetByIDFn == nil {
		panic("fakeProjectRepository.GetByID not implemented")
	}
	return f.GetByIDFn(ctx, id)
}
func (f *fakeProjectRepository) GetByName(ctx context.Context, name string) (*model.Project, error) {
	if f.GetByNameFn == nil {
		panic("fakeProjectRepository.GetByName not implemented")
	}
	return f.GetByNameFn(ctx, name)
}
func (f *fakeProjectRepository) GetByTeam(ctx context.Context, teamID uint) ([]model.Project, error) {
	panic("fakeProjectRepository.GetByTeam not implemented")
}
func (f *fakeProjectRepository) Delete(ctx context.Context, project *model.Project) error {
	if f.DeleteFn == nil {
		panic("fakeProjectRepository.Delete not implemented")
	}
	return f.DeleteFn(ctx, project)
}
func (f *fakeProjectRepository) CountByTeam(ctx context.Context, teamID uint) (int64, error) {
	panic("fakeProjectRepository.CountByTeam not implemented")
}

type fakeTemplateRepository struct{}

func (f *fakeTemplateRepository) Save(ctx context.Context, template *model.Template) error {
	panic("fakeTemplateRepository.Save not implemented")
}
func (f *fakeTemplateRepository) Update(ctx context.Context, template *model.Template) error {
	panic("fakeTemplateRepository.Update not implemented")
}
func (f *fakeTemplateRepository) GetAll(ctx context.Context) ([]model.Template, error) {
	panic("fakeTemplateRepository.GetAll not implemented")
}
func (f *fakeTemplateRepository) GetByID(ctx context.Context, id uint) (*model.Template, error) {
	panic("fakeTemplateRepository.GetByID not implemented")
}
func (f *fakeTemplateRepository) GetByName(ctx context.Context, name string) (*model.Template, error) {
	panic("fakeTemplateRepository.GetByName not implemented")
}
func (f *fakeTemplateRepository) Delete(ctx context.Context, template *model.Template) error {
	panic("fakeTemplateRepository.Delete not implemented")
}

type fakeForge struct{}

func (f *fakeForge) GetAuthenticatedUser(ctx context.Context) (*model.ForgeUser, error) {
	panic("fakeForge.GetAuthenticatedUser not implemented")
}
func (f *fakeForge) GetUserByToken(ctx context.Context, token string) (*model.ForgeUser, error) {
	panic("fakeForge.GetUserByToken not implemented")
}
func (f *fakeForge) IsMemberOfOwner(ctx context.Context, owner, username string) (bool, error) {
	panic("fakeForge.IsMemberOfOwner not implemented")
}
func (f *fakeForge) CreateRepo(ctx context.Context, owner, name string) (*model.Repo, error) {
	panic("fakeForge.CreateRepo not implemented")
}
func (f *fakeForge) GetRepo(ctx context.Context, id int64) (*model.Repo, error) {
	panic("fakeForge.GetRepo not implemented")
}
func (f *fakeForge) GetOrgRepos(ctx context.Context, name string) ([]model.Repo, error) {
	panic("fakeForge.GetOrgRepos not implemented")
}
func (f *fakeForge) DeleteRepo(ctx context.Context, owner, name string) error {
	panic("fakeForge.DeleteRepo not implemented")
}
func (f *fakeForge) CreateWebhook(ctx context.Context, owner, repo, callbackURL, secret, branch string) error {
	panic("fakeForge.CreateWebhook not implemented")
}
func (f *fakeForge) DeleteWebhook(ctx context.Context, owner, repo, callbackURL string) error {
	panic("fakeForge.DeleteWebhook not implemented")
}

type fakeCI struct{}

func (f *fakeCI) ActivateRepo(ctx context.Context, forgeRemoteID int64, slug string) (*model.CIRepo, error) {
	panic("fakeCI.ActivateRepo not implemented")
}
func (f *fakeCI) DeleteRepo(ctx context.Context, ciRepoID int64, slug string) error {
	panic("fakeCI.DeleteRepo not implemented")
}
func (f *fakeCI) GetBuilds(ctx context.Context, ciRepoID int64, slug string) ([]model.Build, error) {
	panic("fakeCI.GetBuilds not implemented")
}

type fakeDeploymentRepository struct {
	GetByAppFn func(ctx context.Context, appID uint) ([]model.Deployment, error)
}

func (f *fakeDeploymentRepository) Save(ctx context.Context, deployment *model.Deployment) error {
	panic("fakeDeploymentRepository.Save not implemented")
}
func (f *fakeDeploymentRepository) GetByApp(ctx context.Context, appID uint) ([]model.Deployment, error) {
	if f.GetByAppFn == nil {
		panic("fakeDeploymentRepository.GetByApp not implemented")
	}
	return f.GetByAppFn(ctx, appID)
}
