package main

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2, 3, 0}})

	if !reflect.DeepEqual(c.head.literals, []int{1, 2, 3, 0}) {
		t.Error("Add Failure")
	}
}

func TestDeepCopy(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2, 3, 0}})

	newC := c.deepCopy()
	if !(c.head != newC.head && reflect.DeepEqual(c.head.literals, newC.head.literals)) ||
		!(c.tail != newC.tail && reflect.DeepEqual(c.tail.literals, newC.tail.literals)) {
		t.Error("DeepCopy Failure")
	}
}

func TestDeleteCommon(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2, 3, 0}})
	c.push(&clause{literals: []int{2, 3, 4, 0}})
	c.push(&clause{literals: []int{3, 4, 5, 0}})

	c.deleteClause(c.head.next)

	if !reflect.DeepEqual(c.head.literals, []int{1, 2, 3, 0}) ||
		!reflect.DeepEqual(c.head.next.literals, []int{3, 4, 5, 0}) ||
		!reflect.ValueOf(c.head.next.next).IsNil() {
		t.Error("Delete Failure")
	}
}

func TestDeleteHead(t *testing.T) {
	c := &cnf{}

	c.push(&clause{literals: []int{1, 2, 3, 0}})
	c.push(&clause{literals: []int{2, 3, 4, 0}})
	c.push(&clause{literals: []int{3, 4, 5, 0}})

	c.deleteClause(c.head)

	if !reflect.DeepEqual(c.head.literals, []int{2, 3, 4, 0}) ||
		!reflect.DeepEqual(c.head.next.literals, []int{3, 4, 5, 0}) {
		t.Error("Delete head Failure")
	}
}

func TestDeleteTail(t *testing.T) {
	c := &cnf{}

	c.push(&clause{literals: []int{1, 2, 3, 0}})
	c.push(&clause{literals: []int{2, 3, 4, 0}})
	c.push(&clause{literals: []int{3, 4, 5, 0}})

	c.deleteClause(c.tail)

	if !reflect.DeepEqual(c.head.literals, []int{1, 2, 3, 0}) ||
		!reflect.DeepEqual(c.head.next.literals, []int{2, 3, 4, 0}) {
		t.Error("delete Tail Failure")
	}
}

func TestSimplifyByUnitRule(t *testing.T) {
	c := &cnf{}

	c.push(&clause{literals: []int{1, 2, -3}})
	c.push(&clause{literals: []int{1, -2}})
	c.push(&clause{literals: []int{-1}})
	c.push(&clause{literals: []int{2, 3}})

	c.simplifyByOneLiteralRule()

	if !reflect.DeepEqual(c.head.literals, []int{}) ||
		!reflect.ValueOf(c.head.next).IsNil() {
		t.Error("Unit Elimination Failure")
	}
}

func TestSimplifyByPureRule(t *testing.T) {
	c := &cnf{}

	c.push(&clause{literals: []int{1, 2}})
	c.push(&clause{literals: []int{-1, 2}})
	c.push(&clause{literals: []int{3, 4}})
	c.push(&clause{literals: []int{-3, -4}})

	c.simplifyByPureLiteralRule()

	if !reflect.DeepEqual(c.head.literals, []int{3, 4}) ||
		!reflect.DeepEqual(c.head.next.literals, []int{-3, -4}) ||
		!reflect.ValueOf(c.head.next.next).IsNil() {
		t.Error("Pure Elimination Failure")
	}
}

func TestParseLiterals(t *testing.T) {
	raw := "1 2 3 4 0"
	if literals, err := parseLiterals(raw); err != nil || !reflect.DeepEqual(literals, []int{1, 2, 3, 4}) {
		t.Error("Parse Literals Failure")
	}
}

func TestRemoveLiteral(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2}})

	c.head.removeLiteral(1)
	if !reflect.DeepEqual(c.head.literals, []int{1}) {
		t.Error("Remove Literal Failure")
	}
}

func TestFindIndex(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2}})
	if i, found := c.head.findIndex(2); i != 1 || !found {
		t.Error("Find Index Failure")
	}
	if _, found := c.head.findIndex(3); found {
		t.Error("Find Index Failure")
	}
}

func TestHasEmptyClause(t *testing.T) {
	c := &cnf{}
	if c.hasEmptyClause() {
		t.Error("Has Empty Clause Failure")
	}
	c.push(&clause{literals: []int{}})
	if !c.hasEmptyClause() {
		t.Error("Has Empty Clause Failure")
	}
}

func TestGetAtomicFormula(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2, 3, 4}})
	c.push(&clause{literals: []int{1, 2, 4}})
	c.push(&clause{literals: []int{2, 3}})
	c.push(&clause{literals: []int{-2, 3}})
	if c.getAtomicFormula() != 3 {
		t.Error("Get Atomic Formula Failure")
	}
}

func TestGetClausesMinLen(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2, 3, 4}})
	c.push(&clause{literals: []int{1, 2, 4}})
	c.push(&clause{literals: []int{2, 3}})
	c.push(&clause{literals: []int{-2, 3}})

	length := c.getClausesMinLen()

	if length != 2 {
		t.Error("Get Min Len Clauses Failure")
	}
}

func TestGetMinLenClauses(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2, 3, 4}})
	c.push(&clause{literals: []int{1, 2, 4}})
	c.push(&clause{literals: []int{2, 3}})
	c.push(&clause{literals: []int{-2, 3}})

	nc := c.getMinLenClauses(c.getClausesMinLen())

	if !reflect.DeepEqual(nc.head.literals, []int{2, 3}) ||
		!reflect.DeepEqual(nc.head.next.literals, []int{-2, 3}) ||
		!reflect.ValueOf(nc.head.next.next).IsNil() {
		t.Error("Get Min Len Clauses Failure")
	}
}

func TestFindCountMaxLiteral(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2, 3, 4}})
	c.push(&clause{literals: []int{1, 2, 4}})
	c.push(&clause{literals: []int{2, -3}})
	c.push(&clause{literals: []int{2, 3}})

	if findCountMaxLiteral(c.makeLiteralsMap()) != 2 {
		t.Error("Find Count Max Literals Failure")
	}
}
func TestGetPureLiterals(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2, 3}})
	c.push(&clause{literals: []int{1, 2}})
	c.push(&clause{literals: []int{2, -3}})
	c.push(&clause{literals: []int{2, 3}})

	if !reflect.DeepEqual(c.getPureLiterals(c.makeLiteralsMap()), []int{1, 2}) {
		t.Error("Get Pure Literals Failure")
	}
}
func TestMakeLiteralsMap(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2, 3}})
	c.push(&clause{literals: []int{1, 2}})
	c.push(&clause{literals: []int{2, -3}})
	c.push(&clause{literals: []int{2, 3}})

	m := c.makeLiteralsMap()
	if m[1].negative != 0 &&
		m[1].positive != 2 &&
		m[2].negative != 0 &&
		m[2].positive != 4 &&
		m[3].negative != 1 &&
		m[3].positive != 2 {
		t.Error("Make Literals Map Failure")
	}

}
