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

	fmt.Printf("Number of possibilities: %d \n", size)
	return size
}

func writePossibleWord(pos [][]int, indicies ...int) [][]int {
	res := make([]int, 0, 10) //set capacity to max word size

	for _, v := range indicies {
		if v != 0 {
			res = append(res, (v - 1))
		}
	}

	return append(pos, res)
}

func getPermutations(permCount int, str []byte) [][]int {
	pos := make([][]int, 0, 9864100) // this will cause the array to be copied immedtialy. perhaps we should use size / something?
	var count = 0

	start := time.Now()

	count, c1, c2, c3, c4, c5, c6, c7, c8, c9, c10 := 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0

	for i := 0; i < 10; i++ {
		c1 = i
		count++
		pos = writePossibleWord(pos, c1)
		for j := 0; j < 10; j++ {
			if j == i {
				continue
			}
			c2 = j
			count++
			pos = writePossibleWord(pos, c1, c2)
			for k := 0; k < 10; k++ {
				if k == j || k == i {
					continue
				}
				c3 = k
				count++
				pos = writePossibleWord(pos, c1, c2, c3)
				for l := 0; l < 10; l++ {
					if l == k || l == j || l == i {
						continue
					}
					c4 = l
					count++
					pos = writePossibleWord(pos, c1, c2, c3, c4)
					for m := 0; m < 10; m++ {
						if m == l || m == k || m == j || m == i {
							continue
						}
						c5 = m
						count++
						pos = writePossibleWord(pos, c1, c2, c3, c4, c5)
						for n := 0; n < 10; n++ {
							if n == m || n == l || n == k || n == j || n == i {
								continue
							}
							c6 = n
							count++
							pos = writePossibleWord(pos, c1, c2, c3, c4, c5, c6)
							for o := 0; o < 10; o++ {
								if o == n || o == m || o == l || o == k || o == j || o == i {
									continue
								}
								c7 = o
								count++
								pos = writePossibleWord(pos, c1, c2, c3, c4, c5, c6, c7)

								for p := 0; p < 10; p++ {
									if p == o || p == n || p == m || p == l || p == k || p == j || p == i {
										continue
									}
									c8 = p
									count++
									pos = writePossibleWord(pos, c1, c2, c3, c4, c5, c6, c7, c8)

									for r := 0; r < 10; r++ {
										if r == p || r == o || r == n || r == m || r == l || r == k || r == j || r == i {
											continue
										}
										c9 = r
										count++
										pos = writePossibleWord(pos, c1, c2, c3, c4, c5, c6, c7, c8, c9)

										for s := 0; s < 10; s++ {
											if s == r || s == p || s == o || s == n || s == m || s == l || s == k || s == j || s == i {
												continue
											}
											c10 = s
											count++
											pos = writePossibleWord(pos, c1, c2, c3, c4, c5, c6, c7, c8, c9, c10)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// rand.Seed(time.Now().UnixNano())
	// strOri := str
	// pos := make([][]int, 0) // this will cause the array to be coppied immedtialy. perhaps we should use size / something?
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

	fmt.Printf("Number of possibilites generated: %d \n", len(pos))
	fmt.Printf("Number of Random iterations: %d \n", count)
	fmt.Printf("Figuring out all the permutations took %s \n", time.Since(start))

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

// for 9 letters
// wojtek@wojtek-pc:~/git/word-generator/src/word-generator$ go run main.go
// Skipped because of length: 33146, Skipped because chars don't exist in provided word: 68028.  Total skipped: 101174
// Number of possibilities: 986409
// Number of Random iterations: 59485463
// Figuring out all the permutations took 23.754992341s
// Traversing tree took 74.642873ms
// Found a total of 264 words.
