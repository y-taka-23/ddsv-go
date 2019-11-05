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

	capacity := 1

	waitConditionVar := func(mutex, cond vars.Name) do.Action {
		return func(vs vars.Shared) (vars.Shared, error) {
			newVars := vs.Clone()
			newVars[mutex] = 0
			newVars[cond] = 1
			return newVars, nil
		}
	}

	producer := func(queue, mutex, over, under vars.Name) deadlock.Process {
		return deadlock.NewProcess().
			EnterAt("0").
			Define(rule.At("0").Only(when.Var(mutex).Is(0)).
				Let("lock", do.Set(1).ToVar(mutex)).MoveTo("1")).
			Define(rule.At("1").Only(when.Var(queue).Is(capacity)).
				Let("wait", waitConditionVar(mutex, over)).MoveTo("3")).
			Define(rule.At("3").Only(when.Var(over).Is(0)).
				Let("wakeup", do.Nothing()).MoveTo("0")).
			Define(rule.At("1").Only(when.Var(queue).IsLessThan(capacity)).
				Let("produce", do.Add(1).ToVar(queue)).MoveTo("4")).
			Define(rule.At("4").
				Let("signal", do.Set(0).ToVar(under)).MoveTo("5")).
			Define(rule.At("5").
				Let("unlock", do.Set(0).ToVar(mutex)).MoveTo("0"))
	}

	consumer := func(queue, mutex, over, under vars.Name) deadlock.Process {
		return deadlock.NewProcess().
			EnterAt("0").
			Define(rule.At("0").Only(when.Var(mutex).Is(0)).
				Let("lock", do.Set(1).ToVar(mutex)).MoveTo("1")).
			Define(rule.At("1").Only(when.Var(queue).Is(0)).
				Let("wait", waitConditionVar(mutex, under)).MoveTo("3")).
			Define(rule.At("3").Only(when.Var(under).Is(0)).
				Let("wakeup", do.Nothing()).MoveTo("0")).
			Define(rule.At("1").Only(when.Var(queue).IsGreaterThan(0)).
				Let("consume", do.Add(-1).ToVar(queue)).MoveTo("4")).
			Define(rule.At("4").
				Let("signal", do.Set(0).ToVar(over)).MoveTo("5")).
			Define(rule.At("5").
				Let("unlock", do.Set(0).ToVar(mutex)).MoveTo("0"))
	}

	system := deadlock.NewSystem().
		Declare(vars.Shared{"que": 0, "mut": 0, "over": 0, "under": 0}).
		Register("P", producer("que", "mut", "over", "under")).
		Register("C", consumer("que", "mut", "over", "under"))

	report, err := deadlock.NewDetector().Detect(system)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	_, err = deadlock.NewPrinter(os.Stdout).Print(report)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
