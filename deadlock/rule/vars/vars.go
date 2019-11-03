// Package vars provides variables shared by multiple processes.
package vars

type Name string

// Shared contains the values of variables at the system's each moment.
type Shared map[Name]int

func (vs Shared) Clone() Shared {
	c := Shared{}
	for x, n := range vs {
		c[x] = n
	}
	return c
}
