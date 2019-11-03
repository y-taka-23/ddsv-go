// Package do provides human-readable DSL for variable mutations.
package do

import (
	"fmt"

	"github.com/y-taka-23/ddsv-go/deadlock/rule/vars"
)

// Action changes the values of shared variables.
// If the specified variable name is undeclared, it returns an error.
type Action func(vars.Shared) (vars.Shared, error)

func Nothing() Action {
	return func(vs vars.Shared) (vars.Shared, error) {
		return vs.Clone(), nil
	}
}

type Operation interface {
	ToVar(vars.Name) Action
}

type copyVar struct {
	name vars.Name
}

func CopyVar(y vars.Name) Operation {
	return copyVar{name: y}
}

func (o copyVar) ToVar(x vars.Name) Action {
	return func(vs vars.Shared) (vars.Shared, error) {
		modified := vs.Clone()
		if _, ok := vs[o.name]; !ok {
			return vars.Shared{}, fmt.Errorf("undeclared variable: %s", o.name)
		}
		if _, ok := modified[x]; !ok {
			return vars.Shared{}, fmt.Errorf("undeclared variable: %s", x)
		}
		modified[x] = vs[o.name]
		return modified, nil
	}
}

type set struct {
	val int
}

func Set(n int) Operation {
	return set{val: n}
}

func (o set) ToVar(x vars.Name) Action {
	return func(vs vars.Shared) (vars.Shared, error) {
		modified := vs.Clone()
		if _, ok := modified[x]; !ok {
			return vars.Shared{}, fmt.Errorf("undeclared variable: %s", x)
		}
		modified[x] = o.val
		return modified, nil
	}
}

type add struct {
	val int
}

func Add(n int) Operation {
	return add{val: n}
}

func (o add) ToVar(x vars.Name) Action {
	return func(vs vars.Shared) (vars.Shared, error) {
		modified := vs.Clone()
		if _, ok := modified[x]; !ok {
			return vars.Shared{}, fmt.Errorf("undeclared variable: %s", x)
		}
		modified[x] = vs[x] + o.val
		return modified, nil
	}
}
