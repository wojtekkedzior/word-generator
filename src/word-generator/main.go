package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
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

func getPermutations(length int) int {
	var size = 0
	for i := (length); i > 0; i-- {
		var p = 1
		for j := length; j >= i; j-- {
			p = p * j
		}
		size = size + p
	}
	return size
}

func (topParent Node) lookup(str []byte) {
	rand.Seed(time.Now().UnixNano())
	strOri := str

	permutations := getPermutations(len(str))

	pos := make([][]int, 0) // this wil cause the array to be coppied immedtialy. perhaps we should use size / something?
	posWord := make(map[string]int)
	var count = 0

	start := time.Now()

	// work out all the possible permutations

	for len(pos) < permutations {
		randWordSize := rand.Intn(len(str) + 1)

		if randWordSize == 0 {
			continue
		}

		word := make([]byte, randWordSize)
		wordAsIndexs := make([]int, randWordSize)

		//store used indexes. they cannot be repeated
		usedIndex := make([]int, randWordSize)

		for i, _ := range usedIndex {
			usedIndex[i] = -1
		}

		for i, _ := range word {
			index := rand.Intn(len(str))

			for containsIndex(usedIndex, index) {
				index = rand.Intn(len(str))
			}

			usedIndex[i] = index
			word[i] = str[index]
			wordAsIndexs[i] = index
		}

		wordAsStr := byteArrToString(word)

		if posWord[wordAsStr] != 1 {
			posWord[wordAsStr] = 1
			pos = append(pos, wordAsIndexs)
		}

		str = strOri
		count++
	}

	fmt.Printf("Number of possibilities: %d \n", len(pos))
	fmt.Printf("Number of Random iterations: %d \n", count)
	fmt.Printf("Figuring out all the permutations took %s \n", time.Since(start))

	var foundWords sync.Map
	var wg sync.WaitGroup

	steps := len(pos) / 1000

	for i := 0; i <= steps; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			var end = 0
			//On the last 1000 th step
			if index == steps {
				end = len(pos)
			} else {
				end = (index + 1) * 1000
			}

			top := &topParent

			for _, v := range pos[(index * 1000):end] {
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
		}(i)
	}

	start = time.Now()
	wg.Wait()

	fmt.Printf("Traversing tree took %s \n", time.Since(start))

	m := map[string]interface{}{}
	foundWords.Range(func(key, value interface{}) bool {
		m[fmt.Sprint(key)] = value
		return true
	})

	fmt.Printf("Found a total of %d words.", len(m))
}

func main() {
	str := "planets" // 7
	//	str := "yoghurts" //8
	//	str := "youngster" //9
	var skippedDueToLength = 0
	var skippedDueToChar = 0

	strDict := make(map[rune]int)

	for _, v := range str {
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

		if len(r) > len(str) {
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

	topParent.lookup([]byte(str))
}
