package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

var lengthPlusSpace, iterations, words = 0, 0, 0

var foundWords sync.Map
var wg sync.WaitGroup
var bufferSize = 200

var numOfResults = 10000

var client pulsar.Client
var err error

// Represents a char.  A link of Nodes can result in a word.
type Node struct {
	IsWord   bool
	Value    rune
	Childern map[rune]*Node
}

type Permutations struct {
	Permutations []perm `json:"permutations,omitempty"`
}

type perm struct {
	Permutation string `json:"permutation,omitempty"`
}

func init() {
	client, err = pulsar.NewClient(pulsar.ClientOptions{
		URL:                     "pulsar://192.168.122.10:4002",
		OperationTimeout:        30 * time.Second,
		ConnectionTimeout:       30 * time.Second,
		MaxConnectionsPerBroker: 10,
	})

	if err != nil {
		fmt.Println("Failed to create client", err)
	}

	defer client.Close()
}

func getNumberOfPermutations(length int) {
	var size = 0
	for i := length; i > 0; i-- {
		var p = 1
		for j := length; j >= i; j-- {
			p = p * j
		}
		size = size + p
	}

	fmt.Printf("Number of possibilities for length of %d is %d \n", length, size)
}

// index   -
// counter -
// results - a channel where a message container multiple permutations will be written to.
// buffer  - contains each permutation as a comma delimited string of integers
// permutations - holds a list of Permutations. The array is marshaled into JSON when the array size exceeds the "bufferSize"
// permutation  - an array of indexes representing a single permutation. This array is cleared after the generated permutation is passed on.
func run(index int, counters []int, results chan string, buffer *bytes.Buffer, permutations *Permutations, permutation []int) {
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

			for _, v := range counters[0 : index+1] {
				if v != 0 { //effectively trimming
					permutation = append(permutation, (v - 1))
				}
			}

			for _, i := range permutation {
				buffer.WriteString(strconv.Itoa(i))
				buffer.WriteString(",")
			}

			permutations.Permutations = append(permutations.Permutations, perm{Permutation: buffer.String()})
			buffer.Reset()
			permutation = permutation[:0]

			if len(permutations.Permutations) > bufferSize {
				jsonData, err := json.Marshal(permutations)
				if err != nil {
					log.Fatal(err)
				}

				results <- string(jsonData)

				permutations.Permutations = permutations.Permutations[:0]
			}

			run(index+1, counters, results, buffer, permutations, permutation)
		}
	}
}

func createChannelReader(id int) chan string {
	var results = make(chan string, numOfResults)
	var producer pulsar.Producer

	go func(id int) {
		producer, err = client.CreateProducer(pulsar.ProducerOptions{
			Topic:                   "t/ns/mercury",
			BatchingMaxPublishDelay: 10 * time.Millisecond, // Maximum time to wait for batching
			BatchingMaxMessages:     uint(500),             // Maximum number of messages in a batch
			BatchingMaxSize:         1024 * 1024 * 10,      // Maximum size of batch (10MB)
			Name:                    fmt.Sprintf("producer-for-channel-%v", +id),
		})

		if err != nil {
			fmt.Println("Failed to create producer: ", err)
		}

		defer producer.Close()

		for msg := range results {
			producer.SendAsync(context.Background(), &pulsar.ProducerMessage{
				Payload: []byte(msg),
			}, func(msgID pulsar.MessageID, msg *pulsar.ProducerMessage, err error) {
				if err != nil {
					fmt.Printf("Failed to send message: %v", err)
					return
				}
			})
		}

		close(results)
	}(id)

	return results
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
			results := createChannelReader(index)

			// var permutations Permutations
			permutations := &Permutations{}

			permutation := make([]int, 0, lengthPlusSpace)

			buffer := bytes.NewBuffer(make([]byte, 0, bufferSize))
			run(index, counters, results, buffer, permutations, permutation)
		}(i)
	}

	wg.Wait()

	// for permutation := range results {
	// 	iterations++

	// 	top := &topParent

	// 	for i, vr := range permutation {
	// 		r, _ := utf8.DecodeRune([]byte{str[vr]})
	// 		n := top.Childern[r]

	// 		if n == nil {
	// 			break
	// 		} else if i == len(permutation)-1 && n.IsWord { // is this a short word?
	// 			var b bytes.Buffer

	// 			for _, va := range permutation {
	// 				r, _ := utf8.DecodeRune([]byte{str[va]})
	// 				b.WriteRune(r)
	// 			}

	// 			foundWords.Store(b.String(), 1)
	// 			words++
	// 			break
	// 		} else {
	// 			top = top.Childern[r]
	// 		}
	// 	}

	// 	top = &topParent
	// }

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
	inputWord := "proselytize" //11
	// inputWord := "abandonwares" //12
	// inputWord := "ventriloquizes" //14
	// inputWord := "kaiserdoms" //10
	// inputWord := "Counterrevolutionary" //20 - 6613313319248080000 possibilites :D
	// inputWord := "planets" //7
	// inputWord := "youngster" //9
	// inputWord := "or" //9
	// inputWord := "helper"

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
