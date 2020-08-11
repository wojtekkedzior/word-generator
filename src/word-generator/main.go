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
	//	"github.com/google/uuid"
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
	//	fmt.Println(s)
	for _, v := range s {
		if v == index {
			return true
		}
	}

	return false
}

func (topParent *Node) lookup() {
	size := 13699
	rand.Seed(42)
	pos := make([][]int, size)
	top := topParent
	lastIndex := 0
	str := "planets"
	strOri := str

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

		if containsUUID(pos, wordAsIndexs) {
			fmt.Println("do nothing for: ", wordAsIndexs)
		} else {
			pos[lastIndex] = wordAsIndexs
			lastIndex++
		}

		str = strOri
	}

	fmt.Println(pos)

	foundsWords := make(map[string]int)

	for _, v := range pos {
		var exist = false

		var b bytes.Buffer

		for _, va := range v {
			r, _ := utf8.DecodeRune([]byte{str[va]})
			b.WriteRune(r)
		}

		for i, vr := range v {
			r, _ := utf8.DecodeRune([]byte{str[vr]})
			n := topParent.Childern[r]
			if n == nil {
				exist = false
				break
			} else if i == len(v)-1 && n.IsWord { // is this a short word?
				exist = true
				break
			} else {
				topParent = topParent.Childern[r]
			}
		}
		if exist {
			var b bytes.Buffer

			for _, va := range v {
				r, _ := utf8.DecodeRune([]byte{str[va]})
				b.WriteRune(r)
			}

			foundsWords[b.String()] = 1

			fmt.Println(b.String(), " - ------------------word exisit")
			exist = false
		}
		topParent = top
	}

	for i, _ := range foundsWords {
		fmt.Println(i)
	}

	fmt.Println(len(foundsWords))
}

func main() {
	f, err := os.Open("/home/wojtek/workspace/word-generator/bin/words_alpha.txt")
	//	f, err := os.Open("/home/wojtek/workspace/word-generator/bin/smallWords.txt")

	if err != nil {
		fmt.Println(err)
	}

	parent := &Node{Childern: make(map[rune]*Node), IsWord: false, Value: ' '}
	topParent := parent

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		r := []rune(scanner.Text())
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

	topParent.lookup()
}
