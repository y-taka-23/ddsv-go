package when

import (
	"fmt"

	"github.com/y-taka-23/ddsv-go/deadlock/rule/vars"
)

type Guard func(vars.Shared) (bool, error)

type Testee struct {
	name vars.Name
}

func Var(x vars.Name) Testee {
	return Testee{name: x}
}

func (t Testee) Is(n int) Guard {
	return func(vs vars.Shared) (bool, error) {
		val, ok := vs[t.name]
		if !ok {
			return false, fmt.Errorf("undeclared variable: %s", t.name)
		}
		return val == n, nil
	}

}

func (t Testee) IsNot(n int) Guard {
	return func(vs vars.Shared) (bool, error) {
		val, ok := vs[t.name]
		if !ok {
			return false, fmt.Errorf("undeclared variable: %s", t.name)
		}
		return val != n, nil
	}
}
