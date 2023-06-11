package suites

import "math/rand"

const memoryCode = `#include <stdio.h>
#include <algorithm>
#include <random>

// 1GB of memory
const int SIZE = (1<<30)/sizeof(int);
int values[SIZE];

int main(){
	int a; scanf("%i", &a);
	
	// Pretty random access
	std::iota(values, values + SIZE, 0);
	std::shuffle(values, values + SIZE, std::mt19937(a));
	
	// Assure nothing will be optimized away
	printf("%i", values[a % SIZE]);
}
`

// Random access 1GB int array with std::shuffle
func MemoryProblem() *PerfTestSet {
	return &PerfTestSet{
		Name:    "MEMORY",
		Count:   10,
		CapTime: 5000,
		Generator: func(r *rand.Rand) []byte {
			return []byte("\n")
		},
		Solution: []byte(memoryCode),
	}
}
