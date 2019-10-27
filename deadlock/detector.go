package deadlock

import "fmt"

type Detector interface {
	Detect(s System) (Report, error)
}

func NewDetector() Detector {
	return detector{}
}

type detector struct{}

func (d detector) Detect(s System) (Report, error) {

	visited := stateSet{}
	transited := []transition{}

	queue := []state{d.initialize(s)}

	for len(queue) > 0 {
		from := queue[0]
		queue = queue[1:]

		if _, ok := visited[from.id()]; ok {
			continue
		}
		visited[from.id()] = from

		ts := []transition{}
		for _, p := range s.Processes() {
			focus, ok := from.locations[p.Id()]
			if !ok {
				return nil,
					fmt.Errorf("location of prosess %s is undefined", p.Id())
			}
			for _, r := range p.Rules()[focus] {
				nextLocs := map[processId]location{}
				for pid, l := range from.locations {
					nextLocs[pid] = l
				}
				nextLocs[p.Id()] = r.Target()

				nextVars := r.Action()(from.sharedVars)
				to := state{locations: nextLocs, sharedVars: nextVars}

				t := transition{
					process: p.Id(),
					label:   r.Label(),
					source:  from.id(),
					target:  to.id(),
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

func (_ detector) initialize(s System) state {
	ls := map[processId]location{}
	for _, p := range s.Processes() {
		ls[p.Id()] = p.EntryPoint()
	}
	return state{
		locations:  ls,
		sharedVars: sharedVars{},
	}
}
