package deadlock

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Report interface {
	Visited() stateSet
	Transited() []transition
}

type report struct {
	visited   stateSet
	transited []transition
}

func (rp report) Visited() stateSet {
	return rp.visited
}

func (rp report) Transited() []transition {
	return rp.transited
}

type state struct {
	locations  locationMap
	sharedVars sharedVars
}

func (s state) id() stateId {
	h := sha1.New()
	serial := fmt.Sprintf("%+v%+v", s.locations, s.sharedVars)
	h.Write([]byte(serial))
	return stateId(fmt.Sprintf("%x", h.Sum(nil)))
}

func (s state) label() label {
	lbls := []string{}
	for pid, l := range s.locations {
		lbls = append(lbls, fmt.Sprintf("%s @ %s", pid, l))
	}
	sort.Strings(lbls)
	return label(strings.Join(lbls, ", "))
}

type transition struct {
	process processId
	label   label
	source  stateId
	target  stateId
}

type Printer interface {
	Print(Report) (int, error)
}

func NewPrinter(w io.Writer) Printer {
	return printer{writer: w}
}

type printer struct {
	writer io.Writer
}

func (pr printer) Print(rp Report) (int, error) {
	written, err := fmt.Fprintln(pr.writer, "digraph {")
	if err != nil {
		return written, err
	}
	for _, s := range rp.Visited() {
		n, err := fmt.Fprintf(
			pr.writer,
			"  \"%s\" [label=\"%s\"]\n",
			s.id(), s.label(),
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
			t.source, t.target, t.process, t.label,
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

type stateId string
type stateSet = map[stateId]state
