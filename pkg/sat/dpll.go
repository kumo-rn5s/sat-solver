package sat

import (
	"github.com/FirosStuart/SAT-solver/pkg/model"
)

/*
1リテラル規則（one literal rule, unit rule）
リテラル L 1つだけの節があれば、L を含む節を除去し、他の節の否定リテラル ¬L を消去する。
*/
func unitElimination(formula *model.CNF) {
	target := 0
	var targetIndex []int

	for index, clause := range formula.Clause {
		if target == 0 && len(clause.Literal) == 1 {
			target = clause.Literal[0]
			targetIndex = append(targetIndex, index)
		}

		if target != 0 {
			if contains(clause.Literal, target) || contains(clause.Literal, -1*target) {
				targetIndex = append(targetIndex, index)
			}
		}
	}
	formula.Clause = remove(formula.Clause, targetIndex)
}

/*
純リテラル規則（pure literal rule, affirmative-nagative rule）
節集合の中に否定と肯定の両方が現れないリテラル（純リテラル） L があれば、L を含む節を除去する。
*/
func pureElimination(formula *model.CNF) {
	target := 0
	var targetIndex []int

	for index, clause := range formula.Clause {
		if target == 0 && len(clause.Literal) == 1 {
			target = clause.Literal[0]
			targetIndex = append(targetIndex, index)
		}

		if target != 0 {
			if !contains(clause.Literal, -1*target) {
				targetIndex = append(targetIndex, index)
			}
		}
	}
	formula.Clause = remove(formula.Clause, targetIndex)
}

/*
分割規則（splitting rule, rule of case analysis）
節集合 F の中に否定と肯定の両方があるリテラル L があれば、そのリテラルを真偽に解釈してえられる2つの節集合に分割する。
*/
func splitting(formula *model.CNF) {

}

/*
DPLL(F):
 1リテラル規則、純リテラル規則などを使い F を単純化
 if F is 空:
     return "充足可能"
 if F is 空節を含む:
     return "充足不能"
 原子論理式 v を選択
 真理値 b を選択 (true or false)
 if DPLL(v = b とした F ) is "充足可能":
     return "充足可能"
 if DPLL(v = ¬b とした F) is "充足可能":
     return "充足可能"
 return "充足不能"
*/

func DPLL(formula model.CNF) bool {
	unitElimination(&formula)
	pureElimination(&formula)
	//splitting(&formula)

	if len(formula.Clause) == 0 {
		return true
	}

	nowVariables := getNowLiteral(formula)

	for literal := range nowVariables {
		formula.ValueSet[literal] = true
		if DPLL(formula) {
			return true
		}
		formula.ValueSet[literal] = false
		if DPLL(formula) {
			return true
		}
	}
	return false
}

func getNowLiteral(formula model.CNF) []int {
	return nil
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func remove(clauses []model.Clauses, indexs []int) []model.Clauses {
	var newClauses []model.Clauses

	for i := 0; i < len(clauses); i++ {
		if !contains(indexs, i) {
			newClauses = append(newClauses, clauses[i])
		}
	}
	return newClauses
}

/*
c Example CNF format file c
p cnf 4 3
1 3 -4 0
4 0
2 -3 0
*/
