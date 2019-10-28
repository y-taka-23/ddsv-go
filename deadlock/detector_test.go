package deadlock_test

import (
	"testing"

	"github.com/y-taka-23/ddsv-go/deadlock"
	"github.com/y-taka-23/ddsv-go/deadlock/rule"
)

func TestDetect(t *testing.T) {

	tests := []struct {
		name string
		in   deadlock.System
		want summary
	}{
		{
			"1-step",
			deadlock.NewSystem().
				Register("P", deadlock.NewProcess().
					EnterAt("0").
					Define(rule.At("0").MoveTo("1"))),
			summary{state: 2, trans: 1, init: true, deadlock: 1},
		},
		{
			"2-step",
			deadlock.NewSystem().
				Register("P", deadlock.NewProcess().
					EnterAt("0").
					Define(rule.At("0").MoveTo("1")).
					Define(rule.At("1").MoveTo("2"))),
			summary{state: 3, trans: 2, init: true, deadlock: 1},
		},
		{
			"1-step 1-step",
			deadlock.NewSystem().
				Register("P", deadlock.NewProcess().
					EnterAt("0").
					Define(rule.At("0").MoveTo("1"))).
				Register("Q", deadlock.NewProcess().
					EnterAt("0").
					Define(rule.At("0").MoveTo("1"))),
			summary{state: 4, trans: 4, init: true, deadlock: 1},
		},
		{
			"loop",
			deadlock.NewSystem().
				Register("P", deadlock.NewProcess().
					EnterAt("0").
					Define(rule.At("0").MoveTo("0"))),
			summary{state: 1, trans: 1, init: true, deadlock: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := summarize(deadlock.NewDetector().Detect(tt.in))
			if got != tt.want {
				t.Fatalf("want %+v, but %+v", tt.want, got)
			}
		})
	}

}

type summary struct {
	state    int
	trans    int
	init     bool
	deadlock int
}

func summarize(rp deadlock.Report) summary {
	vs := rp.Visited()
	_, ok := vs[rp.Initial()]
	return summary{
		state:    len(vs),
		trans:    len(rp.Transited()),
		init:     ok,
		deadlock: len(rp.Deadlocked()),
	}
}
