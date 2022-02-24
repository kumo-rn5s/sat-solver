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

func (cnf *CNF) push(v []int) {
	n := &Clause{Literals: v}
	if cnf.Head == nil {
		cnf.Head = n
		cnf.Tail = n
	} else {
		cnf.Tail.next = n
		n.prev = cnf.Tail
		cnf.Tail = n
	}
}

func (cnf *CNF) delete(clause *Clause) {
	if clause == cnf.Head && clause == cnf.Tail {
		cnf.Head = nil
		cnf.Tail = nil
	} else if clause == cnf.Head {
		newHead := clause.next
		clause.next = nil
		cnf.Head = newHead
	} else if clause == cnf.Tail {
		newTail := clause.prev
		newTail.next = nil
		cnf.Tail = newTail
	} else if clause != nil {
		prev := clause.prev
		next := clause.next

		prev.next = clause.next
		next.prev = clause.prev
	}
}

func (c *Clause) find(literal int) int {
	res := -1
	for index, l := range c.Literals {
		if l == literal {
			res = index
		}
	}
	return res
}

func (c *Clause) remove(index int) {
	var literal []int
	if index == 0 {
		literal = c.Literals[1:]
	} else if index == len(c.Literals)-1 {
		literal = c.Literals[:len(c.Literals)-1]
	} else {
		literal = append(c.Literals[:index], c.Literals[index+1:]...)
	}
	c.Literals = literal
}

func isSkipped(s string) bool {
	return len(s) == 0 || s[0] == '0' || s[0] == 'c' || s[0] == 'p' || s[0] == '%'
}

func (cnf *CNF) parseClause(s string) error {
	var clause = make([]int, 0, len(s)-1)

	for _, v := range strings.Fields(s) {
		if v == "0" {
			break
		}

		num, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("wrong dimacs formats")
		}
		clause = append(clause, num)
	}

	if clause == nil {
		return errors.New("wrong dimacs formats")
	}

	cnf.push(clause)
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
	cnf.showCNF()
	return nil
}

func (cnf *CNF) showCNF() {
	for p := cnf.Head; p != nil; p = p.next {
		fmt.Println(p)
	}
}

/*
1リテラル規則（one literal rule, unit rule）
リテラル L 1つだけの節があれば、L を含む節を除去し、他の節の否定リテラル ¬L を消去する。
*/
func eliminateByUnitRule(cnf *CNF) {

	operation := map[*Clause][]int{}

	targetLiteral := 0
	for p := cnf.Head; p != nil; p = p.next {
		if len(p.Literals) == 1 {
			targetLiteral = p.Literals[0]
			break
		}
	}

	for n := cnf.Head; n != nil && targetLiteral != 0; n = n.next {
		//Lを含む節と¬Lを含む節に、Lと¬LのIndexを出力
		literalIndex := n.find(targetLiteral)
		literalNotIndex := n.find(targetLiteral * (-1))
		if literalIndex*literalNotIndex != 1 {
			operation[n] = []int{literalIndex, literalNotIndex}
		}
	}

	// 統一して削除
	for clause, v := range operation {
		if clause != nil {
			if v[1] != -1 {
				clause.remove(v[1])
			}
			if v[0] != -1 {
				cnf.delete(clause)
			}
		}
	}

	// temporary
	if len(operation) > 0 {
		eliminateByUnitRule(cnf)
	}
}

/*
純リテラル規則（pure literal rule, affirmative-nagative rule）
節集合の中に否定と肯定の両方が現れないリテラル（純リテラル） L があれば、L を含む節を除去する。
*/
func eliminateByPureRule(cnf *CNF) {

	operation := []*Clause{}
	literalMap := make(map[int]bool)
	// literal: true -> Pure
	// literal: false -> Not pure

	for p := cnf.Head; p != nil; p = p.next {
		for _, v := range p.Literals {
			if _, ok := literalMap[v*(-1)]; !ok {
				literalMap[v] = true
			} else {
				literalMap[v*(-1)] = false
			}
		}
	}

	for k, v := range literalMap {
		if v {
			for p := cnf.Head; p != nil; p = p.next {
				literalIndex := p.find(k)
				if literalIndex != -1 {
					operation = append(operation, p)
				}
			}
		}
	}
	// 統一して削除
	for _, clause := range operation {
		cnf.delete(clause)
	}
}

func absInt(v int) int {
	return int(math.Abs(float64(v)))
}

// moms heuristicへの準備
func getAtomicFormula(cnf *CNF) int {
	//出現回数記録
	variables := map[int]int{}
	for n := cnf.Head; n != nil; n = n.next {
		for _, literal := range n.Literals {
			// intを処理するabs()を実装する
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
	for n := cnf.Head; n != nil; n = n.next {
		newInt := []int{}
		newInt = append(newInt, n.Literals...)

		newcnf.push(newInt)
	}
	return newcnf
}

func (cnf *CNF) hasEmptyclause() bool {
	for n := cnf.Head; n != nil; n = n.next {
		if len(n.Literals) == 0 {
			return true
		}
	}
	return false
}

func DPLL(cnf *CNF) bool {
	eliminateByUnitRule(cnf)

	if cnf.Head == nil {
		return true
	}

	if cnf.hasEmptyclause() {
		return false
	}

	variable := getAtomicFormula(cnf)
	cnfBranch1 := cnf.deepCopy()

	cnfBranch1.push([]int{variable})
	if DPLL(cnfBranch1) {
		return true
	}

	cnfBranch2 := cnf.deepCopy()
	cnfBranch2.push([]int{variable * (-1)})
	if DPLL(cnfBranch2) {
		return true
	}

	return false
}

func main() {
	if len(os.Args) == 1 {
		cnf := &CNF{}
		if err := cnf.Parse(os.Stdin); err != nil {
			log.Fatal("Parse Error")
		}
		if DPLL(cnf) {
			fmt.Println("sat")
		} else {
			fmt.Println("unsat")
		}
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
			if DPLL(cnf) {
				fmt.Println("sat")
			} else {
				fmt.Println("unsat")
			}
		}
	}
}
