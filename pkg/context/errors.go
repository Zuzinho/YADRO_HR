package context

import "fmt"

type NoEnoughArgsError struct {
	wait int
	has  int
}

func newNoEnoughArgsError(wait, has int) NoEnoughArgsError {
	return NoEnoughArgsError{
		wait: wait,
		has:  has,
	}
}

func (err NoEnoughArgsError) Error() string {
	return fmt.Sprintf("wait %d args - has %d args", err.wait, err.has)
}

var (
	NoEnoughArgsErr = NoEnoughArgsError{}
)
