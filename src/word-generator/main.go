package main

import (
	"bufio"
	"fmt"
	//	"io/ioutil"
	"math/rand"
	"os"
	//	"strings"
	//	"github.com/davecgh/go-spew/spew"
	//	"errors"
	"bytes"
	//	"sync"
	//	"github.com/google/uuid"
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

func containsUUID(s [][]int, word []int) bool {
	for _, v := range s {
		hits := 0

		if len(v) != len(word) {
			continue
		}

		for j, k := range v {
			if k != word[j] {
				break
			} else {
				hits++
			}
		}

		if hits == len(word) {
			return true
		}
	}

	return false
}

func containsIndex(s []int, index int) bool {
	for _, v := range s {
		if v == index {
			return true
		}
	}

	return false
}

func (topParent Node) lookup(str []byte) {
	rand.Seed(time.Now().UnixNano())
	//	rand.Seed(42)
	//	top := topParent
	lastIndex := 0
	strOri := str
	var size = 0
	var p = 1

	for i := (len(str)); i > 0; i-- {
		for j := len(str); j >= i; j-- {
			p = p * j
		}
		size = size + p
		p = 1
	}

	fmt.Println("size; ", size)
	pos := make([][]int, size)

	var count = 0

	start := time.Now()

	for lastIndex < size {
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

		if !containsUUID(pos, wordAsIndexs) {
			//			fmt.Println("do nothing for: ", wordAsIndexs)
			//		} else {
			pos[lastIndex] = wordAsIndexs
			lastIndex++
		}

		str = strOri

		count++
	}
	fmt.Println(len(pos))
	fmt.Printf("Randon gen count: %d \n", count)
	//	os.Exit(1)

	elapsed := time.Since(start)
	fmt.Printf("Figuring out all the permutatiosn took %s \n", elapsed)

	foundsWords := make(map[string]int)

	steps := len(pos) / 1000
	fmt.Println(steps)
	c := make(chan string)

	fmt.Println("address of topParent: ", &topParent)
	fmt.Println("value of topParent: ", topParent)
	//	fmt.Println("value of topParent: ", *topParent)

	for i := 0; i <= steps; i++ {
		go func(index int, co chan<- string) {
			var end = 0
			//On the last 1000 th step
			if index == steps {
				end = len(pos)
			} else {
				end = (index + 1) * 1000
			}

			var top *Node
			top = &topParent
			//			fmt.Println("address of top2: ", top2)

			//			top := topParent
			fmt.Println("address of top: ", top)
			//			fmt.Println("address of top: ", top)
			//			fmt.Println("Value  top points to: ", *top)

			for _, v := range pos[(index * 1000):end] {
				var exist = false

				for i, vr := range v {
					r, _ := utf8.DecodeRune([]byte{str[vr]})
					n := top.Childern[r]
					fmt.Println(r, n, i)

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

				fmt.Println(exist)

				if exist {
					var b bytes.Buffer

					for _, va := range v {
						r, _ := utf8.DecodeRune([]byte{str[va]})
						b.WriteRune(r)
					}

					co <- fmt.Sprintf(b.String())
					exist = false
				}

				top = &topParent

				//				var b bytes.Buffer
				//				for _, va := range v {
				//					r, _ := utf8.DecodeRune([]byte{str[va]})
				//					b.WriteRune(r)
				//				}
				//				c <- fmt.Sprintf("checking: %s, %d, ", b.String(), index)
			}
		}(i, c)

	}

	start = time.Now()

	for i := 0; i < 159; i++ {
		//		fmt.Println(i)
		foundsWords[<-c] = 1
	}

	elapsed = time.Since(start)
	fmt.Printf("Traversing tree took %s \n", elapsed)

	//	var wg sync.WaitGroup
	//
	//		wg.Add(steps)
	//	for i := 0; i < steps; i++ {
	//
	//		go func() {
	//			defer wg.Done()
	//
	//			//			for _, v := range pos[i:(i * 1000)] {
	//			//				fmt.Println("adad", v)
	//			//			}
	//
	//			for _, v := range pos[(i * 1000) : (i+1)*1000] {
	//				var exist = false
	//
	//				for i, vr := range v {
	//					r, _ := utf8.DecodeRune([]byte{str[vr]})
	//					n := topParent.Childern[r]
	//					if n == nil {
	//						exist = false
	//						break
	//					} else if i == len(v)-1 && n.IsWord { // is this a short word?
	//						exist = true
	//						break
	//					} else {
	//						topParent = topParent.Childern[r]
	//					}
	//				}
	//				if exist {
	//					var b bytes.Buffer
	//
	//					for _, va := range v {
	//						r, _ := utf8.DecodeRune([]byte{str[va]})
	//						b.WriteRune(r)
	//					}
	//
	//					foundsWords[b.String()] = 1
	//					exist = false
	//				}
	//				topParent = top
	//			}
	//
	//		}()
	//
	//	}

	//Works without mutex

	//	for _, v := range pos {
	//		var exist = false
	//
	//		for i, vr := range v {
	//			r, _ := utf8.DecodeRune([]byte{str[vr]})
	//			n := topParent.Childern[r]
	//			if n == nil {
	//				exist = false
	//				break
	//			} else if i == len(v)-1 && n.IsWord { // is this a short word?
	//				exist = true
	//				break
	//			} else {
	//				topParent = topParent.Childern[r]
	//			}
	//		}
	//		if exist {
	//			var b bytes.Buffer
	//
	//			for _, va := range v {
	//				r, _ := utf8.DecodeRune([]byte{str[va]})
	//				b.WriteRune(r)
	//			}
	//
	//			foundsWords[b.String()] = 1
	//			exist = false
	//		}
	//		topParent = top
	//	}

	fmt.Println(len(foundsWords))
}

func main() {
	str := "planets"
	var skippedDueToLength = 0
	var skippedDueToChar = 0

	strDict := make(map[rune]int)

	for _, v := range str {
		strDict[v] = 1
	}

	//	f, err := os.Open("/home/wojtek/workspace/word-generator/bin/words_alpha.txt")
	//	f, err := os.Open("/home/wojtek/workspace/word-generator/bin/smallWords.txt")
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
