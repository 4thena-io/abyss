package dto

type ActivateRepoRequestDTO struct {
	Owner					string
	Name					string
	CloneUrl			string
	ForgeRemoteId	int64
}

type ActivateRepoResponseDTO struct {
	RepoId	int64
	RepoUrl	string
}
