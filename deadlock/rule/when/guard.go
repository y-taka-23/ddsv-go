package when

import (
	"fmt"

	"github.com/y-taka-23/ddsv-go/deadlock/rule/vars"
)

type Guard func(vars.Shared) (bool, error)

func Eq(x vars.Name, n int) Guard {
	return func(vs vars.Shared) (bool, error) {
		val, ok := vs[x]
		if !ok {
			return false, fmt.Errorf("undeclared variable: %s", x)
		}
		return val == n, nil
	}
}

func NotEq(x vars.Name, n int) Guard {
	return func(vs vars.Shared) (bool, error) {
		val, ok := vs[x]
		if !ok {
			return false, fmt.Errorf("undeclared variable: %s", x)
		}
		return val != n, nil
	}
}
