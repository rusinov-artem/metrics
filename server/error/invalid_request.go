package error

import "fmt"

type InvalidRequest struct {
	InnerErr error
	Msg      string
}

func (e InvalidRequest) Error() string {
	return fmt.Sprintf("%s: %s", e.Msg, e.InnerErr)
}

func (e InvalidRequest) Unwrap() error {
	return e.InnerErr
}
