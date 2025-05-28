package gitlab

type InvalidResponseError struct {
	Msg        string
	StatusCode int
}

func (e *InvalidResponseError) Error() string {
	return e.Msg
}
