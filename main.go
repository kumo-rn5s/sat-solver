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
	if clause == c.head && clause == c.tail {
		c.head = nil
		c.tail = nil
	} else if clause == c.head {
		c.head = clause.next
		clause.next.prev = nil
	} else if clause == c.tail {
		c.tail = clause.prev
		clause.prev.next = nil
	} else {
		clause.prev.next, clause.next.prev = clause.next, clause.prev
	}
}

func (c *clause) findIndex(literal int) (int, bool) {
	for index, l := range c.literals {
		if l == literal {
			return index, true
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

func (c *cnf) createClause(l []int) *clause {
	return &clause{literals: append([]int{}, l...)}
}

func (c *cnf) parseLiterals(s string) ([]int, error) {
	var literals = make([]int, 0, len(s)-1)

	for _, v := range strings.Fields(s) {
		if v == string(clauseEND) {
			break
		}

		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.New("wrong dimacs formats")
		}
		literals = append(literals, num)
	}

	if literals == nil {
		return nil, errors.New("wrong dimacs formats")
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
		if literals, err := c.parseLiterals(t); err != nil {
			return err
		} else {
			clause := c.createClause(literals)
			c.push(clause)
		}
	}
	return nil
}

func (c *cnf) deleteClauseByTargetLiteral(target int) {
	for p := c.head; p != nil; p = p.next {
		if _, found := p.findIndex(target); found {
			c.delete(p)
		}
	}
}

func (c *cnf) deleteLiteralFromAllClause(target int) {
	for p := c.head; p != nil; p = p.next {
		if index, found := p.findIndex(target); found {
			p.removeLiteral(index)
		}
	}
}

// リテラル L 1つだけの節があれば、L を含む節を除去し、他の節の否定リテラル ¬L を消去する。
func simplifyByOneLiteralRule(c *cnf) {
	for p := c.head; p != nil; p = p.next {
		if len(p.literals) == 1 {
			c.deleteClauseByTargetLiteral(p.literals[0])
			c.deleteLiteralFromAllClause(-p.literals[0])
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

func (c *cnf) getPureClauseIndex(m map[int]*purity) []int {
	res := []int{}
	for k, v := range m {
		if v.positive == 0 && v.negative > 0 {
			res = append(res, -k)
		} else if v.positive > 0 && v.negative == 0 {
			res = append(res, k)
		}
	}
	return res
}

// 節集合の中に否定と肯定の両方が現れないリテラル（純リテラル） L があれば、L を含む節を除去する。
func simplifyByPureLiteralRule(c *cnf) {
	literalsMap := c.makeLiteralsMap()
	literals := c.getPureClauseIndex(literalsMap)

	for _, v := range literals {
		c.deleteClauseByTargetLiteral(v)
	}
}

func absInt(v int) int {
	return int(math.Abs(float64(v)))
}

func maxLiteralCount(literalsMap map[int]*purity) int {
	maxNumber := 0
	maxInt := -1
	for k, v := range literalsMap {
		if v.positive > maxNumber {
			maxNumber = v.positive
			maxInt = k
		}
		if v.negative > maxNumber {
			maxNumber = v.negative
			maxInt = k
		}

	}
	return maxInt
}

// moms heuristicへの準備
func (c *cnf) getAtomicFormula() int {
	return maxLiteralCount(c.makeLiteralsMap())
}

func (c *cnf) deepCopy() cnf {
	var new cnf
	for p := c.head; p != nil; p = p.next {
		clause := c.createClause(p.literals)
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
	simplifyByOneLiteralRule(c)
	simplifyByPureLiteralRule(c)

	if c.head == nil {
		return true
	}
	if c.hasEmptyClause() {
		return false
	}

	v := c.getAtomicFormula()

	c2 := c.deepCopy()
	clause := c2.createClause([]int{v})
	c2.push(clause)
	if c2.isSatisfied() {
		return true
	}

	c3 := c.deepCopy()
	clause = c3.createClause([]int{-v})
	c3.push(clause)
	return c3.isSatisfied()
}

func main() {
	if len(os.Args) == 1 {
		cnf := &cnf{}
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

			cnf := &cnf{}
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
