package performance

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/natsukagami/kjudge/test/performance/suites"
)

var testList = []*suites.PerfTestSet{
	// suites.BigInputProblem(),
	// suites.BigOutputProblem(),
	// suites.SpawnTimeProblem(),
	suites.TLEProblem(),
	suites.MemoryProblem(),
	suites.CalcProblem(),
}
var sandboxList = []string{"isolate", "raw"}

func getTempDir(b *testing.B) (string, bool) {
	value, exists := os.LookupEnv("TEMP_DIR")
	if exists {
		log.Printf("saving DB to %v", value)
		return value, false
	}
	return b.TempDir(), true
}

func BenchmarkSandboxes(b *testing.B) {
	log.Println("creating test DB")

	ctx, err := NewBenchmarkContext(getTempDir(b))
	if err != nil {
		b.Fatal(err)
	}
	defer ctx.Close()

	for _, testset := range testList {
		log.Printf("creating problem %v", testset.Name)
		if err := ctx.writeProblem(testset); err != nil {
			b.Fatal(err)
		}
	}

	for _, testset := range testList {
		for _, sandboxName := range sandboxList {
			testName := fmt.Sprintf("%v %v", testset.Name, sandboxName)
			b.Run(testName, func(b *testing.B) { RunSingleTest(b, ctx, testset, sandboxName) })
		}
	}
}
