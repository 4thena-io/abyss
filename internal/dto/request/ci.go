package request

type ActivateRepo struct {
	Owner					string
	Name					string
	CloneUrl			string
	ForgeRemoteId	int64
}
