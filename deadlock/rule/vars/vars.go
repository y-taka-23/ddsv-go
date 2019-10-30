package vars

type Name string

type Shared map[Name]int

func (vs Shared) Clone() Shared {
	c := Shared{}
	for x, n := range vs {
		c[x] = n
	}
	return c
}
