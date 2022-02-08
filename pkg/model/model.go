package model

// 一回しかない
type Preambles struct {
	Format       string
	VariablesNum int
	ClausesNum   int
}

type CNF struct {
	Preamble Preambles
	Clauses  List
}
