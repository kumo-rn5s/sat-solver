package main

import (
	"bufio"
	"errors"
	"log"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Preamble struct {
	Format       string
	VariablesNum int
	ClausesNum   int
}

type CNF struct {
	Preamble Preamble
	Head     *Clause
	Tail     *Clause
}

type Clause struct {
	Literals []int
	next     *Clause
	prev     *Clause
}

// Formula List Methods
func (f *CNF) First() *Clause {
	return f.Head
}

func (f *CNF) Push(v []int) {
	n := &Clause{Literals: v}
	if f.Head == nil {
		f.Head = n
	} else {
		f.Tail.next = n
		n.prev = f.Tail
	}
	f.Tail = n
}

func (f *CNF) Delete(clause *Clause) {
	if clause == f.Head {
		newHead := clause.next
		clause.next = nil
		f.Head = newHead
	} else if clause == f.Tail {
		newTail := clause.prev
		newTail.next = nil
		f.Tail = newTail
	} else if clause != nil {
		prev := clause.prev
		next := clause.next

		prev.next = clause.next
		next.prev = clause.prev
	}
}

// Clause Node Methods
func (c *Clause) Next() *Clause {
	return c.next
}

func (c *Clause) Prev() *Clause {
	return c.prev
}

func (c *Clause) Find(literal int) int {
	res := -1
	for index, l := range c.Literals {
		if l == literal {
			res = index
		}
	}
	return res
}

func (c *Clause) Remove(index int) {
	literal := append(c.Literals[:index], c.Literals[index+1:]...)
	c.Literals = literal
}

// Parse CNF Methods
func isComment(s string) bool {
	return s[0:1] == "c"
}

func isPreamble(s string) bool {
	return s[0:1] == "p"
}

func isClause(s string) ([]int, bool) {
	clauseRaw := strings.Fields(s)
	var newClauseRaw = []int{}

	if clauseRaw[len(clauseRaw)-1] != "0" {
		return nil, false
	}
	for i := 0; i < len(clauseRaw)-1; i++ {
		value, err := strconv.Atoi(clauseRaw[i])
		newClauseRaw = append(newClauseRaw, value)
		if err != nil {
			return nil, false
		}
	}
	return newClauseRaw, true
}

func Parse(filename string) (*CNF, error) {
	// Define variables
	formulas := &CNF{}

	// Open file
	fullFilename, _ := filepath.Abs(filename)
	file, err := os.Open(fullFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Start Scan file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw := scanner.Text()

		if len(raw) == 0 {
			continue
		}
		if isComment(raw) {
			log.Println(raw)
		}
		if isPreamble(raw) {
			preambles := strings.Fields(raw)
			if len(preambles) != 4 {
				return nil, errors.New("wrong dimacs formats")
			}
			formulas.Preamble.Format = preambles[1]
			formulas.Preamble.VariablesNum, _ = strconv.Atoi(preambles[2])
			formulas.Preamble.ClausesNum, _ = strconv.Atoi(preambles[3])
		}
		if clauseRaw, result := isClause(raw); result {
			formulas.Push(clauseRaw)
		}
	}
	return formulas, nil
}

/*
1リテラル規則（one literal rule, unit rule）
リテラル L 1つだけの節があれば、L を含む節を除去し、他の節の否定リテラル ¬L を消去する。
*/
func unitElimination(formula *CNF) {

	operation := map[*Clause][]int{}

	targetLiteral := 0
	for n := formula.First(); n != nil; n = n.Next() {
		if len(n.Literals) == 1 {
			targetLiteral = n.Literals[0]
			break
		}
	}

	for n := formula.First(); n != nil && targetLiteral != 0; n = n.Next() {
		//Lを含む節と¬Lを含む節に、Lと¬LのIndexを出力
		literalIndex := n.Find(targetLiteral)
		literalNotIndex := n.Find(targetLiteral * (-1))
		if literalIndex*literalNotIndex != 1 {
			operation[n] = []int{literalIndex, literalNotIndex}
		}
	}

	// 統一して削除
	for clause, value := range operation {
		if clause != nil {
			if value[1] != -1 {
				clause.Remove(value[1])
			}
			if value[0] != -1 {
				formula.Delete(clause)
			}
		}
	}
}

/*
純リテラル規則（pure literal rule, affirmative-nagative rule）
節集合の中に否定と肯定の両方が現れないリテラル（純リテラル） L があれば、L を含む節を除去する。
*/
func pureElimination(formula *CNF) {

	operation := []*Clause{}
	literalMap := make(map[int]bool)
	// literal: true -> Pure
	// literal: false -> Not pure

	for n := formula.First(); n != nil; n = n.Next() {
		for _, l := range n.Literals {
			if _, ok := literalMap[l*(-1)]; !ok {
				literalMap[l] = true
			} else {
				literalMap[l*(-1)] = false
			}
		}
	}

	for key, value := range literalMap {
		if value {
			for n := formula.First(); n != nil; n = n.Next() {
				literalIndex := n.Find(key)
				if literalIndex != -1 {
					operation = append(operation, n)
				}
			}
		}
	}
	// 統一して削除
	for _, clause := range operation {
		formula.Delete(clause)
	}
}

// moms heuristicへの準備
func getAtomicFormula(f *CNF) int {
	//出現回数記録
	variables := map[int]int{}
	for n := f.First(); n != nil; n = n.Next() {
		for _, literal := range n.Literals {
			if value, ok := variables[int(math.Abs(float64(literal)))]; !ok {
				variables[literal] = 1
			} else {
				variables[literal] = value + 1
			}
		}
	}
	maxNumber := 0
	maxInt := -1
	for i, v := range variables {
		if v > maxNumber {
			maxNumber = v
			maxInt = i
		}
	}
	return maxInt
}

func (f *CNF) DeepCopy() *CNF {
	newFormula := &CNF{}
	for n := f.First(); n != nil; n = n.Next() {
		newFormula.Push(n.Literals)
	}
	return newFormula
}

func DPLL(formula *CNF) bool {
	unitElimination(formula)
	//pureElimination(formula)
	//splitting(&formula)

	if formula.Head == nil {
		return true
	}

	for n := formula.First(); n != nil; n = n.Next() {
		if reflect.DeepEqual(n.Literals, []int{}) {
			return false
		}
	}

	variable := getAtomicFormula(formula)

	formulaBranch1 := formula.DeepCopy()
	formulaBranch2 := formula.DeepCopy()

	formulaBranch1.Push([]int{variable})
	if DPLL(formulaBranch1) {
		return true
	}

	formulaBranch2.Push([]int{variable * (-1)})
	if DPLL(formulaBranch2) {
		return true
	}

	return false
}

func main() {
	filename := "./test-unsat-1.cnf"
	formula, err := Parse(filename)
	if err != nil {
		log.Println("UNSATISFIABLE")
		os.Exit(1)
	}
	if DPLL(formula) {
		log.Println("SATISFIABLE")
	} else {
		log.Println("UNSATISFIABLE")
	}
}
