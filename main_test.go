package main

import (
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

func TestDelete(t *testing.T) {
	f := &CNF{}
	testa := []int{1, 2, 3, 0}
	testb := []int{2, 3, 4, 0}
	testc := []int{3, 4, 5, 0}

	f.Push(testa)
	f.Push(testb)
	f.Push(testc)

	f.Delete(f.Tail.prev)

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
		t.Error("Delete Failure")
	}
}

func TestFindClause(t *testing.T) {
	f := &CNF{}
	testa := []int{1, 2, 3, 0}

	f.Push(testa)

	if !f.Head.Find(1) || f.Head.Find(4) {
		t.Error("Clause Find Failure")
	}
}

func TestParse(t *testing.T) {
	filename := "./aim-100-1_6-no-1.cnf"
	formula, err := Parse(filename)
	if err != nil {
		t.Error("CNF Parse Failure")
	}
	if !(formula.Preamble.Format == "cnf") ||
		!(formula.Preamble.VariablesNum == 100) ||
		!(formula.Preamble.ClausesNum == 160) {
		t.Error("Preamble Parse Failure")
	}

	count := 0
	for n := formula.First(); n != nil; n = n.Next() {
		count++
	}
	if !(count == 160) {
		t.Error("Clause Parse Failure")
	}
}
