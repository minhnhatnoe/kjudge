package performance

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/natsukagami/kjudge/test/performance/suites"
	"github.com/natsukagami/kjudge/worker"
	"github.com/natsukagami/kjudge/worker/queue"
	"github.com/natsukagami/kjudge/worker/sandbox"
)

func getTempDir(b *testing.B) (string, bool) {
	value, exists := os.LookupEnv("TEMP_DIR")
	if exists {
		log.Printf("saving DB to %v", value)
		return value, false
	}
	return b.TempDir(), true
}

func RunBenchmark(b *testing.B, testList []*suites.PerfTestSet, sandboxList []string) {
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

func RunSingleTest(b *testing.B, ctx *BenchmarkContext, testset *suites.PerfTestSet, sandboxName string) {
	sandbox, err := worker.NewSandbox(
		sandboxName,
		sandbox.IgnoreWarnings(true),
		sandbox.EnableSandboxLogs(false))
	if err != nil {
		b.Fatal(err)
	}
	
	for i := 0; i < b.N; i++ {
		ctx.writeSolution(testset.Name)
	}

	queue := queue.NewQueue(ctx.db, sandbox, queue.CompileLogs(false), queue.RunLogs(false), queue.ScoreLogs(false))

	b.ResetTimer()
	queue.Run()
	b.StopTimer()

	if err := ctx.assertRunComplete(testset); err != nil {
		b.Fatal(err)
	}
}
