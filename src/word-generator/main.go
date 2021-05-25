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

// Represents a char.  A link of Nodes can result in a word.
type Node struct {
	IsWord   bool
	Value    rune
	Childern map[rune]*Node
}

// Holds all of the permutations
var permutations [][]int

/*
 One counter for each char in the word.  Each counter is max size of the word + 1 (an empty space).
 The counters can be thought of as number wheel on a padlock:

   1 3 2
   2 4 3
 [ 3 5 4 ]  <- we are only looking at the alignment of numbers here.  We spin each wheel one digit at a time and take the value
   4 6 5
   5 7 6

  The permutations are actually indexes of the array that holds the word.
  This way we can convert each permutation into char quickly.
*/
var counters []int

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

func savePermutations(length int, indicies []int) {
	res := make([]int, 0, length)

	for _, v := range indicies {
		if v != 0 {
			res = append(res, (v - 1))
		}
	}

	permutations = append(permutations, res)
}

func run(length, index, limit, count int) {
	if index == length || count == limit {
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
			savePermutations(length, counters[0:index+1])
			count++
			run(length, index+1, limit, count)
		}
	}
}

func getPermutations(permCount int, str []byte) {
	realPermCount := getNumberOfPermutations(len(str) + 1)
	permutations = make([][]int, 0, realPermCount)
	length := len(str) + 1
	counters = make([]int, length)
	start := time.Now()

	for i := 0; i < len(counters); i++ {
		counters[i] = 0
	}

	for i := 0; i < length; i++ {
		run(length, i, realPermCount, 0)
	}

	fmt.Printf("Number of possibilites generated: %d \n", len(permutations))
	fmt.Printf("Time to generete all permutations %s \n", time.Since(start))
}

func (topParent Node) lookup(str []byte) {
	//work out the number of all possible permutations
	permCount := getNumberOfPermutations(len(str))

	// work out all the possible permutations
	getPermutations(permCount, str)

	var foundWords sync.Map
	var wg sync.WaitGroup

	start := time.Now()

	segmentSize := 10000

	//each routine will process 1000 permutations
	for index := 0; index <= permCount/segmentSize; index++ {
		wg.Add(1)

		var end = 0
		//On the last 1000 th step
		if index == permCount%segmentSize {
			end = permCount
		} else {
			end = (index + 1) * segmentSize
		}

		go func(start, end int) {
			defer wg.Done()
			top := &topParent

			for _, v := range permutations[start:end] {
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
		}(index*segmentSize, end)
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
	// inputWord := "dogs" // 4
	//	inputWord := "yoghurts" //8
	inputWord := "youngster" //9
	// inputWord := "abcdefghij" //9
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
