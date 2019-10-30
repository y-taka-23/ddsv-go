package do

import (
	"fmt"

	"github.com/y-taka-23/ddsv-go/deadlock/rule/vars"
)

type Action func(vars.Shared) (vars.Shared, error)

func Copy(y, x vars.Name) Action {
	return func(vs vars.Shared) (vars.Shared, error) {
		modified := vs.Clone()
		if _, ok := vs[y]; !ok {
			return vars.Shared{}, fmt.Errorf("undeclared variable: %s", y)
		}
		if _, ok := modified[x]; !ok {
			return vars.Shared{}, fmt.Errorf("undeclared variable: %s", x)
		}
		modified[x] = vs[y]
		return modified, nil
	}
}

func Set(n int, x vars.Name) Action {
	return func(vs vars.Shared) (vars.Shared, error) {
		modified := vs.Clone()
		if _, ok := modified[x]; !ok {
			return vars.Shared{}, fmt.Errorf("undeclared variable: %s", x)
		}
		modified[x] = n
		return modified, nil
	}
}

func Add(n int, x vars.Name) Action {
	return func(vs vars.Shared) (vars.Shared, error) {
		modified := vs.Clone()
		if _, ok := modified[x]; !ok {
			return vars.Shared{}, fmt.Errorf("undeclared variable: %s", x)
		}
		modified[x] = vs[x] + n
		return modified, nil
	}
}
