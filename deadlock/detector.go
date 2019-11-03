// Package deadlock provides a simple algorithm checker.
// It contains a DSL to describe multi-process system
// and a detector for deadlocks in the concurrent product of the processes.
package deadlock

import (
	"github.com/y-taka-23/ddsv-go/deadlock/rule"
)

// Detector searches the state space of given system
// and reports presence of deadlocks.
type Detector interface {
	Detect(s System) (Report, error)
}

func NewDetector() Detector {
	return detector{}
}

type detector struct{}

func (d detector) Detect(s System) (Report, error) {

	visited := StateSet{}
	transited := TransitionSet{}
	accepting := StateSet{}
	deadlocked := StateSet{}
	traces := TransitionSet{}

	initial := d.initialize(s)
	queue := []State{initial}

	for len(queue) > 0 {
		from := queue[0]
		queue = queue[1:]

		if _, ok := visited[from.Id()]; ok {
			continue
		}
		visited[from.Id()] = from

		nexts := 0
		for _, p := range s.Processes() {
			// The locations of every processes are
			// certainly defined inductively
			focus, _ := from.Locations()[p.Id()]
			for _, r := range p.Rules()[focus] {

				fireable, err := r.Guard()(from.SharedVars())
				if err != nil {
					return report{}, err
				}
				if !fireable {
					continue
				}

				nextLocs := map[ProcessId]rule.Location{}
				for pid, l := range from.Locations() {
					nextLocs[pid] = l
				}
				nextLocs[p.Id()] = r.Target()

				nextVars, err := r.Action()(from.SharedVars())
				if err != nil {
					return report{}, err
				}

				to := state{
					locations:  nextLocs,
					sharedVars: nextVars,
					upstream:   "",
				}

				t := transition{
					process: p.Id(),
					label:   r.Label(),
					source:  from.Id(),
					target:  to.Id(),
				}
				transited[t.Id()] = t
				nexts++

				// assume that state.Id() is independent from state.upstream
				to.upstream = t.Id()
				queue = append(queue, to)
			}
		}

		if nexts == 0 {
			if acceptable(s, from) {
				accepting[from.Id()] = from
				continue
			}
			deadlocked[from.Id()] = from
			up := from.Upstream()
			for up != "" {
				// states and transitions in the path are certainly registered
				t, _ := transited[up]
				traces[up] = t
				prev, _ := visited[t.Source()]
				up = prev.Upstream()
			}
		}

	}

	return report{
		visited:    visited,
		transited:  transited,
		initial:    initial.Id(),
		accepting:  accepting,
		deadlocked: deadlocked,
		traces:     traces,
	}, nil

}

func (_ detector) initialize(s System) State {
	ls := LocationSet{}
	for _, p := range s.Processes() {
		ls[p.Id()] = p.EntryPoint()
	}
	return state{
		locations:  ls,
		sharedVars: s.InitVars(),
		upstream:   "",
	}
}

func acceptable(s System, state State) bool {
	for _, p := range s.Processes() {
		// The locations is certainly defined
		focus, _ := state.Locations()[p.Id()]
		found := false
		for _, l := range p.HaltingPoints() {
			if focus == l {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
