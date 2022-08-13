package main

import (
	"bufio"
	"fmt"
	"os"
)

type Blunder struct {
	face            byte
	walking         byte
	breaker         bool
	inversePriority bool
}

type Coord struct {
	row int
	col int
}

type Step struct {
	blunder Blunder
	pos     Coord
}

type Transition func(Step) Step

type FSM map[byte]Transition

type Path []Step

type Map [][]byte

const SOUTH byte = 0
const EAST byte = 1
const NORTH byte = 2
const WEST byte = 3
const LOOP byte = 4
const STAY byte = 5

var PRIORITY []byte
var MESSAGES []string
var BlenderPathSimu FSM
var LoopStep Step
var invalidCoord Coord

func init() {
	MESSAGES = []string{SOUTH: "SOUTH", EAST: "EAST", NORTH: "NORTH", WEST: "WEST", LOOP: "LOOP", STAY: "STAY"}
	BlenderPathSimu = FSM{'#': Turn, 'X': Break, '$': StepForward,
		'S': StepAndSouth, 'E': StepAndEast, 'N': StepAndNort, 'W': StepAndWest,
		'B': Drink, 'I': Reverse, 'T': StepForward, ' ': StepForward}
	LoopStep = Step{Blunder{STAY, LOOP, false, false}, Coord{0, 0}}
	invalidCoord = Coord{-1, -1}
}

// #, X, @, $, S, E, N, W, B, I, T
// Default priority S E N W

func (c Coord) NextCoord(face byte) Coord {
	switch face {
	case SOUTH:
		return Coord{row: c.row + 1, col: c.col}
	case EAST:
		return Coord{row: c.row, col: c.col + 1}
	case NORTH:
		return Coord{row: c.row - 1, col: c.col}
	case WEST:
		return Coord{row: c.row, col: c.col - 1}
	}
	panic(face)
}

func Turn(lastStep Step) (newStep Step) {
	newStep = Step{blunder: lastStep.blunder, pos: lastStep.pos}
	newStep.blunder.walking = STAY
	if lastStep.blunder.walking == STAY {
		if newStep.blunder.inversePriority {
			newStep.blunder.face = ((newStep.blunder.face + 3) % 4)
		} else {
			newStep.blunder.face = ((newStep.blunder.face + 1) % 4)
		}
	} else {
		if newStep.blunder.inversePriority {
			newStep.blunder.face = WEST
		} else {
			newStep.blunder.face = SOUTH
		}
	}
	return
}

func Break(lastStep Step) (newStep Step) {
	if lastStep.blunder.breaker {
		return StepForward(lastStep)
	}
	return Turn(lastStep)
}

func Drink(lastStep Step) Step {
	newStep := Step{lastStep.blunder, lastStep.pos}
	newStep.blunder.breaker = !newStep.blunder.breaker
	return StepForward(newStep)
}

func Reverse(lastStep Step) Step {
	newStep := Step{lastStep.blunder, lastStep.pos}
	newStep.blunder.inversePriority = !newStep.blunder.inversePriority
	return StepForward(newStep)
}

func StepAndSouth(lastStep Step) Step {
	return StepAndTurn(lastStep, SOUTH)
}

func StepAndEast(lastStep Step) Step {
	return StepAndTurn(lastStep, EAST)
}

func StepAndNort(lastStep Step) Step {
	return StepAndTurn(lastStep, NORTH)
}

func StepAndWest(lastStep Step) Step {
	return StepAndTurn(lastStep, WEST)
}

func StepAndTurn(lastStep Step, turn byte) Step {
	newStep := StepForward(lastStep)
	newStep.blunder.face = turn
	return newStep
}

func StepForward(lastStep Step) (newStep Step) {
	newBlunder := lastStep.blunder
	newBlunder.walking = lastStep.blunder.face
	return Step{blunder: newBlunder, pos: lastStep.pos.NextCoord(lastStep.blunder.face)}
}

func Simulate(game Map) Path {
	var startPos, lastPos Coord
	portails := make([]Coord, 2)
	portails[0] = invalidCoord
	for row, line := range game {
		for col, zone := range line {
			switch zone {
			case '@':
				startPos = Coord{row, col}
			case '$':
				lastPos = Coord{row, col}
			case 'T':
				if portails[0] == invalidCoord {
					portails[0] = Coord{row, col}
				} else {
					portails[1] = Coord{row, col}
				}
			default:
			}
		}
	}
	game[startPos.row][startPos.col] = ' '
	step := Step{Blunder{SOUTH, STAY, false, false}, startPos}
	loop := make(Path, 0)
	path := make(Path, 0)
	for {

		nextCoord := step.pos.NextCoord(step.blunder.face)
		nextZone := game[nextCoord.row][nextCoord.col]
		step = BlenderPathSimu[nextZone](step)
		// if current pos is X replace by space and clear loop slice
		if game[step.pos.row][step.pos.col] == 'X' {
			game[step.pos.row][step.pos.col] = ' '
			loop = make(Path, 0)
		}
		// if current pos is T replace current pos by teleported position
		if game[step.pos.row][step.pos.col] == 'T' {
			if step.pos == portails[0] {
				step.pos = portails[1]
			} else {
				step.pos = portails[0]
			}
		}
		// if current pos in loop slice return "LOOP"
		if CheckLoop(step, loop) {
			return Path{0: LoopStep}
		}
		// add current step to path and loop
		loop = append(loop, step)
		path = append(path, step)
		// if current pos is $ then return path
		if step.pos == lastPos {
			return path
		}
	}
}

func CheckLoop(step Step, previousSteps Path) bool {
	for _, ps := range previousSteps {
		if ps == step {
			return true
		}
	}
	return false
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var L, C int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &L, &C)
	fmt.Fprintf(os.Stderr, "%v %v\n", L, C)
	game := make(Map, L)
	for i := 0; i < L; i++ {
		scanner.Scan()
		row := scanner.Text()
		fmt.Fprintf(os.Stderr, "%v\n", row)
		game[i] = []byte(row)
	}
	path := Simulate(game)
	for _, s := range path {
		fmt.Fprintf(os.Stderr, "Step: %v, (%v)%v\n", MESSAGES[s.blunder.face], MESSAGES[s.blunder.walking], s.pos)
	}
	for _, step := range path {
		if step.blunder.walking != STAY {
			fmt.Println(MESSAGES[step.blunder.walking])
		}
	}
	// fmt.Fprintln(os.Stderr, "Debug messages...")
}
