package git

import gogit "github.com/go-git/go-git/v5"

type GitClient struct{}

func New() *GitClient {
	return &GitClient{}
}

func (c *GitClient) CloneTemplate(repoURL, token, localPath string) error {
	_, err := gogit.PlainClone(localPath, false, &gogit.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *GitClient) PushRepo(token, localPath string) error {
	r, err := gogit.PlainOpen(localPath)
	if err != nil {
		return err
	}

	return r.Push(&gogit.PushOptions{})
}
