package errutils

import (
	"fmt"
	"github.com/pkg/errors"
)

func PanicToError(p interface{}) error {
	switch v := p.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		return fmt.Errorf("%v", p)
	}
}

