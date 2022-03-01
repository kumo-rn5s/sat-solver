package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type CNF struct {
	Head *Clause
	Tail *Clause
}

type Clause struct {
	Literals []int
	next     *Clause
	prev     *Clause
}

var CLAUSEEND = "0"

func (cnf *CNF) push(clause *Clause) {
	if cnf.Head == nil && cnf.Tail == nil {
		cnf.Head = clause
	} else {
		clause.prev = cnf.Tail
		clause.prev.next = clause
	}
	cnf.Tail = clause
}

func (cnf *CNF) delete(clause *Clause) {
	if clause == cnf.Head && clause == cnf.Tail {
		cnf.Head = nil
		cnf.Tail = nil
	} else if clause == cnf.Head {
		cnf.Head = clause.next
		clause.next.prev = nil
	} else if clause == cnf.Tail {
		cnf.Tail = clause.prev
		clause.prev.next = nil
	} else {
		clause.prev.next, clause.next.prev = clause.next, clause.prev
	}
}

func (c *Clause) findIndex(literal int) (int, bool) {
	for index, l := range c.Literals {
		if l == literal {
			return index, true
		}
	}
	return -1, false
}

func (c *Clause) remove(index int) {
	c.Literals = append(c.Literals[:index], c.Literals[index+1:]...)
}

func isSkipped(s string) bool {
	return len(s) == 0 || s[0] == '0' || s[0] == 'c' || s[0] == 'p' || s[0] == '%'
}

func (cnf *CNF) parseClause(s string) error {
	var literals = make([]int, 0, len(s)-1)

	for _, v := range strings.Fields(s) {
		if v == CLAUSEEND {
			break
		}

		num, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("wrong dimacs formats")
		}
		literals = append(literals, num)
	}

	if literals == nil {
		return errors.New("wrong dimacs formats")
	}

	cnf.push(&Clause{Literals: literals})
	return nil
}

func (cnf *CNF) Parse(f *os.File) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if isSkipped(t) {
			continue
		}
		if err := cnf.parseClause(t); err != nil {
			return err
		}
	}
	return nil
}

func (cnf *CNF) ShowCNF() {
	for p := cnf.Head; p != nil; p = p.next {
		fmt.Println(p)
	}
}

func (cnf *CNF) deleteClause(target int) {
	for p := cnf.Head; p != nil; p = p.next {
		if _, found := p.findIndex(target); found {
			cnf.delete(p)
		}
	}
}

func (cnf *CNF) deleteLiteral(target int) {
	for p := cnf.Head; p != nil; p = p.next {
		if index, found := p.findIndex(target); found {
			p.remove(index)
		}
	}
}

/*
1リテラル規則（one literal rule, unit rule）
リテラル L 1つだけの節があれば、L を含む節を除去し、他の節の否定リテラル ¬L を消去する。
*/
func simplifyByUnitRule(cnf *CNF) {
	for p := cnf.Head; p != nil; p = p.next {
		if len(p.Literals) == 1 {
			cnf.deleteClause(p.Literals[0])
			cnf.deleteLiteral(-p.Literals[0])
			p.next = cnf.Head
		}
	}
}

/*
純リテラル規則（pure literal rule, affirmative-nagative rule）
節集合の中に否定と肯定の両方が現れないリテラル（純リテラル） L があれば、L を含む節を除去する。
*/
type purity struct {
	positive bool
	negative bool
}

func upsertPurityMap(m map[int]purity, clause *Clause) map[int]purity {
	for _, v := range clause.Literals {
		newPurity := purity{}
		if old, ok := m[absInt(v)]; ok {
			newPurity = old
		}
		if v > 0 {
			newPurity.positive = true
		} else {
			newPurity.negative = true
		}
		m[absInt(v)] = newPurity
	}
	return m
}

func simplifyByPureRule(cnf *CNF) {
	literalPurityMap := make(map[int]purity)

	for p := cnf.Head; p != nil; p = p.next {
		literalPurityMap = upsertPurityMap(literalPurityMap, p)
	}

	for k, v := range literalPurityMap {
		if v.positive != v.negative {
			if v.positive {
				cnf.deleteClause(k)
			} else {
				cnf.deleteClause(-k)
			}
		}
	}
}

func absInt(v int) int {
	return int(math.Abs(float64(v)))
}

// moms heuristicへの準備
func getAtomicFormula(cnf *CNF) int {
	variables := map[int]int{}
	for p := cnf.Head; p != nil; p = p.next {
		for _, literal := range p.Literals {
			if v, ok := variables[absInt(literal)]; !ok {
				v++
				variables[absInt(literal)] = v
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

func (cnf *CNF) deepCopy() *CNF {
	newcnf := &CNF{}
	for p := cnf.Head; p != nil; p = p.next {
		newcnf.push(&Clause{Literals: append([]int{}, p.Literals...)})
	}
	return newcnf
}

func (cnf *CNF) hasEmptyclause() bool {
	for p := cnf.Head; p != nil; p = p.next {
		if len(p.Literals) == 0 {
			return true
		}
	}
	return false
}

func dpll(cnf *CNF) bool {
	simplifyByUnitRule(cnf)
	simplifyByPureRule(cnf)

	if cnf.Head == nil {
		return true
	}

	if cnf.hasEmptyclause() {
		return false
	}

	variable := getAtomicFormula(cnf)
	cnfBranch1 := cnf.deepCopy()

	cnfBranch1.push(&Clause{Literals: []int{variable}})
	if dpll(cnfBranch1) {
		return true
	}

	cnfBranch2 := cnf.deepCopy()
	cnfBranch2.push(&Clause{Literals: []int{-variable}})
	return dpll(cnfBranch2)
}

func (cnf *CNF) IsSatisfied() {
	if dpll(cnf) {
		fmt.Println("sat")
	} else {
		fmt.Println("unsat")
	}
}

func main() {
	if len(os.Args) == 1 {
		cnf := &CNF{}
		if err := cnf.Parse(os.Stdin); err != nil {
			log.Fatal("Parse Error")
		}

		cnf.IsSatisfied()
	} else {
		for i := 1; i < len(os.Args); i++ {
			f, err := os.Open(os.Args[i])
			if err != nil {
				log.Fatal("Parse Multi File Error")
			}
			defer f.Close()

			cnf := &CNF{}
			if err := cnf.Parse(f); err != nil {
				log.Fatal("Parse Error")
			}

			cnf.IsSatisfied()
		}
	}
}
