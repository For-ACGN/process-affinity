package affinity

import (
	"fmt"

	"github.com/pkg/errors"
)

func newError(name string, err error, v ...interface{}) error {
	if err != nil {
		return errors.Errorf("%s: %s, because %s", name, fmt.Sprint(v...), err)
	}
	return errors.Errorf("%s: %s", name, fmt.Sprint(v...))
}

func newErrorf(name string, err error, format string, v ...interface{}) error {
	if err != nil {
		return errors.Errorf("%s: %s, because %s", name, fmt.Sprintf(format, v...), err)
	}
	return errors.Errorf("%s: %s", name, fmt.Sprintf(format, v...))
}
