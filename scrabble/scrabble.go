package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/
func Permutations(nbElems int) [][]int {
	return buildPermuts(nbElems, buildNElems(nbElems), make([][]int, 0))
}

func GetAllPermuts(elems []int) [][]int {
	return buildPermuts(len(elems), elems, make([][]int, 0))
}

func buildNElems(n int) []int {
	elems := make([]int, n)
	for i := range elems {
		elems[i] = i
	}
	return elems
}

func buildPermuts(k int, elems []int, cumulatedPermuts [][]int) (permuts [][]int) {
	permuts = cumulatedPermuts
	if k == 1 {
		newPermut := make([]int, len(elems))
		copy(newPermut, elems)
		permuts = append(permuts, newPermut)
		return
	}
	permuts = buildPermuts(k-1, elems, permuts)
	for i := 0; i < k-1; i++ {
		if k%2 == 0 {
			elems[i], elems[k-1] = elems[k-1], elems[i]
		} else {
			elems[0], elems[k-1] = elems[k-1], elems[0]
		}
		permuts = buildPermuts(k-1, elems, permuts)
	}
	return
}

func rotate(elems []any) {
	e0 := elems[0]
	copy(elems, elems[1:])
	elems[len(elems)-1] = e0
}

func Combinations(k int, n int) [][]int {
	return buildCombinations(k, buildNElems(n), make([]int, 0), make([][]int, 0))
}

func buildCombinations(k int, elems []int, selectedElem []int, cumulatedCombis [][]int) (combi [][]int) {
	combi = cumulatedCombis
	if k > len(elems) {
		return cumulatedCombis
	} else if k == 0 {
		newCombi := make([]int, len(selectedElem))
		copy(newCombi, selectedElem)
		combi = append(combi, newCombi)
	} else {
		for i, e := range elems {
			remainingElems := elems[i+1:]
			newSelectedElems := make([]int, len(selectedElem)+1)
			copy(newSelectedElems, selectedElem)
			newSelectedElems[len(selectedElem)] = e
			combi = buildCombinations(k-1, remainingElems, newSelectedElems, combi)
		}
	}
	return
}

func eval(points map[rune]int, word string) int {
	score := 0
	for _, r := range word {
		score += points[r]
	}
	return score
}

func getPossibleWord(letters string) []string {
	combis := make([][]int, 0)
	nb := len(letters)
	for i := 0; i < nb; i++ {
		combis = append(combis, Combinations(i+1, nb)...)
	}
	permuts := make([][]int, 0)
	for _, c := range combis {
		permuts = append(permuts, GetAllPermuts(c)...)
	}
	list := make([]string, 0)
	for _, p := range permuts {
		word := make([]byte, 0)
		for _, i := range p {
			word = append(word, letters[i])
		}
		list = append(list, string(word))
	}
	return list
}

func scrabble() {
	var lettersPointsMap map[rune]int
	lettersPointsMap = map[rune]int{'a': 1, 'b': 3, 'c': 3, 'd': 2,
		'e': 1, 'f': 4, 'g': 2, 'h': 4, 'i': 1, 'j': 8, 'k': 5, 'l': 1,
		'm': 3, 'n': 1, 'o': 1, 'p': 3, 'q': 10, 'r': 1, 's': 1, 't': 1,
		'u': 1, 'v': 4, 'w': 4, 'x': 8, 'y': 4, 'z': 10}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	var N int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &N)
	dictionary := make([]string, N)
	for i := 0; i < N; i++ {
		scanner.Scan()
		W := scanner.Text()
		dictionary[i] = W
	}
	scanner.Scan()
	LETTERS := scanner.Text()
	_ = LETTERS
	combi := getPossibleWord(LETTERS)
	words := make(sort.StringSlice, len(combi))
	copy(words, combi)
	words.Sort()
	bestWord := ""
	bestScore := 0
	fmt.Fprintf(os.Stderr, "dictionary : %v\n", dictionary)
	fmt.Fprintf(os.Stderr, "words : %v\n", words)
	for _, word := range dictionary {
		// fmt.Fprintf(os.Stderr, "testing : %v\n", word)
		i := sort.SearchStrings(words, word)
		fmt.Fprintf(os.Stderr, "testing : %v, %d\n", word, i)
		if i < len(words) && words[i] == word {
			score := eval(lettersPointsMap, word)
			if score > bestScore {
				fmt.Fprintf(os.Stderr, "selecting : %v\n", word)
				bestWord = word
				bestScore = score
			}
		}
	}
	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Println(bestWord) // Write answer to stdout
}
