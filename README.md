# word-generator

Problem: Work out the number of possible combinations of letters in a given word.
Example: Using the British English dictionary on Ubunut (/usr/share/dict/british-english) we end up with 159 words for an input word "planes"

Details:

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

Therefore, all the letters that make up the word "planets" can be arranged in 13699 different ways.  Note: my solution treats duplicate letters as unique. 


