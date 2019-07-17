package syncd

import (
	"fmt"
	"strings"
)

type multiError []error

func (me *multiError) Add(err error) {
	*me = append(*me, err)
}

func (me multiError) Error() string {
	switch n := len(me); n {
	case 0:
		return ""
	case 1:
		return me[0].Error()
	default:
		msgs := make([]string, 0, n+1)

		msgs = append(msgs, fmt.Sprintf("multiple errors (%d):", n))
		for i, err := range me {
			msgs = append(msgs, fmt.Sprintf("%d: %s", i+1, err.Error()))
		}

		return strings.Join(msgs, "\n\t")
	}
}
