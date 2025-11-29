package constant

type Status string

const (
	StatusNew      Status = "NEW"
	StatusActive   Status = "ACTIVE"
	StatusDisabled Status = "DISABLED"
	StatusDeleted  Status = "DELETED"
)
