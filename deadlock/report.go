package deadlock

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/y-taka-23/ddsv-go/deadlock/rule"
	"github.com/y-taka-23/ddsv-go/deadlock/rule/vars"
)

type StateId string

type StateSet map[StateId]State

type LocationSet map[ProcessId]rule.Location

type State interface {
	Id() StateId
	Locations() LocationSet
	SharedVars() vars.Shared
	Upstream() TransitionId
}

type state struct {
	locations  LocationSet
	sharedVars vars.Shared
	upstream   TransitionId
}

// assume that s.Id() is independent from s.upstream
func (s state) Id() StateId {
	h := sha1.New()
	serial := fmt.Sprintf("%+v%+v", s.locations, s.sharedVars)
	h.Write([]byte(serial))
	return StateId(fmt.Sprintf("%x", h.Sum(nil)))
}

func (s state) Locations() LocationSet {
	return s.locations
}

func (s state) SharedVars() vars.Shared {
	return s.sharedVars
}

func (s state) Upstream() TransitionId {
	return s.upstream
}

type TransitionId string

type TransitionSet map[TransitionId]Transition

type Transition interface {
	Id() TransitionId
	Process() ProcessId
	Label() rule.Label
	Source() StateId
	Target() StateId
}

type transition struct {
	process ProcessId
	label   rule.Label
	source  StateId
	target  StateId
}

func (t transition) Id() TransitionId {
	h := sha1.New()
	serial := fmt.Sprintf("%+v", t)
	h.Write([]byte(serial))
	return TransitionId(fmt.Sprintf("%x", h.Sum(nil)))
}

func (t transition) Process() ProcessId {
	return t.process
}

func (t transition) Label() rule.Label {
	return t.label
}

func (t transition) Source() StateId {
	return t.source
}

func (t transition) Target() StateId {
	return t.target
}

type Report interface {
	Visited() StateSet
	Transited() TransitionSet
	Initial() StateId
	Accepting() StateSet
	Deadlocked() StateSet
	Traces() TransitionSet
}

type report struct {
	visited    StateSet
	transited  TransitionSet
	initial    StateId
	accepting  StateSet
	deadlocked StateSet
	traces     TransitionSet
}

func (rp report) Visited() StateSet {
	return rp.visited
}

func (rp report) Transited() TransitionSet {
	return rp.transited
}

func (rp report) Initial() StateId {
	return rp.initial
}

func (rp report) Accepting() StateSet {
	return rp.accepting
}

func (rp report) Deadlocked() StateSet {
	return rp.deadlocked
}

func (rp report) Traces() TransitionSet {
	return rp.traces
}

type Printer struct {
	writer io.Writer
}

func NewPrinter(w io.Writer) Printer {
	return Printer{writer: w}
}

func (pr Printer) Print(rp Report) (int, error) {
	written, err := fmt.Fprintln(pr.writer, "digraph {")
	if err != nil {
		return written, err
	}
	for _, s := range rp.Visited() {
		n := 0
		if s.Id() == rp.Initial() {
			n, err = pr.printInitial(s)
		} else if _, ok := rp.Accepting()[s.Id()]; ok {
			n, err = pr.printAccepting(s)
		} else if _, ok := rp.Deadlocked()[s.Id()]; ok {
			n, err = pr.printDeadlocked(s)
		} else {
			n, err = pr.printState(s)
		}
		written += n
		if err != nil {
			return written, err
		}
	}
	for _, t := range rp.Transited() {
		n := 0
		if _, ok := rp.Traces()[t.Id()]; ok {
			n, err = pr.printTrace(t)
		} else {
			n, err = pr.printTransition(t)
		}
		written += n
		if err != nil {
			return written, err
		}

	}
	n, err := fmt.Fprintln(pr.writer, "}")
	written += n
	if err != nil {
		return written, err
	}
	return written, nil
}

func (pr Printer) printState(s State) (int, error) {
	return fmt.Fprintf(
		pr.writer,
		"  \"%s\" [label=\"%s\"]\n",
		s.Id(), stateLabel(s),
	)
}

func (pr Printer) printInitial(s State) (int, error) {
	return fmt.Fprintf(
		pr.writer,
		"  \"%s\" [label=\"%s\", fillcolor=\"#AAFFFF\", style=\"solid,filled\"];\n",
		s.Id(), stateLabel(s),
	)
}

func (pr Printer) printAccepting(s State) (int, error) {
	return fmt.Fprintf(
		pr.writer,
		"  \"%s\" [label=\"%s\", peripheries=2];\n",
		s.Id(), stateLabel(s),
	)
}

func (pr Printer) printDeadlocked(s State) (int, error) {
	return fmt.Fprintf(
		pr.writer,
		"  \"%s\" [label=\"%s\", fillcolor=\"#FFAAAA\", style=\"solid,filled\"];\n",
		s.Id(), stateLabel(s),
	)
}

func (pr Printer) printTransition(t Transition) (int, error) {
	return fmt.Fprintf(
		pr.writer,
		"  \"%s\" -> \"%s\" [label=\"%s.%s\"];\n",
		t.Source(), t.Target(), t.Process(), t.Label(),
	)
}

func (pr Printer) printTrace(t Transition) (int, error) {
	return fmt.Fprintf(
		pr.writer,
		"  \"%s\" -> \"%s\" [label=\"%s.%s\", color=\"#FF0000\", fontcolor=\"#FF0000\"];\n",
		t.Source(), t.Target(), t.Process(), t.Label(),
	)
}

func stateLabel(s State) string {
	ss := []string{}
	for pid, l := range s.Locations() {
		ss = append(ss, fmt.Sprintf("%s @ %s", pid, l))
	}
	vs := []string{}
	for x, n := range s.SharedVars() {
		vs = append(vs, fmt.Sprintf("%s = %d", x, n))
	}
	sort.Strings(ss)
	sort.Strings(vs)
	return strings.Join(ss, ", ") + "\\n" + strings.Join(vs, ", ")
}
