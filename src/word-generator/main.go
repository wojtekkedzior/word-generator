package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
	"unicode/utf8"
)

func init() {
}

type Node struct {
	IsWord   bool
	Value    rune
	Childern map[rune]*Node
}

func containsIndex(s []int, index int) bool {
	for _, v := range s {
		if v == index {
			return true
		}
	}

	return false
}

func byteArrToString(byteArr []byte) string {
	var b bytes.Buffer
	for _, va := range byteArr {
		r, _ := utf8.DecodeRune([]byte{va})
		b.WriteRune(r)
	}

	return b.String()
}

func getNumberOfPermutations(length int) int {
	var size = 0
	for i := length; i > 0; i-- {
		var p = 1
		for j := length; j >= i; j-- {
			p = p * j
		}
		size = size + p
	}

	fmt.Printf("Number of possibilities for length of %d is %d \n", length, size)
	return size
}

func writePossibleWordArray(length int, indicies []int) {
	res := make([]int, 0, length) //set capacity to max word size

	for _, v := range indicies {
		if v != 0 {
			res = append(res, (v - 1))
		}
	}

	pos = append(pos, res)
}

func run(length int, index int) {
	if index == length {
		return
	}

	for k := 0; k < length; k++ {
		shouldContinue := true
		for j := 0; j < index; j++ {
			if k == counters[j] {
				shouldContinue = false
				break
			}
		}

		if shouldContinue {
			counters[index] = k
			writePossibleWordArray(length, counters[0:index+1])
			count++
			run(length, index+1)
		}
	}
}

var pos [][]int
var counters []int
var count int

func getPermutations(permCount int, str []byte) [][]int {
	realPermCount := getNumberOfPermutations(len(str) + 1)
	pos = make([][]int, 0, realPermCount)
	start := time.Now()
	length := len(str) + 1
	counters = make([]int, length)

	for i := 0; i < len(counters); i++ {
		counters[i] = 0
	}

	for i := 0; i < length; i++ {
		run(length, i)
	}

	fmt.Printf("Number of possibilites generated: %d \n", len(pos)) //this should match realPermCount.  The recursive call generetes extra, which is a buug
	fmt.Printf("Number of iterations to generete all permutations: %d \n", count)
	fmt.Printf("Time to generete all permutations %s \n", time.Since(start))

	return pos
}

func bruteforce(str []byte) [][]int {

	// rand.Seed(time.Now().UnixNano())
	// strOri := str
	realPermCount := getNumberOfPermutations(len(str) + 1)
	pos = make([][]int, 0, realPermCount)

	// posWord := make(map[string]int)
	// var count = 0

	// start := time.Now()

	// for len(pos) < permCount {
	// 	randWordSize := rand.Intn(len(str) + 1)

	// 	if randWordSize == 0 {
	// 		continue
	// 	}

	// 	word := make([]byte, randWordSize)
	// 	wordAsIndexs := make([]int, randWordSize)

	// 	//store used indexes. they cannot be repeated
	// 	usedIndex := make([]int, randWordSize)

	// 	for i, _ := range usedIndex {
	// 		usedIndex[i] = -1
	// 	}

	// 	for i, _ := range word {
	// 		index := rand.Intn(len(str))

	// 		for containsIndex(usedIndex, index) {
	// 			index = rand.Intn(len(str))
	// 		}

	// 		usedIndex[i] = index
	// 		word[i] = str[index]
	// 		wordAsIndexs[i] = index
	// 	}

	// 	wordAsStr := byteArrToString(word)

	// 	if posWord[wordAsStr] != 1 {
	// 		posWord[wordAsStr] = 1
	// 		pos = append(pos, wordAsIndexs)
	// 	}

	// 	str = strOri
	// 	count++
	// }

	// fmt.Printf("Number of possibilites generated: %d \n", len(pos))
	// fmt.Printf("Number of Random iterations: %d \n", count)
	// fmt.Printf("Figuring out all the permutations took %s \n", time.Since(start))

	return pos
}

func (topParent Node) lookup(str []byte) {
	//work out the number of all possible permutations
	permCount := getNumberOfPermutations(len(str))

	// work out all the possible permutations
	pos := getPermutations(permCount, str)

	//brute-force
	//func bruteforce

	var foundWords sync.Map
	var wg sync.WaitGroup

	start := time.Now()

	//each routine will process 1000 permutations
	for index := 0; index <= permCount/10000; index++ {
		wg.Add(1)

		var end = 0
		//On the last 1000 th step
		if index == permCount%10000 {
			end = permCount
		} else {
			end = (index + 1) * 10000
		}

		go func(start, end int) {
			defer wg.Done()
			top := &topParent

			for _, v := range pos[start:end] {
				var exist = false

				for i, vr := range v {
					r, _ := utf8.DecodeRune([]byte{str[vr]})
					n := top.Childern[r]

					if n == nil {
						exist = false
						break
					} else if i == len(v)-1 && n.IsWord { // is this a short word?
						exist = true
						break
					} else {
						top = top.Childern[r]
					}
				}

				if exist {
					var b bytes.Buffer

					for _, va := range v {
						r, _ := utf8.DecodeRune([]byte{str[va]})
						b.WriteRune(r)
					}

					foundWords.Store(b.String(), 1)
					exist = false
				}

				top = &topParent
			}
		}(index*10000, end)
	}

	wg.Wait()

	fmt.Printf("Traversing tree took %s \n", time.Since(start))

	//sync.Map doesn't have a way of revealing it's size, so have to convert it to a normal map or list
	var finalResult []string

	foundWords.Range(func(key, value interface{}) bool {
		finalResult = append(finalResult, fmt.Sprint(key))
		return true
	})

	fmt.Printf("Found a total of %d words.", len(finalResult))
}

func main() {
	// inputWord := "planets" // 7
	// inputWord := "timers" // 4
	//	inputWord := "yoghurts" //8
	inputWord := "youngster" //9
	var skippedDueToLength, skippedDueToChar = 0, 0

	strDict := make(map[rune]int)

	for _, v := range inputWord {
		strDict[v] = 1
	}

	f, err := os.Open("/usr/share/dict/british-english")

	if err != nil {
		fmt.Println(err)
	}

	parent := &Node{Childern: make(map[rune]*Node), IsWord: false, Value: ' '}
	topParent := parent

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		r := []rune(scanner.Text())

		if len(r) > len(inputWord) {
			skippedDueToLength++
			continue
		}

		cont := true

		for _, v := range r {
			if strDict[v] != 1 {
				skippedDueToChar++
				cont = false
				break
			}
		}

		if !cont {
			continue
		}

		for _, v := range r {
			if parent.Childern[v] == nil {
				node := &Node{Childern: make(map[rune]*Node), IsWord: false, Value: v}
				parent.Childern[v] = node
			}

			parent = parent.Childern[v]
		}
		parent.IsWord = true
		//start from the root for the new word
		parent = topParent
	}

	fmt.Printf("Skipped because of length: %d, Skipped because chars don't exist in provided word: %d.  Total skipped: %d \n", skippedDueToLength, skippedDueToChar, (skippedDueToChar + skippedDueToLength))

	topParent.lookup([]byte(inputWord))
}
