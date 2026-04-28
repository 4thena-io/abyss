package response

type Deployment struct {
	ID          uint   `json:"id"`
	AppID       uint   `json:"appId"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
	Commit      string `json:"commit"`
	TriggeredBy string `json:"triggeredBy"`
	Duration    int64  `json:"duration"`
	DeployedAt  string `json:"deployedAt"`
}
