package model

type Claims struct {
	Key      int
	Literals []int
}

type Node struct {
	Claims
	next *Node
	prev *Node
}

func (n *Node) Next() *Node {
	return n.next
}

func (n *Node) Prev() *Node {
	return n.prev
}

func (n *Node) Find(literal int) bool {
	for _, l := range n.Literals {
		if l == literal {
			return true
		}
	}
	return false
}
