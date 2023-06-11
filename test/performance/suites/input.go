package suites

import "math/rand"

const bigInputCode = `#include <stdio.h>
int main(){
	int r = 0;
	while (char c = getchar()){
		if ('a' <= c && c <= 'b') r++;
		else break;
	}
	printf("%i", r);
}
`

// 500MB input to compare disk read time. O(1) memory.
// Problem: Given a string, print it's length.
func BigInputProblem() *PerfTestSet {
	maxSize := 500000000 // Around 500MB
	return &PerfTestSet{
		Name:    "INPUT",
		Count:   5,
		CapTime: 100,
		Generator: func(r *rand.Rand) []byte {
			strSize := maxSize - r.Intn(10)
			input := make([]byte, strSize)
			for i := range input {
				input[i] = byte(r.Intn(26) + 'A')
			}
			return input
		},
		Solution: []byte(bigInputCode),
	}
}
