package model

type Preambles struct {
	Format       string
	VariablesNum int
	ClausesNum   int
}
type Clauses struct {
	Literal []int
}

type CNF struct {
	Preamble Preambles
	Clause   []Clauses
	ValueSet map[int]bool
}
