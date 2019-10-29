package rule

import (
	"fmt"
)

type Location string

type Label string

type RuleSet map[Location][]Rule

type Rule interface {
	Source() Location
	Target() Location
	Label() Label
	Action() Action
	Do(Label, Action) Rule
	MoveTo(Location) Rule
}

func At(l Location) Rule {
	return rule{
		source: l,
		target: l,
		label:  "",
		action: func(sv SharedVars) (SharedVars, error) {
			return clone(sv), nil
		},
	}
}

type rule struct {
	source Location
	target Location
	label  Label
	action Action
}

func (r rule) Source() Location {
	return r.source
}

func (r rule) Target() Location {
	return r.target
}

func (r rule) Label() Label {
	return r.label
}

func (r rule) Action() Action {
	return r.action
}

func (r rule) Do(lbl Label, a Action) Rule {
	r.label = lbl
	r.action = a
	return r
}

func (r rule) MoveTo(l Location) Rule {
	r.target = l
	return r
}

type VarName string

type SharedVars map[VarName]int

type Action func(SharedVars) (SharedVars, error)

func Copy(y, x VarName) Action {
	return func(sv SharedVars) (SharedVars, error) {
		modified := clone(sv)
		if _, ok := sv[y]; !ok {
			return SharedVars{}, fmt.Errorf("undeclared variable: %s", y)
		}
		if _, ok := modified[x]; !ok {
			return SharedVars{}, fmt.Errorf("undeclared variable: %s", x)
		}
		modified[x] = sv[y]
		return modified, nil
	}
}

func Set(n int, x VarName) Action {
	return func(sv SharedVars) (SharedVars, error) {
		modified := clone(sv)
		if _, ok := modified[x]; !ok {
			return SharedVars{}, fmt.Errorf("undeclared variable: %s", x)
		}
		modified[x] = n
		return modified, nil
	}
}

func Add(n int, x VarName) Action {
	return func(sv SharedVars) (SharedVars, error) {
		modified := clone(sv)
		if _, ok := modified[x]; !ok {
			return SharedVars{}, fmt.Errorf("undeclared variable: %s", x)
		}
		modified[x] = sv[x] + n
		return modified, nil
	}
}

func clone(sv SharedVars) SharedVars {
	c := map[VarName]int{}
	for k, v := range sv {
		c[k] = v
	}
	return c
}
