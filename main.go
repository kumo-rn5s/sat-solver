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
		f.Tail = n
	} else {
		f.Tail.next = n
		n.prev = f.Tail
		f.Tail = n
		f.Tail.next = nil
	}

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

// Parse CNF Methods
func isComment(s string) bool {
	return s[0:1] == "c"
}

func isPreamble(s string) bool {
	return s[0:1] == "p"
}

func (f *CNF) parseClause(s string) bool {
	clauseRaw := strings.Fields(s)
	var newClauseRaw = []int{}

	if clauseRaw[len(clauseRaw)-1] != "0" {
		return false
	}

	for _, i := range clauseRaw {
		if i == "0" {
			break
		}
		value, err := strconv.Atoi(i)
		if err != nil {
			return false
		}
		newClauseRaw = append(newClauseRaw, value)
	}
	f.Push(newClauseRaw)
	return true
}

func (f *CNF) Parse(filename string) error {
	// Define variables
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Start Scan file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw := scanner.Text()

		if len(raw) == 0 {
			continue
		} else if isComment(raw) {
			fmt.Println(raw)
		} else if isPreamble(raw) {
			preambles := strings.Fields(raw)
			if len(preambles) != 4 {
				return errors.New("wrong dimacs formats")
			}
			f.Preamble.Format = preambles[1]
			f.Preamble.VariablesNum, _ = strconv.Atoi(preambles[2])
			f.Preamble.ClausesNum, _ = strconv.Atoi(preambles[3])
		} else {
			if res := f.parseClause(raw); !res {
				return errors.New("wrong dimacs formats")
			}
		}
	}
	return nil
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

func absInt(value int) int {
	return int(math.Abs(float64(value)))
}

// moms heuristicへの準備
func getAtomicFormula(f *CNF) int {
	//出現回数記録
	variables := map[int]int{}
	for n := f.First(); n != nil; n = n.Next() {
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
	for n := f.First(); n != nil; n = n.Next() {
		newInt := []int{}
		newInt = append(newInt, n.Literals...)

		newFormula.Push(newInt)
	}
	return newFormula
}

func (f *CNF) hasEmptyclause() bool {
	for n := f.First(); n != nil; n = n.Next() {
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
	for i := 1; i < len(os.Args); i++ {
		formula := &CNF{}
		if err := formula.Parse(os.Args[i]); err != nil {
			log.Fatal("Parse Error")
		}
		if DPLL(formula) {
			fmt.Println("SATISFIABLE")
		} else {
			fmt.Println("UNSATISFIABLE")
		}
	}

}
