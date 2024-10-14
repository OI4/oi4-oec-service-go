package container

import "fmt"

type Error struct {
	Message string
	Err     error
}

func (w *Error) Error() string {
	return fmt.Sprintf("%s: %v", w.Message, w.Err)
}
