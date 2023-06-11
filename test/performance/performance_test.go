package performance

import "testing"

func BenchmarkAll(b *testing.B) {
	RunBenchmark(b, AllTests, AllSandbox)
}
