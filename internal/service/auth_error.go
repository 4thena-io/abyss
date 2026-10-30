package service

// AuthStageError carries a safe diagnostic stage. Its cause may contain secrets
// returned by an OAuth provider and must never be written to logs or responses.
type AuthStageError struct {
	Stage string
	Err   error
}

func (e *AuthStageError) Error() string { return "authentication failed during " + e.Stage }
func (e *AuthStageError) Unwrap() error { return e.Err }
