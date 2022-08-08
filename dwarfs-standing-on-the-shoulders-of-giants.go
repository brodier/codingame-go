package main

import (
	"fmt"
	"os"
)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

type Tree map[int]*Node

type Node struct {
	id          int
	influencers []*Node
	influencees []*Node
	deepChan    chan int
	maxDeep     int
}

func NewNode(id int) *Node {
	n := Node{id: id}
	n.influencees = make([]*Node, 0)
	n.influencers = make([]*Node, 0)
	n.deepChan = make(chan int)
	n.maxDeep = 0
	return &n
}

func applyInfluence(from *Node, to *Node) {
	from.influencees = append(from.influencees, to)
	to.influencers = append(to.influencers, from)
}

func (n *Node) evalNbInfluencees(done chan int) {
	fmt.Fprintf(os.Stderr, "Start eval nb influencees for %d\n", n.id)
	if len(n.influencees) == 0 {
		for _, influencer := range n.influencers {
			fmt.Fprintf(os.Stderr, "sending one from %d to %d\n", n.id, influencer.id)
			influencer.deepChan <- 1
		}
	} else {
		var deep int
		for _, influencee := range n.influencees {
			deep = <-n.deepChan
			fmt.Fprintf(os.Stderr, "recieve %d from %d on %d\n", deep, influencee.id, n.id)
			if deep+1 > n.maxDeep {
				(*n).maxDeep = deep + 1
			}
		}
		for _, influencer := range n.influencers {
			influencer.deepChan <- (*n).maxDeep
		}
	}
	fmt.Fprintf(os.Stderr, "Done eval nb influencees for %d\n", n.id)
	done <- 1
}

func (n Node) String() string {
	var s string
	s = fmt.Sprintf("%d:>[", n.id)
	for _, i := range n.influencees {
		s += fmt.Sprintf("%d ", i.id)
	}
	s += "],<["
	for _, i := range n.influencers {
		s += fmt.Sprintf("%d ", i.id)
	}
	s += "]"
	s += fmt.Sprintf(" %d ", n.maxDeep)
	return s
}

func (t Tree) evalDeep() {
	done := make(chan int, len(t))
	for _, n := range t {
		go n.evalNbInfluencees(done)
	}
	for i := 0; i < len(t); i++ {
		_ = <-done
	}
}

func main() {
	// n: the number of relationships of influence
	var n int
	fmt.Scan(&n)
	tree := make(Tree, 0)
	for i := 0; i < n; i++ {
		// x: a relationship of influence between two people (x influences y)
		var x, y int
		fmt.Scan(&x, &y)
		if _, ok := tree[x]; !ok {
			tree[x] = NewNode(x)
		}
		if _, ok := tree[y]; !ok {
			tree[y] = NewNode(y)
		}
		applyInfluence(tree[x], tree[y])
	}
	fmt.Fprintf(os.Stderr, "%v\n", tree)
	tree.evalDeep()
	maxDeep := 0
	fmt.Fprintf(os.Stderr, "%v\n", tree)
	for _, n := range tree {
		if maxDeep < n.maxDeep {
			maxDeep = n.maxDeep
		}
	}
	fmt.Fprintf(os.Stderr, "%v\n", tree)

	// The number of people involved in the longest succession of influences
	fmt.Printf("%v\n", maxDeep)
}
