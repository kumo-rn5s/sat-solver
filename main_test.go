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

func TestDelete(t *testing.T) {
	c := &cnf{}
	c.push(&clause{literals: []int{1, 2, 3, 0}})
	c.push(&clause{literals: []int{2, 3, 4, 0}})
	c.push(&clause{literals: []int{3, 4, 5, 0}})

	c.delete(c.head.next)

	if !reflect.DeepEqual(c.head.literals, []int{1, 2, 3, 0}) ||
		!reflect.DeepEqual(c.head.next.literals, []int{3, 4, 5, 0}) ||
		!reflect.ValueOf(c.head.next.next).IsNil() {
		t.Error("Delete Failure")
	}
}

func TestDeletehead(t *testing.T) {
	c := &cnf{}

	c.push(&clause{literals: []int{1, 2, 3, 0}})
	c.push(&clause{literals: []int{2, 3, 4, 0}})
	c.push(&clause{literals: []int{3, 4, 5, 0}})

	c.delete(c.head)

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

	c.delete(c.tail)

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
