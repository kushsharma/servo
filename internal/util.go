package internal

import (
	"errors"
	"fmt"
)

// ErrMerge merges array of errors and return a single error
func ErrMerge(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	str := ""
	for idx, err := range errs {
		str += fmt.Sprintf("%d. %s", idx+1, err.Error())
	}
	return errors.New(str)
}
