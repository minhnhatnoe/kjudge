package suites

import "math/rand"

const tleCode = `#include <time.h>
int main(){
	while (time(NULL));
}
`

// Problem with code that will TLE and not optimized away
// Problem: Loop until C time() returns 0
func TLEProblem() *PerfTestSet {
	return &PerfTestSet{
		Name:    "TLE",
		Count:   100,
		CapTime: 500,
		Generator: func(r *rand.Rand) []byte {
			return []byte("\n")
		},
		Solution: []byte(tleCode),
	}
}
