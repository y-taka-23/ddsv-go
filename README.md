Toy Deadlock Detector
=====================

[![GoDoc](https://godoc.org/github.com/y-taka-23/ddsv-go/deadlock?status.svg)](https://godoc.org/github.com/y-taka-23/ddsv-go/deadlock)
[![CircleCI](https://circleci.com/gh/y-taka-23/ddsv-go.svg?style=svg)](https://circleci.com/gh/y-taka-23/ddsv-go)

This package aims to provide a DSL to represent processes as [finate state machines](https://en.wikipedia.org/wiki/Finite-state_machine) and their concurrent composition. A detector traverses all possible states of the concurrent system, and reports on deadlocks, namely states in which no process can take the next step. Also, the package provides [Graphviz](https://www.graphviz.org/) style outputs, so you can intuitively view the state space of your system.


Example: Dining Philosophers Problem
------------------------------------

The [dining philosophers problem](https://en.wikipedia.org/wiki/Dining_philosophers_problem) is one of the best-known examples of conccerent programming. In this model, some philosophers are sitting on a round table and forks are served between each philosophers. A pasta bawl is also served at the centre of the table, but philosophers have to hold both of left/right forks to help themselves. Here the philosophers are analogues of processes/threads, and the forks are that of shared resources.

<img src="/assets/table.png" width=300px align="right" alt="philosophers and forks around a table">

In a naive implementation of this setting, for example, all philosophers act as following:

1. Pick up the fork in his left side
2. Pick up the fork in his right side
3. Eat the pasta
4. Put down the fork in his right hand
5. Put down the fork in his left hand

When multiple philosophers act like this concurrently, as you noticed, it results in a __deadlock__. Let's model the situation and detect the deadlocked state by this package.

As the simpleerest case, assume that only two philosophers sitting on the table. We define two processes `P1`, `P2` to represent the philosophers, and two shared variables `f1`, `f2` for forks. The fork `f1` is in `P1`'s left side, and the `f2` is in his right side.



```golang
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
			// Pick up the fork in his left side
			Define(rule.At("0").Only(when.Var(left).Is(0)).
				Let("up_l", do.Set(me).ToVar(left)).MoveTo("1")).
			// Pick up the fork in his right side
			Define(rule.At("1").Only(when.Var(right).Is(0)).
				Let("up_r", do.Set(me).ToVar(right)).MoveTo("2")).
			// Put down the fork in his right side
			Define(rule.At("2").Only(when.Var(right).Is(me)).
				Let("down_r", do.Set(0).ToVar(right)).MoveTo("3")).
			// Put down the fork in his left side
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
```

<img src="/assets/trace_bad.png" height=500px alt="transition graph which has a deadlocked state">

The graph shows an error case in which `P1` gets `f1` then `P2` gets `f2`. The red state is a deadlock, `P1` waits `f2` and `P2` waits `f1` respectively forever.

Then, how can we solve the deadlock problem? One idea is to let philosophers put down his first fork, if his second fork is occupied by another philosopher, and try again. Add the following lines in the definition of `philo`. Run the detector again, and you see the deadlock state disappears.

```golang
// Discard the fork in his left side
Define(rule.At("1").Only(when.Var(right).IsNot(0)).
	Let("down_l", do.Set(0).ToVar(left)).MoveTo("0")).
```

<img src="/assets/trace_good.png" height=500px alt="transition graph without the deadlock">

More examples are demonstrated in [examples](/examples) directory. Check it out!

Acknowledgement
---------------

* [Multi-thread programming by DYI deadlock detector]() by [@hatsugai](https://github.com/hatsugai)
