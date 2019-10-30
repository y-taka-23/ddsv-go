package deadlock_test

import (
	"errors"
	"testing"

	"github.com/y-taka-23/ddsv-go/deadlock"
	"github.com/y-taka-23/ddsv-go/deadlock/rule"
	"github.com/y-taka-23/ddsv-go/deadlock/rule/do"
	"github.com/y-taka-23/ddsv-go/deadlock/rule/vars"
)

func TestDetect(t *testing.T) {

	tests := []struct {
		name      string
		in        deadlock.System
		want      summary
		wantError bool
	}{
		{
			"1-step",
			deadlock.NewSystem().
				Register("P", deadlock.NewProcess().
					EnterAt("0").
					Define(rule.At("0").MoveTo("1"))),
			summary{state: 2, trans: 1, init: true, deadlock: 1},
			false,
		},
		{
			"2-step",
			deadlock.NewSystem().
				Register("P", deadlock.NewProcess().
					EnterAt("0").
					Define(rule.At("0").MoveTo("1")).
					Define(rule.At("1").MoveTo("2"))),
			summary{state: 3, trans: 2, init: true, deadlock: 1},
			false,
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
			false,
		},
		{
			"loop",
			deadlock.NewSystem().
				Register("P", deadlock.NewProcess().
					EnterAt("0").
					Define(rule.At("0").MoveTo("0"))),
			summary{state: 1, trans: 1, init: true, deadlock: 0},
			false,
		},
		{
			"declared var",
			deadlock.NewSystem().
				Declare(vars.Shared{"x": 0}).
				Register("P", deadlock.NewProcess().
					EnterAt("0").
					Define(rule.At("0").Let("", do.Set(1, "x")).MoveTo("1"))),
			summary{state: 2, trans: 1, init: true, deadlock: 1},
			false,
		},
		{
			"undeclared var",
			deadlock.NewSystem().
				Register("P", deadlock.NewProcess().
					EnterAt("0").
					Define(rule.At("0").Let("", do.Set(1, "x")).MoveTo("1"))),
			summary{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := deadlock.NewDetector().Detect(tt.in)
			if tt.wantError && errors.Is(err, nil) {
				t.Fatalf("want error, but has no error")
			}
			if !tt.wantError && !errors.Is(err, nil) {
				t.Fatalf("want no error, but has error %v", err)
			}
			if summarize(got) != tt.want {
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
