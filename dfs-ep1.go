package main

import (
	"fmt"
	"os"
)

type Link struct {
	from int
	to   int
}

type Node struct {
	id     int
	linkTo []int
}

type Network struct {
	nodes map[int]Node
	exits []int
}

func contains(n int, list []int) bool {
	for _, i := range list {
		if n == i {
			return true
		}
	}
	return false
}

func NewLink(n1, n2 int) Link {
	if n1 < n2 {
		return Link{n1, n2}
	} else {
		return Link{n2, n1}
	}
}

func (l Link) String() string {
	return fmt.Sprintf("%v %v", l.from, l.to)
}

func (n *Node) linkWith(otherNodeId int) {
	n.linkTo = append(n.linkTo, otherNodeId)
}

func (n *Node) cut(otherNodeId int) {
	newLinks := make([]int, 0)
	for _, l := range n.linkTo {
		if l != otherNodeId {
			newLinks = append(newLinks, l)
		}
	}
	n.linkTo = newLinks
}

func (n *Node) accessTo(nodes []int) []int {
	access := make([]int, 0)
	for _, ln := range n.linkTo {
		for _, e := range nodes {
			if ln == e {
				access = append(access, e)
			}
		}
	}
	return access
}

func (nw *Network) addNodeIfNotExist(n int) {
	if _, found := nw.nodes[n]; !found {
		nw.nodes[n] = Node{n, make([]int, 0)}
	}
}

func (nw *Network) linkNodes(n1 int, n2 int) {
	node := nw.nodes[n1]
	(&node).linkWith(n2)
	nw.nodes[n1] = node
	node = nw.nodes[n2]
	(&node).linkWith(n1)
	nw.nodes[n2] = node
}

func (nw *Network) cutAccessToDoubleExitsNode(endingNodes *[]int) (Link, bool) {
	accessToBlock := make([]int, 0)
	for _, n := range *endingNodes {
		if nw.countExit(n) > 1 {
			accessToBlock = append(accessToBlock, n)
		}
	}
	if len(accessToBlock) > 0 {
		return nw.cutPriorityLink(accessToBlock), true
	}
	return Link{}, false
}

func (nw *Network) addExit(e int) {
	nw.exits = append(nw.exits, e)
}

func (nw *Network) countExit(n int) (nbExit int) {
	nbExit = 0
	for _, a := range nw.nodes[n].linkTo {
		if contains(a, nw.exits) {
			nbExit++
		}
	}
	return
}

func (nw *Network) deepSearch(usedNodes, endingNodes *[]int) bool {
	fmt.Fprintf(os.Stderr, "deepSearch :\n-used : %v\n-end  :%v\n", *usedNodes, *endingNodes)
	futureEndingNodes := make([]int, 0)
	oneExitNode := make([]int, 0)
	for _, en := range *endingNodes {
		if nw.countExit(en) == 1 {
			oneExitNode = append(oneExitNode, en)
			for _, n := range nw.nodes[en].linkTo {
				if !contains(n, futureEndingNodes) && !contains(n, *usedNodes) && !contains(n, *endingNodes) {
					futureEndingNodes = append(futureEndingNodes, n)
				}
			}
		}
	}
	*usedNodes = append(*usedNodes, oneExitNode...)
	newEndingNode := make([]int, 0)
	for _, n := range *endingNodes {
		if !contains(n, oneExitNode) {
			newEndingNode = append(newEndingNode, n)
		}
	}
	newEndingNode = append(newEndingNode, futureEndingNodes...)
	*endingNodes = newEndingNode
	return len(futureEndingNodes) > 0
}

func (nw *Network) cut(n1, n2 int) Link {
	var node Node
	node = nw.nodes[n1]
	node.cut(n2)
	nw.nodes[n1] = node
	node = nw.nodes[n2]
	node.cut(n1)
	nw.nodes[n2] = node
	return Link{n1, n2}
}

func (nw *Network) cutNwPriorityLink() Link {
	nodeList := make([]int, 0)
	for n, _ := range nw.nodes {
		nodeList = append(nodeList, n)
	}
	return nw.cutPriorityLink(nodeList)
}

func (nw *Network) cutPriorityLink(nodes []int) Link {
	bestNode := nw.nodes[nodes[0]]
	bestExits := make([]int, 0)
	for _, n := range nodes {
		node := nw.nodes[n]
		exits := node.accessTo(nw.exits)
		if len(exits) > len(bestExits) {
			bestNode = node
			bestExits = exits
		}
	}
	if len(bestExits) > 0 {
		return nw.cut(bestNode.id, bestExits[0])
	}
	fmt.Fprint(os.Stderr, "Should never pass here\n")
	return Link{}
}

func (nw *Network) searchPriorityLink(usedNodes, endingNodes *[]int) Link {
	fmt.Fprintf(os.Stderr, "searchPriorityLink :\nused : %v\n ending : %v\n", *usedNodes, *endingNodes)
	if cut, found := nw.cutAccessToDoubleExitsNode(endingNodes); found {
		return cut
	}
	for nw.deepSearch(usedNodes, endingNodes) {
	}
	if cut, found := nw.cutAccessToDoubleExitsNode(endingNodes); found {
		return cut
	}
	futureEndingNodes := make([]int, 0)
	*usedNodes = append(*usedNodes, *endingNodes...)
	for _, nId := range *endingNodes {
		for _, nextNode := range nw.nodes[nId].linkTo {
			if !contains(nextNode, *usedNodes) {
				futureEndingNodes = append(futureEndingNodes, nextNode)
			}
		}
	}
	if len(futureEndingNodes) == 0 {
		return nw.cutNwPriorityLink()
	}
	endingNodes = &futureEndingNodes
	return nw.searchPriorityLink(usedNodes, endingNodes)
}

func (nw *Network) agentAt(a int) Link {
	agNode := nw.nodes[a]
	agExits := agNode.accessTo(nw.exits)
	if len(agExits) > 0 {
		return NewLink(a, agExits[0])
	}
	usedNodes := make([]int, 1)
	usedNodes[0] = agNode.id
	endingNodes := make([]int, len(agNode.linkTo))
	copy(endingNodes, agNode.linkTo)
	return nw.searchPriorityLink(&usedNodes, &endingNodes)
}

func main() {
	// N: the total number of nodes in the level, including the gateways
	// L: the number of links
	// E: the number of exit gateways
	var N, L, E int
	fmt.Scan(&N, &L, &E)
	network := Network{make(map[int]Node, N), make([]int, 0)}
	_ = network
	for i := 0; i < L; i++ {
		// N1: N1 and N2 defines a link between these nodes
		var N1, N2 int
		fmt.Scan(&N1, &N2)
		network.addNodeIfNotExist(N1)
		network.addNodeIfNotExist(N2)
		network.linkNodes(N1, N2)
	}

	for i := 0; i < E; i++ {
		// EI: the index of a gateway node
		var EI int
		fmt.Scan(&EI)
		network.addExit(EI)
	}

	fmt.Fprintf(os.Stderr, "%v\n", network)

	for {
		// SI: The index of the node on which the Bobnet agent is positioned this turn
		var SI int
		fmt.Scan(&SI)

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		// Example: 0 1 are the indices of the nodes you wish to sever the link between
		fmt.Println(network.agentAt(SI))
	}
}
