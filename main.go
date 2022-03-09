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

type cnf struct {
	head *clause
	tail *clause
}

type clause struct {
	literals []int
	next     *clause
	prev     *clause
}

type purity struct {
	negative int
	positive int
}

const (
	breakPoint = '%'
	clauseEND  = '0'
	comment    = 'c'
	preamble   = 'p'
)

type satSolver interface {
	isSatisfied() bool
}

var _ satSolver = (*cnf)(nil)

func newCNF() *cnf {
	return &cnf{}
}

func (c *cnf) push(clause *clause) {
	if c.head == nil && c.tail == nil {
		c.head = clause
	} else {
		clause.prev = c.tail
		clause.prev.next = clause
	}
	c.tail = clause
}

func (c *cnf) delete(clause *clause) {
	if clause != c.head && clause != c.tail {
		clause.prev.next = clause.next
		clause.next.prev = clause.prev
	} else {
		if clause == c.head {
			c.head = clause.next
			if c.head != nil {
				c.head.prev = nil
			}
		}
		if clause == c.tail {
			c.tail = clause.prev
			if c.tail != nil {
				c.tail.next = nil
			}
		}
	}
}

func (c *clause) findIndex(literal int) (int, bool) {
	for i, l := range c.literals {
		if l == literal {
			return i, true
		}
	}
	return 0, false
}

func (c *clause) removeLiteral(index int) {
	c.literals = append(c.literals[:index], c.literals[index+1:]...)
}

func isSkipped(s string) bool {
	return len(s) == 0 || s[0] == clauseEND || s[0] == comment || s[0] == preamble
}

func isBreakPoint(s string) bool {
	return s[0] == breakPoint
}

func createClause(literals []int) *clause {
	return &clause{literals: append([]int{}, literals...)}
}

func (c *cnf) parseLiterals(s string) ([]int, error) {
	raw := strings.Fields(s)
	literals := make([]int, len(raw)-1)

	for i, v := range raw {
		if v == string(clauseEND) {
			break
		}

		l, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.New("wrong dimacs formats")
		}
		literals[i] = l
	}
	return literals, nil
}

func (c *cnf) parseDIMACS(f *os.File) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if isSkipped(t) {
			continue
		}
		if isBreakPoint(t) {
			break
		}
		literals, err := c.parseLiterals(t)
		if err != nil {
			return err
		}
		clause := createClause(literals)
		c.push(clause)
	}
	return nil
}

func (c *cnf) deleteAllClausesByLiteral(literal int) {
	for p := c.head; p != nil; p = p.next {
		if _, found := p.findIndex(literal); found {
			c.delete(p)
		}
	}
}

func (c *cnf) deleteLiteralFromAllClauses(literal int) {
	for p := c.head; p != nil; p = p.next {
		if i, found := p.findIndex(literal); found {
			p.removeLiteral(i)
		}
	}
}

// リテラル L 1つだけの節があれば、L を含む節を除去し、他の節の否定リテラル ¬L を消去する。
func (c *cnf) simplifyByOneLiteralRule() {
	for p := c.head; p != nil; p = p.next {
		if len(p.literals) == 1 {
			c.deleteAllClausesByLiteral(p.literals[0])
			c.deleteLiteralFromAllClauses(-p.literals[0])
			p.next = c.head
		}
	}
}

func (c *cnf) makeLiteralsMap() map[int]*purity {
	m := make(map[int]*purity)

	for p := c.head; p != nil; p = p.next {
		for _, l := range p.literals {
			k := absInt(l)
			if _, ok := m[k]; !ok {
				m[k] = &purity{}
			}
			if l > 0 {
				m[k].positive++
			} else {
				m[k].negative++
			}
		}
	}
	return m
}

func (c *cnf) getPureLiterals(m map[int]*purity) []int {
	pureLiterals := []int{}
	for k, v := range m {
		if v.positive == 0 && v.negative > 0 {
			pureLiterals = append(pureLiterals, -k)
		} else if v.positive > 0 && v.negative == 0 {
			pureLiterals = append(pureLiterals, k)
		}
	}
	return pureLiterals
}

// 節集合の中に否定と肯定の両方が現れないリテラル（純リテラル） L があれば、L を含む節を除去する。
func (c *cnf) simplifyByPureLiteralRule() {
	literalsMap := c.makeLiteralsMap()
	pureLiterals := c.getPureLiterals(literalsMap)

	for _, l := range pureLiterals {
		c.deleteAllClausesByLiteral(l)
	}
}

func absInt(n int) int {
	return int(math.Abs(float64(n)))
}

func findCountMaxLiteral(m map[int]*purity) int {
	maxCount := 0
	literal, _ := strconv.Atoi(string(clauseEND))
	for k, v := range m {
		if v.positive >= maxCount {
			maxCount = v.positive
			literal = k
		}
		if v.negative >= maxCount {
			maxCount = v.negative
			literal = -k //Key of purity map has literals' absolute value
		}
	}
	return literal
}

func (c *cnf) getClausesMinLen() int {
	minCount := math.MaxInt
	for p := c.head; p != nil; p = p.next {
		length := len(p.literals)
		if length <= minCount {
			minCount = length
		}
	}
	return minCount
}

func (c *cnf) getMinLenClauses(minCount int) *cnf {
	nc := newCNF()
	for p := c.head; p != nil; p = p.next {
		if len(p.literals) == minCount {
			nc.push(createClause(p.literals))
		}
	}
	return nc
}

func (c *cnf) getAtomicFormula() int {
	nc := c.getMinLenClauses(c.getClausesMinLen())
	return findCountMaxLiteral(nc.makeLiteralsMap())
}

func (c *cnf) deepCopy() *cnf {
	new := newCNF()
	for p := c.head; p != nil; p = p.next {
		clause := createClause(p.literals)
		new.push(clause)
	}
	return new
}

func (c *cnf) hasEmptyClause() bool {
	for p := c.head; p != nil; p = p.next {
		if len(p.literals) == 0 {
			return true
		}
	}
	return false
}

func (c *cnf) isSatisfied() bool {
	c.simplifyByOneLiteralRule()
	c.simplifyByPureLiteralRule()

	if c.head == nil {
		return true
	}
	if c.hasEmptyClause() {
		return false
	}

	v := c.getAtomicFormula()

	c2 := c.deepCopy()
	clause := createClause([]int{v})
	c2.push(clause)
	if c2.isSatisfied() {
		return true
	}

	c3 := c.deepCopy()
	clause = createClause([]int{-v})
	c3.push(clause)
	return c3.isSatisfied()
}

func main() {
	if len(os.Args) == 1 {
		cnf := newCNF()
		if err := cnf.parseDIMACS(os.Stdin); err != nil {
			log.Fatal("Parse Error")
		}
		if cnf.isSatisfied() {
			fmt.Println("sat")
		} else {
			fmt.Println("unsat")
		}
	} else {
		for i := 1; i < len(os.Args); i++ {
			f, err := os.Open(os.Args[i])
			if err != nil {
				log.Fatal("Parse Multiple File Error")
			}
			defer f.Close()

			cnf := newCNF()
			if err := cnf.parseDIMACS(f); err != nil {
				log.Fatal("Parse Error")
			}

			if cnf.isSatisfied() {
				fmt.Println("sat")
			} else {
				fmt.Println("unsat")
			}
		}
	}
}
