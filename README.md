# word-generator

Problem: Work out the number of possible combinations of letters in a given word.
Example: Using the British English dictionary on Ubuntu (/usr/share/dict/british-english) we end up with 159 words for an input word "planes"

## Details:

Given the word "planets", we have a total of 7 letters to work with. We need to find out the number of permutations we can make. Not only from all the 7 words, but also words of any length less then 7. 
To do this we use http://discrete.openmathbooks.org/dmoi2/sec_counting-combperm.html (also see https://en.wikipedia.org/wiki/Summation) which gives us something like:

    (7, 7)  ->  7*6*5*4*3*2*1 (in other words: 7!)
    (7, 6)  ->  7*6*5*4*3*2
    (7, 5)  ->  7*6*5*4*3
    (7, 4)  ->  7*6*5*4
    (7, 3)  ->  7*6*5
    (7, 2)  ->  7*6
    (7, 1)  ->  7

Rewritten to:

    7! + (7*6*5*4*3*2) + (7*6*5*4*3) + (7*6*5*4) + (7*6*5) + (7*6) + 7
    5040 + 5040 + 2520 + 840 + 210 + 42 +7 = 13699

Therefore, all the letters that make up the word "planets" can be arranged in 13699 different ways.  Note: my solution treats duplicate letters as unique, but this is broken because I'm using hashing to decrease lookup time. 

## Aproach

### Brute force

Initially I figured that brute forcing all the permuatation would be done within a reasonable amount of time, however soon enough it dawned on me that as the sample size increases the brute force approach grows expenentiall inefficent.  I found that brute force takes a long time to guess the last 25% of permutations. This is becuase after some time it starts to generete alrady known permutations.  This way the odds of finding the very last permutation decrease significantly as the word size increases.  With small word sizes this is not a problem. 

### Recursion

This appraoch tries to minimize the number of iterations so that they match the number of permutations.  In other words each cycle should generete a new and unique permutation. 

## Results

The results are broken down by CPU and method.  I've selected a 9-character long word for the test.

### "youngster" - brute force - Intel 4930K clocked at 4.30 GHz 

      Skipped because of length: 33146, Skipped because chars don't exist in provided word: 68028.  Total skipped: 101174 
      Number of possibilities: 986409
      Number of Random iterations: 39691600
      Figuring out all the permutations took 34.957715773s
      Traversing tree took 21.390154ms
      Found a total of 264 words.

### "youngster" - brute force - AMD 5800X at 4.75 GHz

      Skipped because of length: 33146, Skipped because chars don't exist in provided word: 68028.  Total skipped: 101174
      Number of possibilities: 986409
      Number of Random iterations: 59485463
      Figuring out all the permutations took 23.754992341s
      Traversing tree took 74.642873ms
      Found a total of 264 words.

### "youngster" - recursive lookup - AMD 5800X at 4.75 GHz

      Skipped because of length: 33146, Skipped because chars don't exist in provided word: 68028.  Total skipped: 101174 
      Number of possibilities for length of 9 is 986409
      Number of possibilities for length of 10 is 9864100
      Number of possibilites generated: 10976173
      Number of iterations to generete all permutations: 10976173
      Time to generete all permutations 881.606448ms
      Traversing tree took 5.53395ms
      Found a total of 264 words.



### Some more results highlighting just how inefficient the brute force approach is as word size increases

Times are from an Intel 4930K clocked at 4.30% Ghz 

### planets (7 chars)

    Skipped because of length: 64144, Skipped because chars don't exist in provided word: 37370.  Total skipped: 101514 
    Number of possibilities: 13699 
    Number of Random iterations: 319454 
    Figuring out all the permutations took 219.055435ms 
    Traversing tree took 963.318µs 
    Found a total of 159 words.
        
### yoghurts (8 chars)

    Skipped because of length: 47961, Skipped because chars don't exist in provided word: 53725.  Total skipped: 101686 
    Number of possibilities: 109600 
    Number of Random iterations: 3451321 
    Figuring out all the permutations took 2.653343008s 
    Traversing tree took 4.319075ms 
    Found a total of 133 words.

    
### youngster (9 chars)

    Skipped because of length: 33146, Skipped because chars don't exist in provided word: 68028.  Total skipped: 101174 
    Number of possibilities: 986409 
    Number of Random iterations: 39691600 
    Figuring out all the permutations took 34.957715773s 
    Traversing tree took 21.390154ms 
    Found a total of 264 words.
    

## Conclusion

Bruteforce is easy to implement and works well up words 8-chars long. However anything over that length we see exponential slowdown.  Using recursion to figure out all the permutations is  clearly the best way forward.  Although my recursion solution generetes some duplicated permutations, it is still a magnatide faster than the brute force solution.  23 seconds vs 881 ms. When trying to figure out why the recursion is generating more permutations I wrote the loops by hand and iterated exactly the number of times expected for a given length.  For the 9-char word the solution finished in 660 ms. Therefore the extra 1 million duplicated permutations cost in the vacinity of 200 ms when the generation permutation. This penatly is not see when traversing the tree in order to find words.  Both cases took 5.5 ms to complete.


