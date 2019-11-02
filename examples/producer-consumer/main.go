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

	capacity := 2

	base := func(mutex, me, you vars.Name) deadlock.Process {
		return deadlock.NewProcess().
			EnterAt("0").
			Define(rule.At("0").Only(when.Var(mutex).Is(0)).
				Let("lock", do.Set(1).ToVar(mutex)).MoveTo("1")).
			Define(rule.At("2").Only(when.Var(me).Is(1)).
				Let("wakeup", do.Set(0).ToVar(me)).MoveTo("0")).
			Define(rule.At("3").
				Let("signal", do.Set(1).ToVar(you)).MoveTo("4")).
			Define(rule.At("4").
				Let("unlock", do.Set(0).ToVar(mutex)).MoveTo("0"))
	}

	producer := func(buffer, mutex, me, you vars.Name) deadlock.Process {
		return base(mutex, me, you).
			Define(rule.At("1").Only(when.Var(buffer).Is(capacity)).
				Let("wait", do.Set(0).ToVar(mutex)).MoveTo("2")).
			Define(rule.At("1").Only(when.Var(buffer).IsLessThan(capacity)).
				Let("produce", do.Add(1).ToVar(buffer)).MoveTo("3"))
	}

	consumer := func(buffer, mutex, me, you vars.Name) deadlock.Process {
		return base(mutex, me, you).
			Define(rule.At("1").Only(when.Var(buffer).Is(0)).
				Let("wait", do.Set(0).ToVar(mutex)).MoveTo("2")).
			Define(rule.At("1").Only(when.Var(buffer).IsGreaterThan(0)).
				Let("consume", do.Add(-1).ToVar(buffer)).MoveTo("3"))
	}

	system := deadlock.NewSystem().
		Declare(vars.Shared{"buf": 0, "mut": 0, "prod": 0, "cons": 0}).
		Register("P", producer("buf", "mut", "prod", "cons")).
		Register("C", consumer("buf", "mut", "cons", "prod"))

	report, err := deadlock.NewDetector().Detect(system)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	_, err = deadlock.NewPrinter(os.Stdout).Print(report)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
