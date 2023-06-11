package suites

import (
	"fmt"
	"math/rand"
)

const bigOutputCode = `#include <stdio.h>
int main(){
	int a; scanf("%i", &a);
	while (a--) putchar('A');
}
`

// 500MB output problem to compare output syscall time
// Problem: Input one number n, then output 'A' n times.
func BigOutputProblem() *PerfTestSet {
	maxSize := 500000000 // Around 500MB
	return &PerfTestSet{
		Name:    "OUTPUT",
		Count:   100,
		CapTime: 100,
		Generator: func(r *rand.Rand) []byte {
			outputSize := maxSize - r.Intn(100)
			input := fmt.Sprintf("%v", outputSize)
			return []byte(input)
		},
		Solution: []byte(bigOutputCode),
	}
}
