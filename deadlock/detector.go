package deadlock

import (
	"fmt"

	"github.com/y-taka-23/ddsv-go/deadlock/rule"
)

type Detector interface {
	Detect(s System) (Report, error)
}

func NewDetector() Detector {
	return detector{}
}

type detector struct{}

func (d detector) Detect(s System) (Report, error) {

	visited := StateSet{}
	transited := []Transition{}

	queue := []State{d.initialize(s)}

	for len(queue) > 0 {
		from := queue[0]
		queue = queue[1:]

		if _, ok := visited[from.Id()]; ok {
			continue
		}
		visited[from.Id()] = from

		ts := []Transition{}
		for _, p := range s.Processes() {
			focus, ok := from.Locations()[p.Id()]
			if !ok {
				return nil,
					fmt.Errorf("location of prosess %s is undefined", p.Id())
			}
			for _, r := range p.Rules()[focus] {
				nextLocs := map[ProcessId]rule.Location{}
				for pid, l := range from.Locations() {
					nextLocs[pid] = l
				}
				nextLocs[p.Id()] = r.Target()

				nextVars := r.Action()(from.SharedVars())
				to := state{locations: nextLocs, sharedVars: nextVars}

				t := transition{
					process: p.Id(),
					label:   r.Label(),
					source:  from.Id(),
					target:  to.Id(),
				}
				ts = append(ts, t)
				queue = append(queue, to)
			}
		}

		transited = append(transited, ts...)
	}

	return report{
		visited:   visited,
		transited: transited,
	}, nil

}

func (_ detector) initialize(s System) State {
	ls := LocationSet{}
	for _, p := range s.Processes() {
		ls[p.Id()] = p.EntryPoint()
	}
	return state{
		locations:  ls,
		sharedVars: rule.SharedVars{},
	}
}
