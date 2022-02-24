package main

import (
	"log"
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	f := &CNF{}
	testa := []int{1, 2, 3, 0}
	testb := []int{2, 3, 4, 0}
	testc := []int{3, 4, 5, 0}

	f.Push(testa)
	f.Push(testb)
	f.Push(testc)

	if !reflect.DeepEqual(f.Head.Literals, testa) ||
		!reflect.DeepEqual(f.Head.next.Literals, testb) ||
		!reflect.DeepEqual(f.Head.next.next.Literals, testc) {
		t.Error("Add Failure")
	}
}

func TestDeepCopy(t *testing.T) {
	f := &CNF{}
	testa := []int{2, -3, 0}
	testb := []int{-2, 0}
	testc := []int{3, 2, 0}

	f.Push(testa)
	f.Push(testb)
	f.Push(testc)

	newF := f.DeepCopy()
	log.Println(newF.Head.Literals)
	log.Println(newF.Head.next.Literals)
	log.Println(newF.Head.next.next.Literals)
	if !(f.Head != newF.Head && reflect.DeepEqual(f.Head.Literals, newF.Head.Literals)) {
		t.Error("DeepCopy1 Failure")
	}
	if !(f.Head.next != newF.Head.next && reflect.DeepEqual(f.Head.next.Literals, newF.Head.next.Literals)) {
		t.Error("DeepCopy2 Failure")
	}
	if !(f.Head.next.next != newF.Head.next.next && reflect.DeepEqual(f.Head.next.next.Literals, newF.Head.next.next.Literals)) {
		t.Error("DeepCopy3 Failure")
	}
}

func TestDelete(t *testing.T) {
	f := &CNF{}
	testa := []int{1, 2, 3, 0}
	testb := []int{2, 3, 4, 0}
	testc := []int{3, 4, 5, 0}

	f.Push(testa)
	f.Push(testb)
	f.Push(testc)

	f.Delete(f.Head.next)

	if !reflect.DeepEqual(f.Head.Literals, testa) ||
		!reflect.DeepEqual(f.Head.next.Literals, testc) ||
		!reflect.ValueOf(f.Head.next.next).IsNil() {
		t.Error("Delete Failure")
	}
}

func TestDeleteHead(t *testing.T) {
	f := &CNF{}
	testa := []int{1, 2, 3, 0}
	testb := []int{2, 3, 4, 0}
	testc := []int{3, 4, 5, 0}

	f.Push(testa)
	f.Push(testb)
	f.Push(testc)

	f.Delete(f.Head)

	if !reflect.DeepEqual(f.Head.Literals, testb) ||
		!reflect.DeepEqual(f.Head.next.Literals, testc) {
		t.Error("Delete Head Failure")
	}
}

func TestDeleteTail(t *testing.T) {
	f := &CNF{}
	testa := []int{1, 2, 3, 0}
	testb := []int{2, 3, 4, 0}
	testc := []int{3, 4, 5, 0}

	f.Push(testa)
	f.Push(testb)
	f.Push(testc)

	f.Delete(f.Tail)

	if !reflect.DeepEqual(f.Head.Literals, testa) ||
		!reflect.DeepEqual(f.Head.next.Literals, testb) {
		t.Error("Delete Tail Failure")
	}
}

func TestUnitElimination(t *testing.T) {
	formula := &CNF{}

	formula.Push([]int{1, 2, -3})
	formula.Push([]int{1, -2})
	formula.Push([]int{-1})
	formula.Push([]int{2, 3})

	unitElimination(formula)

	if !reflect.DeepEqual(formula.Head.Literals, []int{2, -3}) ||
		!reflect.DeepEqual(formula.Head.next.Literals, []int{-2}) ||
		!reflect.DeepEqual(formula.Head.next.next.Literals, []int{2, 3}) ||
		!reflect.ValueOf(formula.Head.next.next.next).IsNil() {
		t.Error("First Unit Elimination Failure")
	}

	unitElimination(formula)
	if !reflect.DeepEqual(formula.Head.Literals, []int{-3}) ||
		!reflect.DeepEqual(formula.Head.next.Literals, []int{3}) ||
		!reflect.ValueOf(formula.Head.next.next).IsNil() {
		t.Error("Second Unit Elimination Failure")
	}
	unitElimination(formula)
	if !reflect.DeepEqual(formula.Head.Literals, []int{}) ||
		!reflect.ValueOf(formula.Head.next).IsNil() {
		t.Error("Third Unit Elimination Failure")
	}
}

func TestPureElimination(t *testing.T) {
	formula := &CNF{}

	formula.Push([]int{1, 2})
	formula.Push([]int{-1, 2})
	formula.Push([]int{3, 4})
	formula.Push([]int{-3, -4})

	pureElimination(formula)

	if !reflect.DeepEqual(formula.Head.Literals, []int{3, 4}) ||
		!reflect.DeepEqual(formula.Head.next.Literals, []int{-3, -4}) ||
		!reflect.ValueOf(formula.Head.next.next).IsNil() {
		t.Error("First Pure Elimination Failure")
	}
}

func TestGetAtomicFormula(t *testing.T) {
	formula := &CNF{}

	formula.Push([]int{1, 2})
	formula.Push([]int{5, 4})
	formula.Push([]int{3, -5})
	formula.Push([]int{5, -6})

	result := getAtomicFormula(formula)

	if result != 5 {
		t.Error("Get Atomic Formula Error")
	}
}

func TestDPLL(t *testing.T) {
	formula := &CNF{}

	formula.Push([]int{1, 2, -3})
	formula.Push([]int{1, -2})
	formula.Push([]int{-1})
	formula.Push([]int{2, 3})
	if DPLL(formula) {
		t.Error("DPLL Error")
	}
}
