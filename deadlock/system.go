package deadlock

type process struct {
	id    processId
	rules map[location][]rule
}

type rule struct {
	label  label
	target location
	guard  guard
	action action
}

type System interface{}

type system struct {
	processes     []process
	initLocations locationMap
	initVars      sharedVars
}

type processId string
type location string
type locationMap map[processId]location

type label string
type guard func(sharedVars) sharedVars
type action func(sharedVars) sharedVars

type varName string

type sharedVars map[varName]int
