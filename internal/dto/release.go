package dto

type CreateReleaseRequestDto struct {
	Tag string `json:"tag"`
}

type CreateReleaseRecordDto struct {
	Tag string
	App string
}
