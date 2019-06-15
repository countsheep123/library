package obj

import "fmt"

type InvalidRequest struct {
	Msg string
}

func (e InvalidRequest) Error() string {
	return fmt.Sprintf("invalid request: %s", e.Msg)
}

type NotFound struct{}

func (e NotFound) Error() string {
	return "not found"
}

type Internal struct {
	Msg string
}

func (e Internal) Error() string {
	return fmt.Sprintf("internal error: %s", e.Msg)
}

type Duplicate struct{}

func (e Duplicate) Error() string {
	return "duplicated"
}
