package suites

import (
	"fmt"
	"math/rand"
)

const calcCode = `#include <stdio.h>
const int MOD = 1'000'000'007;
int main(){
	int a; scanf("%i", &a);
	int result = 2;
	while (a--)
		result = 1LL * result * result % MOD;
	printf("%i", result);
}
`

// Heavy number crunching problem
// Problem: Input one number n, then output 2^(2^n) mod 1e9+7
// Solution: O(n).
// I'm aware that there are better solutions using math magic.
func CalcProblem() *PerfTestSet {
	maxOps := 1000000000 // 1e9
	return &PerfTestSet{
		Name:    "CALC",
		Count:   100,
		CapTime: 5000,
		Generator: func(r *rand.Rand) []byte {
			ops := maxOps - r.Intn(100)
			input := fmt.Sprintf("%v", ops)
			return []byte(input)
		},
		Solution: []byte(calcCode),
	}
}
