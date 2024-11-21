package cmd

import (
	"fmt"
)

type FieldValidator struct {
	valid bool
	err   error
}

func (fv *FieldValidator) InArray(value string, array []string) {
	fv.valid = false
	for _, v := range array {
		if value == v {
			fv.valid = true
			return
		}
	}
	fv.err = fmt.Errorf("value '%s' is not in the allowed list", value)
}

func (fv *FieldValidator) IsValid() bool {
	return fv.valid
}

func (fv *FieldValidator) Error() error {
	return fv.err
}
