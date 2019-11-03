/*
 * The when package provides human-readable DSL for conditional transitions.
 */
package when

import (
	"fmt"

	"github.com/y-taka-23/ddsv-go/deadlock/rule/vars"
)

// Guard determines whether a transition is fireable
// under the given values of shared variables.
// If the specified variable name is undeclared, it returns an error.
type Guard func(vars.Shared) (bool, error)

type Testee struct {
	name vars.Name
}

func Var(x vars.Name) Testee {
	return Testee{name: x}
}

func (t Testee) Is(n int) Guard {
	return t.check(func(x, y int) bool { return x == y }, n)
}

func (t Testee) IsNot(n int) Guard {
	return t.check(func(x, y int) bool { return x != y }, n)
}

func (t Testee) IsLessThan(n int) Guard {
	return t.check(func(x, y int) bool { return x < y }, n)
}

func (t Testee) IsGreaterThan(n int) Guard {
	return t.check(func(x, y int) bool { return x > y }, n)
}

func (t Testee) check(op func(x, y int) bool, n int) Guard {
	return func(vs vars.Shared) (bool, error) {
		val, ok := vs[t.name]
		if !ok {
			return false, fmt.Errorf("undeclared variable: %s", t.name)
		}
		return op(val, n), nil
	}
}
