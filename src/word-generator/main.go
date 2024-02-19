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

// var segmentSize = 10000
var lengthPlusSpace, iterations, words = 0, 0, 0

var foundWords sync.Map
var wg sync.WaitGroup
var results = make(chan []int, 10000)

// Represents a char.  A link of Nodes can result in a word.
type Node struct {
	IsWord   bool
	Value    rune
	Childern map[rune]*Node
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

func run(index int, counters []int, c chan []int) {
	if index == lengthPlusSpace {
		return
	}

	for j := 0; j < lengthPlusSpace; j++ {
		shouldContinue := true //skipping 'self'
		for k := 0; k < index; k++ {
			if j == counters[k] {
				shouldContinue = false
				break
			}
		}
		if shouldContinue {
			counters[index] = j

			permutation := make([]int, 0, lengthPlusSpace)
			for _, v := range counters[0 : index+1] {
				if v != 0 { //effectively trimming
					permutation = append(permutation, (v - 1))
				}
			}

			c <- permutation
			run(index+1, counters, c)
		}
	}
}

func (topParent Node) lookup(str []byte) {
	//+1 for the extra space char
	lengthPlusSpace = len(str) + 1

	//work out the number of all possible permutations
	getNumberOfPermutations(len(str))

	// work out all the possible permutations
	getNumberOfPermutations(lengthPlusSpace)

	start := time.Now()

	for i := 0; i < lengthPlusSpace; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			counters := make([]int, lengthPlusSpace)
			run(index, counters, results)
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for permutation := range results {
		iterations++

		top := &topParent

		for i, vr := range permutation {
			r, _ := utf8.DecodeRune([]byte{str[vr]})
			n := top.Childern[r]

			if n == nil {
				break
			} else if i == len(permutation)-1 && n.IsWord { // is this a short word?
				var b bytes.Buffer

				for _, va := range permutation {
					r, _ := utf8.DecodeRune([]byte{str[va]})
					b.WriteRune(r)
				}

				foundWords.Store(b.String(), 1)
				words++
				break
			} else {
				top = top.Childern[r]
			}
		}

		top = &topParent
	}

	var finalResult []string

	foundWords.Range(func(key, value interface{}) bool {
		finalResult = append(finalResult, fmt.Sprint(key))
		return true
	})

	fmt.Printf("Found a total of %d words. \n", len(finalResult))
	fmt.Println(iterations)
	fmt.Println(words)
	// fmt.Printf("Number of possibilites generated: %d \n", len(permutations))
	fmt.Printf("Time to generete all permutations %s \n", time.Since(start))
}

func main() {
	// inputWord := "proselytize" //11
	// inputWord := "abandonwares" //12
	// inputWord := "ventriloquizes" //14
	inputWord := "kaiserdoms" //10
	// inputWord := "Counterrevolutionary" //20 - 6613313319248080000 possibilites :D
	// inputWord := "planets" //7
	// inputWord := "youngster" //9
	// inputWord := "or" //9

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

// func test(c chan int) int {
// 	mycount := 0
// 	for i := 0; i < 1000000; i++ {
// 		// count++
// 		c <- i
// 		mycount++
// 	}
// 	return mycount
// }

// for complete != 8 {
// 	select {
// 	case <-c:
// 		// fmt.Println(complete)
// 	}
// }
