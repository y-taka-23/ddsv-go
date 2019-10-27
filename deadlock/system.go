package deadlock

type Process interface {
	Id() processId
	EntryPoint() location
	Rules() ruleSet
	EnterAt(location) Process
	Define(Rule) Process
}

func NewProcess() Process {
	return process{
		id:    "",
		rules: ruleSet{},
	}
}

type process struct {
	id         processId
	entryPoint location
	rules      ruleSet
}

func (p process) Id() processId {
	return p.id
}

func (p process) EntryPoint() location {
	return p.entryPoint
}

func (p process) EnterAt(l location) Process {
	p.entryPoint = l
	return p
}

func (p process) Define(r Rule) Process {
	rs, ok := p.rules[r.Source()]
	if !ok {
		p.rules[r.Source()] = []Rule{r}
		return p
	}
	p.rules[r.Source()] = append(rs, r)
	return p
}

func (p process) Rules() ruleSet {
	return p.rules
}

type Rule interface {
	Source() location
	Target() location
	Label() label
	Action() action
	Do(label) Rule
	MoveTo(location) Rule
}

func At(l location) Rule {
	return rule{
		source: l,
		target: l,
		label:  "",
		action: func(sv sharedVars) sharedVars { return sv },
	}
}

type rule struct {
	source location
	target location
	label  label
	action action
}

func (r rule) Source() location {
	return r.source
}

func (r rule) Target() location {
	return r.target
}

func (r rule) Label() label {
	return r.label
}

func (r rule) Action() action {
	return r.action
}

func (r rule) Do(lbl label) Rule {
	r.label = lbl
	return r
}

func (r rule) MoveTo(l location) Rule {
	r.target = l
	return r
}

type System interface {
	Processes() []Process
	Register(processId, Process) System
}

func NewSystem() System {
	return system{
		processes: []Process{},
	}
}

type system struct {
	processes []Process
	initVars  sharedVars
}

func (s system) Processes() []Process {
	return s.processes
}

func (s system) Register(pid processId, p Process) System {
	registered := process{
		id:         pid,
		rules:      p.Rules(),
		entryPoint: p.EntryPoint(),
	}
	s.processes = append(s.processes, registered)
	return s
}

type processId string
type location string
type ruleSet map[location][]Rule
type locationMap map[processId]location

type label string
type action func(sharedVars) sharedVars

type varName string

type sharedVars map[varName]int
