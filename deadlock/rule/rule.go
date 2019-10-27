package rule

type Location string

type Label string

type RuleSet map[Location][]Rule

type VarName string

type SharedVars map[VarName]int

type Action func(SharedVars) SharedVars

type Rule interface {
	Source() Location
	Target() Location
	Label() Label
	Action() Action
	Do(Label) Rule
	MoveTo(Location) Rule
}

func At(l Location) Rule {
	return rule{
		source: l,
		target: l,
		label:  "",
		action: func(sv SharedVars) SharedVars { return sv },
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

func (r rule) Do(lbl Label) Rule {
	r.label = lbl
	return r
}

func (r rule) MoveTo(l Location) Rule {
	r.target = l
	return r
}
