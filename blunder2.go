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
func selectRoom(maxCashByRoom []int, visitedRooms map[int]bool) int {
	var maxCashRoomId int
	maxCashRoomId = 0
	for id, cash := range maxCashByRoom {
		if !visitedRooms[id] && cash > maxCashByRoom[maxCashRoomId] {
			maxCashRoomId = id
		}
	}
	return maxCashRoomId
}

func GetMaxCash(building Building) int {
	maxCashByRoom := make([]int, len(building))
	visitedRoom := make(map[int]bool, 0)
	var currentRoom int
	maxCashByRoom[currentRoom] = building[0].cash
	for len(visitedRoom) < len(building) {
		currentRoom = selectRoom(maxCashByRoom, visitedRoom) // Select current room (max cash value)
		nb1 := building[currentRoom].door1
		nb2 := building[currentRoom].door2
		if maxCashByRoom[nb1] < building[nb1].cash+maxCashByRoom[currentRoom] {
			maxCashByRoom[nb1] = building[nb1].cash + maxCashByRoom[currentRoom]
			delete(visitedRoom, nb1)
		}
		if maxCashByRoom[nb2] < building[nb2].cash+maxCashByRoom[currentRoom] {
			maxCashByRoom[nb2] = building[nb2].cash + maxCashByRoom[currentRoom]
			delete(visitedRoom, nb2)
		}
		visitedRoom[currentRoom] = true
	}
	return maxCashByRoom[len(building)-1]
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var N int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &N)
	building := make(Building, N+1)
	for i := 0; i < N; i++ {
		scanner.Scan()
		building[i] = NewRoom(scanner.Text(), N)
	}
	building[N] = NewRoom("E 0 E E", N)
	fmt.Fprintf(os.Stderr, "Building %v\n", building)
	fmt.Printf("%v\n", GetMaxCash(building)) // Write answer to stdout
}
