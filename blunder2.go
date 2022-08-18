package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 * Implementation : expecting room present in seq order from 0 to N
 **/

type Room struct {
	cash  int
	door1 int
	door2 int
}
type void struct{}

type Set map[int]void

var member void

const EXIT = -1

type Building []Room

func NewRoom(room string, idExit int) Room {
	roomParams := strings.Split(room, " ")
	var newRoom Room
	if cash, err := strconv.Atoi(roomParams[1]); err != nil {
		panic(err)
	} else {
		newRoom = Room{cash, 0, 0}
	}
	if door1, err := strconv.Atoi(roomParams[2]); err != nil {
		newRoom.door1 = idExit
	} else {
		newRoom.door1 = door1
	}
	if door2, err := strconv.Atoi(roomParams[3]); err != nil {
		newRoom.door2 = idExit
	} else {
		newRoom.door2 = door2
	}
	return newRoom
}
func selectRoom(visitedRooms Set, parentsGraph []Set) int {
	for i := 0; i < len(parentsGraph); i++ {
		if _, v := visitedRooms[i]; !v {
			selectable := true
			for pId := range parentsGraph[i] {
				if _, v := visitedRooms[pId]; !v {
					selectable = false
					break
				}
			}
			if selectable {
				return i
			}
		}
	}
	return len(parentsGraph)
}

func GetMaxRoomCash(maxCashByRoom []int, parents Set) int {
	var maxCash int
	for pId := range parents {
		if maxCash < maxCashByRoom[pId] {
			maxCash = maxCashByRoom[pId]
		}
	}
	return maxCash
}

func GetMaxCash(building Building, parentsGraph []Set) int {
	maxCashByRoom := make([]int, len(building))
	maxCashByRoom[0] = building[0].cash
	visitedRoom := make(Set, 0)
	visitedRoom[0] = member
	var currentRoom int
	for {
		currentRoom = selectRoom(visitedRoom, parentsGraph)
		if currentRoom == len(parentsGraph) {
			break
		}
		if maxCash := GetMaxRoomCash(maxCashByRoom, parentsGraph[currentRoom]); maxCash > 0 {
			maxCashByRoom[currentRoom] = maxCash + building[currentRoom].cash
		}
		visitedRoom[currentRoom] = member
	}
	return maxCashByRoom[len(building)-1]
}

func blunder2() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var N int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &N)
	building := make(Building, N+1)
	neighboors := make([]Set, N+1)
	for i := 0; i < N+1; i++ {
		neighboors[i] = make(Set, 0)
	}

	for i := 0; i < N; i++ {
		scanner.Scan()
		room := NewRoom(scanner.Text(), N)
		neighboors[room.door1][i] = member
		neighboors[room.door2][i] = member
		building[i] = room
	}
	building[N] = NewRoom("E 0 E E", N)
	fmt.Fprintf(os.Stderr, "Building %v\n", building)
	fmt.Printf("%v\n", GetMaxCash(building, neighboors)) // Write answer to stdout
}
