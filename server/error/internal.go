package error

import "fmt"

type Internal struct {
	InnerErr error
	Msg      string
}

func (e Internal) Error() string {
	return fmt.Sprintf("%s: %s", e.Msg, e.InnerErr)
}

func (e Internal) Unwrap() error {
	return e.InnerErr
}
