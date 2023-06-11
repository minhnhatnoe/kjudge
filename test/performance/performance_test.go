package performance

import (
	"testing"

	"github.com/natsukagami/kjudge/test/performance/suites"
)

var allTests = []*suites.PerfTestSet{
	suites.BigInputProblem(),
	suites.BigOutputProblem(),
	suites.SpawnTimeProblem(),
	suites.TLEProblem(),
	suites.MemoryProblem(),
	suites.CalcProblem(),
}
var allSandbox = []string{"isolate", "raw"}

func BenchmarkAll(b *testing.B) {
	RunBenchmark(b, allTests, allSandbox)
}
