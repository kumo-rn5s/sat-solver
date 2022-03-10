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
	negative uint
	positive uint
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

func (c *cnf) deleteClause(clause *clause) {
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

func newClause(literals []int) *clause {
	return &clause{literals: append([]int{}, literals...)}
}

func parseLiterals(s string) ([]int, error) {
	raw := strings.Fields(s)
	literals := make([]int, 0, len(raw)-1)

	for _, v := range raw {
		if v == string(clauseEND) {
			break
		}
		l, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.New("wrong dimacs formats")
		}
		literals = append(literals, l)
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
		literals, err := parseLiterals(t)
		if err != nil {
			return err
		}
		c.push(newClause(literals))
	}
	return nil
}

func (c *cnf) deleteAllClausesByLiteral(literal int) {
	for p := c.head; p != nil; p = p.next {
		if _, found := p.findIndex(literal); found {
			c.deleteClause(p)
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

// リテラルLが1つだけの節があれば、Lを含む節を除去し、他の節の否定リテラル¬Lを消去する。
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
	pureLiterals := make([]int, 0, len(m))
	for k, v := range m {
		if v.positive > 0 && v.negative == 0 {
			pureLiterals = append(pureLiterals, k)
		} else if v.positive == 0 && v.negative > 0 {
			pureLiterals = append(pureLiterals, -k) // Key of purity map has literals' absolute value.
		}
	}
	return pureLiterals
}

// 節集合の中に肯定と否定の両方が現れないリテラル(純リテラル)Lがあれば、Lを含む節を除去する。
func (c *cnf) simplifyByPureLiteralRule() {
	literalsMap := c.makeLiteralsMap()
	pureLiterals := c.getPureLiterals(literalsMap)
	for _, l := range pureLiterals {
		c.deleteAllClausesByLiteral(l)
	}
}

func absInt(n int) int {
	if n > 0 {
		return n
	} else {
		return -n
	}
}

func findCountMaxLiteral(m map[int]*purity) int {
	var maxCount uint
	literal, _ := strconv.Atoi(string(clauseEND))
	for k, v := range m {
		if v.positive > maxCount {
			maxCount = v.positive
			literal = k
		}
		if v.negative > maxCount {
			maxCount = v.negative
			literal = -k // Key of purity map has literals' absolute value.
		}
	}
	return literal
}

func (c *cnf) getClausesMinLen() int {
	minLen := math.MaxInt
	for p := c.head; p != nil; p = p.next {
		len := len(p.literals)
		if len < minLen {
			minLen = len
		}
	}
	return minLen
}

func (c *cnf) getMinLenClauses(minLen int) *cnf {
	cnf := newCNF()
	for p := c.head; p != nil; p = p.next {
		if len(p.literals) == minLen {
			cnf.push(newClause(p.literals))
		}
	}
	return cnf
}

func (c *cnf) getAtomicFormula() int {
	minLen := c.getClausesMinLen()
	if minLen == 0 {
		log.Fatal("Illegal Length")
	}
	cnf := c.getMinLenClauses(minLen)
	return findCountMaxLiteral(cnf.makeLiteralsMap())
}

func (c *cnf) deepCopy() *cnf {
	cnf := newCNF()
	for p := c.head; p != nil; p = p.next {
		cnf.push(newClause(p.literals))
	}
	return cnf
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
	c2.push(newClause([]int{v}))
	if c2.isSatisfied() {
		return true
	}

	c3 := c.deepCopy()
	c3.push(newClause([]int{-v}))
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
