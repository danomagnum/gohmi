package hmitags

import "fmt"

type ErrTagNotFound struct {
	key   string
	other string
}

func (e ErrTagNotFound) Error() string {
	return fmt.Sprintf("tag '%s' not found (%s)", e.key, e.other)
}
