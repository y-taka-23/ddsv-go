package deadlock

import (
	"github.com/y-taka-23/ddsv-go/deadlock/rule"
)

type ProcessId string

type Process interface {
	Id() ProcessId
	EntryPoint() rule.Location
	Rules() rule.RuleSet
	EnterAt(rule.Location) Process
	Define(rule.Rule) Process
}

func NewProcess() Process {
	return process{
		id:         "",
		entryPoint: "",
		rules:      rule.RuleSet{},
	}
}

type process struct {
	id         ProcessId
	entryPoint rule.Location
	rules      rule.RuleSet
}

func (p process) Id() ProcessId {
	return p.id
}

func (p process) EntryPoint() rule.Location {
	return p.entryPoint
}

func (p process) EnterAt(l rule.Location) Process {
	p.entryPoint = l
	return p
}

func (p process) Define(r rule.Rule) Process {
	rs, ok := p.rules[r.Source()]
	if !ok {
		p.rules[r.Source()] = []rule.Rule{r}
		return p
	}
	p.rules[r.Source()] = append(rs, r)
	return p
}

func (p process) Rules() rule.RuleSet {
	return p.rules
}

type System interface {
	InitVars() rule.SharedVars
	Processes() []Process
	Declare(rule.SharedVars) System
	Register(ProcessId, Process) System
}

func NewSystem() System {
	return system{
		initVars:  rule.SharedVars{},
		processes: []Process{},
	}
}

type system struct {
	initVars  rule.SharedVars
	processes []Process
}

func (s system) InitVars() rule.SharedVars {
	return s.initVars
}

func (s system) Processes() []Process {
	return s.processes
}

func (s system) Declare(decls rule.SharedVars) System {
	sv := rule.SharedVars{}
	for x, n := range decls {
		sv[x] = n
	}
	s.initVars = sv
	return s
}

func (s system) Register(pid ProcessId, p Process) System {
	registered := process{
		id:         pid,
		entryPoint: p.EntryPoint(),
		rules:      p.Rules(),
	}
	s.processes = append(s.processes, registered)
	return s
}
