/*
 * The rule package provides human-readable DSL to define transition rules.
 */
package rule

import (
	"github.com/y-taka-23/ddsv-go/deadlock/rule/do"
	"github.com/y-taka-23/ddsv-go/deadlock/rule/vars"
	"github.com/y-taka-23/ddsv-go/deadlock/rule/when"
)

// Location represents the program counter of each process.
type Location string

// Label is human-readable name of each transition.
type Label string

type RuleSet map[Location][]Rule

// Rule defines transition rules of the process.
type Rule interface {
	Source() Location
	Target() Location
	Guard() when.Guard
	Label() Label
	Action() do.Action
	Only(when.Guard) Rule
	Let(Label, do.Action) Rule
	MoveTo(Location) Rule
}

func At(l Location) Rule {
	return rule{
		source: l,
		target: l,
		label:  "",
		guard: func(_ vars.Shared) (bool, error) {
			return true, nil
		},
		action: func(vs vars.Shared) (vars.Shared, error) {
			return vs.Clone(), nil
		},
	}
}

type rule struct {
	source Location
	target Location
	guard  when.Guard
	label  Label
	action do.Action
}

func (r rule) Source() Location {
	return r.source
}

func (r rule) Target() Location {
	return r.target
}

func (r rule) Guard() when.Guard {
	return r.guard
}

func (r rule) Label() Label {
	return r.label
}

func (r rule) Action() do.Action {
	return r.action
}

func (r rule) Only(g when.Guard) Rule {
	r.guard = g
	return r
}

func (r rule) Let(lbl Label, a do.Action) Rule {
	r.label = lbl
	r.action = a
	return r
}

func (r rule) MoveTo(l Location) Rule {
	r.target = l
	return r
}
