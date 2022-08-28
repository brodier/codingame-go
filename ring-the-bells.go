package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Cycle []int

type Permutation []Cycle

type void struct{}
type Set map[int]void

var member void

type BellInPermut struct {
	bellId       int
	bellPosition int
}

func (c Cycle) getNewPos(currPos int) (newPos int) {
	newPos = currPos
	for i, p := range c {
		if p == currPos {
			if i == len(c)-1 {
				return c[0]
			} else {
				return c[i+1]
			}
		}
	}
	return
}

func Compile(p Permutation) []BellInPermut {
	bellsSet := make(Set, 0)
	for _, c := range p {
		for _, b := range c {
			bellsSet[b] = member
		}
	}
	resultPermut := make([]BellInPermut, 0)
	for bell := range bellsSet {
		bellPermut := BellInPermut{bell, bell}
		lastPermut := len(p) - 1
		for i := range p {
			bellPermut.bellPosition = p[lastPermut-i].getNewPos(bellPermut.bellPosition)
		}
		resultPermut = append(resultPermut, bellPermut)
	}
	return resultPermut
}

func ParsePermutation(permut string) Permutation {
	permutation := make(Permutation, 0)
	curCycle := make(Cycle, 0)
	for _, r := range permut {
		switch r {
		case '(':
			curCycle = make(Cycle, 0)
		case ')':
			permutation = append(permutation, curCycle)
		case ' ':
			continue
		default:
			curCycle = append(curCycle, int(r-'0'))
		}
	}
	fmt.Fprintf(os.Stderr, "Permut : %v\n", permutation)
	return permutation
}

func BuildCycle(from, to int, bellMoves []BellInPermut, cycle *Cycle) {
	if from == to {
		return
	}
	for _, bm := range bellMoves {
		if bm.bellId == from {
			*cycle = append(*cycle, bm.bellId)
			BuildCycle(bm.bellPosition, to, bellMoves, cycle)
			return
		}
	}
}
func contains(n int, list []int) bool {
	for _, i := range list {
		if n == i {
			return true
		}
	}
	return false
}

func removeBellsUsed(c Cycle, bellMoves *[]BellInPermut) {
	newBellMoves := make([]BellInPermut, 0)
	for _, b := range *bellMoves {
		if !contains(b.bellId, c) {
			newBellMoves = append(newBellMoves, b)
		}
	}
	*bellMoves = newBellMoves
}

func NextCycle(bellMoves *[]BellInPermut) Cycle {
	fmt.Fprintf(os.Stderr, "Call next cycle  : %v\n", *bellMoves)
	moves := *bellMoves

	sort.Slice(moves, func(i, j int) bool {
		return moves[i].bellId < moves[j].bellId
	})

	if moves[0].bellId == moves[0].bellPosition {
		*bellMoves = moves[1:]
		fmt.Fprintf(os.Stderr, "Return empty cycle\n")
		return make(Cycle, 0)
	} else {
		c := make(Cycle, 1)
		c[0] = moves[0].bellId
		BuildCycle(moves[0].bellPosition, moves[0].bellId, *bellMoves, &c)
		removeBellsUsed(c, bellMoves)
		fmt.Fprintf(os.Stderr, "return %v\n", c)
		return c
	}
}

func BuildPrimaryPermut(bellMoves []BellInPermut) Permutation {
	// todo : sort bellMoves by bellId
	p := make(Permutation, 0)
	for len(bellMoves) > 0 {
		c := NextCycle(&bellMoves)
		if len(c) > 0 {
			p = append(p, c)
		}
	}
	return p
}

func (c Cycle) String() string {
	if len(c) == 0 {
		return ""
	}
	last := len(c) - 1
	s := ""
	for i, e := range c {
		s += fmt.Sprintf("%d", e)
		if i < last {
			s += " "
		}
	}
	return s
}
func (p Permutation) String() string {
	if len(p) == 0 {
		return "()"
	}
	s := ""
	for _, c := range p {
		s += fmt.Sprintf("(%v)", c)
	}
	return s
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	scanner.Scan()
	p := ParsePermutation(scanner.Text())
	fmt.Fprintf(os.Stderr, "Permut: %v\n", p)
	c := Compile(p)
	fmt.Fprintf(os.Stderr, "Bells moves: %v\n", c)
	pp := BuildPrimaryPermut(c)
	fmt.Fprintf(os.Stderr, "Final Permut: %v\n", pp)
	fmt.Printf("%v\n", pp) // Write answer to stdout
}
