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
func selectRoom(maxCashByRoom []int, visitedRooms Set) int {
	var maxCashRoomId int
	maxCashRoomId = 0
	for id, cash := range maxCashByRoom {
		if _, ok := visitedRooms[id]; !ok && cash > maxCashByRoom[maxCashRoomId] {
			maxCashRoomId = id
		}
	}
	return maxCashRoomId
}

func GetMaxCash(building Building, graph []Set) int {
	N := len(building) - 1
	maxCashByRoom := make([]int, len(building))
	visitedRoom := make(Set, 0)
	currentRoom := N
	for room := range graph[N] {
		maxCashByRoom[room] = building[room].cash
	}
	for currentRoom > 0 {
		currentRoom = selectRoom(maxCashByRoom, visitedRoom) // Select current room (max cash value)
		for room := range graph[currentRoom] {
			if maxCashByRoom[room] < building[room].cash+maxCashByRoom[currentRoom] {
				maxCashByRoom[room] = building[room].cash + maxCashByRoom[currentRoom]
			}
		}
		visitedRoom[currentRoom] = member
	}
	return maxCashByRoom[0]
}

func main() {
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
