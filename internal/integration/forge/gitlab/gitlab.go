package gitlab

import (
	"context"
	"fmt"
	"net/http"

	gl "gitlab.com/gitlab-org/api/client-go"
	"github.com/4thena-io/abyss/internal/model"
)

type GitlabForge struct {
	client *gl.Client
	host   string
}

func NewGitlabForge(host, token string) (*GitlabForge, error) {
	opts := []gl.ClientOptionFunc{gl.WithBaseURL(host)}
	client, err := gl.NewClient(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}
	return &GitlabForge{client: client, host: host}, nil
}

func (f *GitlabForge) GetAuthenticatedUser(ctx context.Context) (*model.ForgeUser, error) {
	user, _, err := f.client.Users.CurrentUser(gl.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}
	name := user.Name
	if name == "" {
		name = user.Username
	}
	return &model.ForgeUser{
		ID:        int64(user.ID),
		Username:  user.Username,
		FullName:  name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}

func (f *GitlabForge) GetUserByToken(ctx context.Context, token string) (*model.ForgeUser, error) {
	opts := []gl.ClientOptionFunc{gl.WithBaseURL(f.host)}
	client, err := gl.NewClient(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create user client: %w", err)
	}
	user, _, err := client.Users.CurrentUser(gl.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	name := user.Name
	if name == "" {
		name = user.Username
	}
	return &model.ForgeUser{
		ID:        int64(user.ID),
		Username:  user.Username,
		FullName:  name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}

func (f *GitlabForge) IsMemberOfOwner(ctx context.Context, owner, username string) (bool, error) {
	members, resp, err := f.client.Groups.ListGroupMembers(owner, nil, gl.WithContext(ctx))
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to list group members: %w", err)
	}
	for _, m := range members {
		if m.Username == username {
			return true, nil
		}
	}
	return false, nil
}

func (f *GitlabForge) CreateRepo(ctx context.Context, owner, name string) (*model.Repo, error) {
	// Determine namespace by group name
	groups, _, err := f.client.Groups.SearchGroup(owner, gl.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to find group %s: %w", owner, err)
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("group %s not found", owner)
	}
	nsID := groups[0].ID

	proj, _, err := f.client.Projects.CreateProject(&gl.CreateProjectOptions{
		Name:        gl.Ptr(name),
		NamespaceID: gl.Ptr(nsID),
		Visibility:  gl.Ptr(gl.InternalVisibility),
	}, gl.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return &model.Repo{
		ID:       int64(proj.ID),
		Name:     proj.Name,
		FullName: proj.PathWithNamespace,
		URL:      proj.WebURL,
		CloneURL: proj.HTTPURLToRepo,
	}, nil
}

func (f *GitlabForge) GetRepo(ctx context.Context, id int64) (*model.Repo, error) {
	proj, _, err := f.client.Projects.GetProject(int(id), nil, gl.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	return &model.Repo{
		ID:       int64(proj.ID),
		Name:     proj.Name,
		FullName: proj.PathWithNamespace,
		URL:      proj.WebURL,
		CloneURL: proj.HTTPURLToRepo,
	}, nil
}

func (f *GitlabForge) GetOrgRepos(ctx context.Context, name string) ([]model.Repo, error) {
	var allProjects []*gl.Project
	opts := &gl.ListGroupProjectsOptions{ListOptions: gl.ListOptions{PerPage: 100}}
	for {
		projects, resp, err := f.client.Groups.ListGroupProjects(name, opts, gl.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to list group projects: %w", err)
		}
		allProjects = append(allProjects, projects...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	res := make([]model.Repo, len(allProjects))
	for i, proj := range allProjects {
		res[i] = model.Repo{
			ID:       int64(proj.ID),
			Name:     proj.Name,
			FullName: proj.PathWithNamespace,
			URL:      proj.WebURL,
			CloneURL: proj.HTTPURLToRepo,
		}
	}
	return res, nil
}

func (f *GitlabForge) DeleteRepo(ctx context.Context, owner, name string) error {
	pid := owner + "/" + name
	_, err := f.client.Projects.DeleteProject(pid, nil, gl.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	return nil
}

func (f *GitlabForge) CreateWebhook(ctx context.Context, owner, repo, callbackURL, secret string) error {
	pid := owner + "/" + repo
	pushEvents := true
	_, _, err := f.client.Projects.AddProjectHook(pid, &gl.AddProjectHookOptions{
		URL:        gl.Ptr(callbackURL),
		PushEvents: &pushEvents,
		Token:      gl.Ptr(secret),
	}, gl.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}
	return nil
}
