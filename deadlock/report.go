package deadlock

type Report interface {
}

type report struct {
	states      stateSet
	transitions []transition
	initial     stateId
	deadlocked  stateSet
}

type state struct {
	locations  locationMap
	sharedVars sharedVars
}

type transition struct {
	process processId
	label   label
	source  stateId
	target  stateId
}

type stateId string
type stateSet = map[stateId]state
