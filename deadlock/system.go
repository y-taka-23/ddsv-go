package deadlock

import (
	"github.com/y-taka-23/ddsv-go/deadlock/rule"
	"github.com/y-taka-23/ddsv-go/deadlock/rule/vars"
)

type ProcessId string

// Process represents a single process in a concurrent system.
type Process interface {
	Id() ProcessId
	EntryPoint() rule.Location
	Rules() rule.RuleSet
	HaltingPoints() []rule.Location
	EnterAt(rule.Location) Process
	Define(rule.Rule) Process
	HaltAt(...rule.Location) Process
}

func NewProcess() Process {
	return process{
		id:            "",
		entryPoint:    "",
		rules:         rule.RuleSet{},
		haltingPoints: []rule.Location{},
	}
}

type process struct {
	id            ProcessId
	entryPoint    rule.Location
	rules         rule.RuleSet
	haltingPoints []rule.Location
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

func (p process) HaltAt(ls ...rule.Location) Process {
	p.haltingPoints = ls
	return p
}

func (p process) HaltingPoints() []rule.Location {
	return p.haltingPoints
}

// System represents a set of processes.
// In the deadlock detection, they act concurrently
// accessing the pre-declared global shared variables.
type System interface {
	InitVars() vars.Shared
	Processes() []Process
	Declare(vars.Shared) System
	Register(ProcessId, Process) System
}

func NewSystem() System {
	return system{
		initVars:  vars.Shared{},
		processes: []Process{},
	}
}

type system struct {
	initVars  vars.Shared
	processes []Process
}

func (s system) InitVars() vars.Shared {
	return s.initVars
}

func (s system) Processes() []Process {
	return s.processes
}

func (s system) Declare(decls vars.Shared) System {
	vs := vars.Shared{}
	for x, n := range decls {
		vs[x] = n
	}
	s.initVars = vs
	return s
}

func (s system) Register(pid ProcessId, p Process) System {
	registered := process{
		id:            pid,
		entryPoint:    p.EntryPoint(),
		rules:         p.Rules(),
		haltingPoints: p.HaltingPoints(),
	}
	s.processes = append(s.processes, registered)
	return s
}
