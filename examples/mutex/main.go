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

	proc := func(global, local, mutex vars.Name) deadlock.Process {
		return deadlock.NewProcess().
			EnterAt("0").
			Define(rule.At("0").Only(when.Var(mutex).Is(0)).
				Let("lock", do.Set(1).ToVar(mutex)).MoveTo("1")).
			Define(rule.At("1").
				Let("read", do.CopyVar(global).ToVar(local)).MoveTo("2")).
			Define(rule.At("2").
				Let("incr", do.Add(1).ToVar(local)).MoveTo("3")).
			Define(rule.At("3").
				Let("write", do.CopyVar(local).ToVar(global)).MoveTo("4")).
			Define(rule.At("4").
				Let("unlock", do.Set(0).ToVar(mutex)).MoveTo("5")).
			HaltAt("5")
	}

	system := deadlock.NewSystem().
		Declare(vars.Shared{"var": 0, "tmp1": 0, "tmp2": 0, "mut": 0}).
		Register("P", proc("var", "tmp1", "mut")).
		Register("Q", proc("var", "tmp2", "mut"))

	report, err := deadlock.NewDetector().Detect(system)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	_, err = deadlock.NewPrinter(os.Stdout).Print(report)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
