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

## Results

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

This solution is great for combinations up 8 characters in length.  Anything more than that costs way too much in terms of Inefficiency when brute forcing the combinations.  Although using multiple go routines would probably speed things up even more, ultimately a much better solution is needed to figure out all the possible combinations.
