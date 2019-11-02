package main

import (
	"fmt"
	"os"

	"github.com/y-taka-23/ddsv-go/deadlock"
	"github.com/y-taka-23/ddsv-go/deadlock/rule"
	"github.com/y-taka-23/ddsv-go/deadlock/rule/do"
	"github.com/y-taka-23/ddsv-go/deadlock/rule/vars"
	"github.com/y-taka-23/ddsv-go/deadlock/rule/when"
)

func main() {

	philo := func(me int, left, right vars.Name) deadlock.Process {
		return deadlock.NewProcess().
			EnterAt("0").
			Define(rule.At("0").Only(when.Var(left).Is(0)).
				Let("up_l", do.Set(me).ToVar(left)).MoveTo("1")).
			Define(rule.At("1").Only(when.Var(right).Is(0)).
				Let("up_r", do.Set(me).ToVar(right)).MoveTo("2")).
			// comment in the lines to avoid deadlocks
			//Define(rule.At("1").Only(when.Var(right).IsNot(0)).
			//	Let("down_l", do.Set(0).ToVar(left)).MoveTo("0")).
			Define(rule.At("2").Only(when.Var(right).Is(me)).
				Let("down_r", do.Set(0).ToVar(right)).MoveTo("3")).
			Define(rule.At("3").Only(when.Var(left).Is(me)).
				Let("down_l", do.Set(0).ToVar(left)).MoveTo("0"))
	}

	system := deadlock.NewSystem().
		Declare(vars.Shared{"f1": 0, "f2": 0}).
		Register("P1", philo(1, "f1", "f2")).
		Register("P2", philo(1, "f2", "f1"))

	report, err := deadlock.NewDetector().Detect(system)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	_, err = deadlock.NewPrinter(os.Stdout).Print(report)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
