package response

type Build struct {
	ID        int64     `json:"id"`
	Number    int64     `json:"number"`
	Status    string    `json:"status"`
	Branch    string    `json:"branch"`
	Commit    string    `json:"commit"`
	Duration  int64     `json:"duration"`
	Link      string    `json:"link"`
}
