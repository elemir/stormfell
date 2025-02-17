package model

type Terrain struct {
	Name        string
	Transitions map[string][]int
}
