package model

type List struct {
	head *Node
	tail *Node
}

func (l *List) First() *Node {
	return l.head
}

func (l *List) Push(v Claims) *List {
	n := &Node{Claims: v}
	if l.head == nil {
		l.head = n
	} else {
		l.tail.next = n
		n.prev = l.tail
	}
	l.tail = n
	return l
}

func (l *List) Find(key int) *Node {
	found := false
	var ret *Node = nil
	for n := l.First(); n != nil && !found; n = n.Next() {
		if n.Claims.Key == key {
			found = true
			ret = n
		}
	}
	return ret
}

func (l *List) Delete(key int) bool {
	success := false
	node2del := l.Find(key)
	if node2del != nil {
		prev_node := node2del.prev
		next_node := node2del.next

		//削除
		prev_node.next = node2del.next
		next_node.prev = node2del.prev
		success = true
	}
	return success
}
