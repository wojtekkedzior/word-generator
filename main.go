package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

// TODO - actualPermutations is always more than what is expected. For example for the word 'planets', it being a word of 7 characters, we get a maximum permutation count of 109600 (7+1 for the space). But actual count is 135411, which is wrong.
// this can be seen with the inputword of 3 letters, such as 'pot' Some permutation appear three times.
var lengthPlusSpace, actualPermutations = 0, 0

var wg sync.WaitGroup
var bufferSize = 200

var numOfResults = 10000

var client pulsar.Client
var err error

// Represents a char.  A link of Nodes can result in a word.
// type Node struct {
// 	IsWord   bool
// 	Value    rune
// 	Childern map[rune]*Node
// }

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
func run(index int, counters []int, results chan string, buffer *bytes.Buffer, permutations *Permutations) {
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
					buffer.WriteString(strconv.Itoa((v - 1)))
					buffer.WriteString(",")
				}
			}
			actualPermutations++
			permutations.Permutations = append(permutations.Permutations, perm{Permutation: buffer.String()})
			buffer.Reset()

			if len(permutations.Permutations) > bufferSize {
				jsonData, err := json.Marshal(permutations)
				if err != nil {
					log.Fatal(err)
				}

				results <- string(jsonData)

				permutations.Permutations = permutations.Permutations[:0]
			}

			run(index+1, counters, results, buffer, permutations)
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
			BatchingMaxMessages:     uint(600),             // Maximum number of messages in a batch
			BatchingMaxSize:         1024 * 1024 * 5,       // Maximum size of batch (10MB)
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

func lookup(str []byte) {
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

			buffer := bytes.NewBuffer(make([]byte, 0, bufferSize))
			run(index, counters, results, buffer, permutations)
		}(i)
	}

	wg.Wait()

	fmt.Printf("actual number of permutations: %d \n", actualPermutations)
	fmt.Printf("Time to generate all permutations %s \n", time.Since(start))
}

func main() {
	// all times include producing messages to Pulsar using 1 channel per character (include an extra character for a space)

	//7 - Time to generate all permutations 32.523004ms
	inputWord := "planets"

	//9 - Time to generate all permutations 2.449610907s
	// inputWord := "youngster"

	//10 - Time to generete all permutations 28.532050481s
	// inputWord := "kaiserdoms"

	//11 - Time to generate all permutations 6m17.08988457s
	// inputWord := "proselytize"

	//12 - ~120 minutes
	// inputWord := "abandonwares"

	//14 - n/a
	// inputWord := "ventriloquizes"

	//20 - 6613313319248080000 possibilites :D
	// inputWord := "Counterrevolutionary"

	//---------------------
	// TODO - this needs to be moved into java.
	//---------------------

	// var skippedDueToLength, skippedDueToChar = 0, 0

	// strDict := make(map[rune]int)
	// for _, v := range inputWord {
	// 	strDict[v] = 1
	// }

	// f, err := os.Open("/usr/share/dict/british-english")

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// parent := &Node{Childern: make(map[rune]*Node), IsWord: false, Value: ' '}
	// topParent := parent

	// scanner := bufio.NewScanner(f)
	// for scanner.Scan() {
	// 	r := []rune(scanner.Text())

	// 	if len(r) > len(inputWord) {
	// 		skippedDueToLength++
	// 		continue
	// 	}

	// 	cont := true

	// 	for _, v := range r {
	// 		if strDict[v] != 1 {
	// 			skippedDueToChar++
	// 			cont = false
	// 			break
	// 		}
	// 	}

	// 	if !cont {
	// 		continue
	// 	}

	// 	for _, v := range r {
	// 		if parent.Childern[v] == nil {
	// 			node := &Node{Childern: make(map[rune]*Node), IsWord: false, Value: v}
	// 			parent.Childern[v] = node
	// 		}

	// 		parent = parent.Childern[v]
	// 	}
	// 	parent.IsWord = true
	// 	//start from the root for the new word
	// 	parent = topParent
	// }

	// fmt.Printf("Skipped because of length: %d, Skipped because chars don't exist in provided word: %d.  Total skipped: %d \n", skippedDueToLength, skippedDueToChar, (skippedDueToChar + skippedDueToLength))
	// topParent.

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

	lookup([]byte(inputWord))
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
