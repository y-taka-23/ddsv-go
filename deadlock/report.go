package deadlock

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/y-taka-23/ddsv-go/deadlock/rule"
)

type StateId string

type StateSet map[StateId]State

type LocationSet map[ProcessId]rule.Location

type State interface {
	Id() StateId
	Locations() LocationSet
	SharedVars() rule.SharedVars
}

type state struct {
	locations  LocationSet
	sharedVars rule.SharedVars
}

func (s state) Id() StateId {
	h := sha1.New()
	serial := fmt.Sprintf("%+v%+v", s.locations, s.sharedVars)
	h.Write([]byte(serial))
	return StateId(fmt.Sprintf("%x", h.Sum(nil)))
}

func (s state) Locations() LocationSet {
	return s.locations
}

func (s state) SharedVars() rule.SharedVars {
	return s.sharedVars
}

type Transition interface {
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
	Transited() []Transition
}

type report struct {
	visited   StateSet
	transited []Transition
}

func (rp report) Visited() StateSet {
	return rp.visited
}

func (rp report) Transited() []Transition {
	return rp.transited
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
		n, err := fmt.Fprintf(
			pr.writer,
			"  \"%s\" [label=\"%s\"]\n",
			s.Id(), stateLabel(s),
		)
		written += n
		if err != nil {
			return written, err
		}
	}
	for _, t := range rp.Transited() {
		n, err := fmt.Fprintf(
			pr.writer,
			"  \"%s\" -> \"%s\" [label=\"%s.%s\"]\n",
			t.Source(), t.Target(), t.Process(), t.Label(),
		)
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
