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

// Formula List Methods
func (f *CNF) First() *Clause {
	return f.Head
}

func (f *CNF) Push(v []int) {
	n := &Clause{Literals: v}
	if f.Head == nil {
		f.Head = n
		f.Tail = n
	} else {
		f.Tail.next = n
		n.prev = f.Tail
		f.Tail = n
	}

}

func (f *CNF) Delete(clause *Clause) {
	if clause == f.Head && clause == f.Tail {
		f.Head = nil
		f.Tail = nil
	} else if clause == f.Head {
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

func (cnf *CNF) parseClause(s string) []int {
	var clause = make([]int, 0, len(s)-1)

	for _, l := range strings.Fields(s) {
		if l == "0" {
			break
		}

		v, err := strconv.Atoi(l)
		if err != nil {
			return nil
		}
		clause = append(clause, v)
	}
	return clause
}

func (cnf *CNF) Parse(f *os.File) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if isSkipped(t) {
			continue
		}
		clause := cnf.parseClause(t)
		if clause == nil {
			return errors.New("wrong dimacs formats")
		}
		cnf.Push(clause)
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
func unitElimination(formula *CNF) {

	operation := map[*Clause][]int{}

	targetLiteral := 0
	for n := formula.First(); n != nil; n = n.next {
		if len(n.Literals) == 1 {
			targetLiteral = n.Literals[0]
			break
		}
	}

	for n := formula.First(); n != nil && targetLiteral != 0; n = n.next {
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

	// temporary
	if len(operation) > 0 {
		unitElimination(formula)
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

	for n := formula.First(); n != nil; n = n.next {
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
			for n := formula.First(); n != nil; n = n.next {
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

func absInt(value int) int {
	return int(math.Abs(float64(value)))
}

// moms heuristicへの準備
func getAtomicFormula(f *CNF) int {
	//出現回数記録
	variables := map[int]int{}
	for n := f.First(); n != nil; n = n.next {
		for _, literal := range n.Literals {
			// intを処理するabs()を実装する
			if value, ok := variables[absInt(literal)]; !ok {
				value++
				variables[absInt(literal)] = value
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
	for n := f.First(); n != nil; n = n.next {
		newInt := []int{}
		newInt = append(newInt, n.Literals...)

		newFormula.Push(newInt)
	}
	return newFormula
}

func (f *CNF) hasEmptyclause() bool {
	for n := f.First(); n != nil; n = n.next {
		if len(n.Literals) == 0 {
			return true
		}
	}
	return false
}

func DPLL(formula *CNF) bool {
	unitElimination(formula)

	if formula.Head == nil {
		return true
	}

	if formula.hasEmptyclause() {
		return false
	}

	variable := getAtomicFormula(formula)
	formulaBranch1 := formula.DeepCopy()

	formulaBranch1.Push([]int{variable})
	if DPLL(formulaBranch1) {
		return true
	}

	formulaBranch2 := formula.DeepCopy()
	formulaBranch2.Push([]int{variable * (-1)})
	if DPLL(formulaBranch2) {
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
				log.Fatal("hoge")
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
