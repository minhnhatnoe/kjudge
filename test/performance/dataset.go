package performance

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/natsukagami/kjudge/db"
	"github.com/natsukagami/kjudge/models"
	"github.com/natsukagami/kjudge/test/performance/suites"
	"github.com/pkg/errors"
)

type BenchmarkDataset struct {
	tdir     string
	remove   bool
	db       *db.DB
	contest  *models.Contest
	problems map[string]*models.Problem
}

var AllTests = []*suites.PerfTestSet{
	suites.BigInputProblem(),
	suites.BigOutputProblem(),
	suites.SpawnTimeProblem(),
	suites.TLEProblem(),
	suites.MemoryProblem(),
	suites.CalcProblem(),
}
var AllSandbox = []string{"isolate", "raw"}

func NewDataset(tmpDir string, testList []*suites.PerfTestSet) (*BenchmarkDataset, error) {
	removeDB := false
	if tmpDir == "" {
		createdDir, err := os.MkdirTemp(os.TempDir(), "kjudge_test")
		if err != nil {
			log.Fatal(err)
		}
		tmpDir = createdDir
		removeDB = true
	}
	log.Printf("saving DB to %v", tmpDir)

	benchDB, err := db.New(filepath.Join(tmpDir, "bench.db"))
	if err != nil {
		return nil, err
	}

	ctx := &BenchmarkDataset{tdir: tmpDir, remove: removeDB, db: benchDB, problems: make(map[string]*models.Problem)}

	if err := ctx.writeContest(); err != nil {
		ctx.Close()
		return nil, err
	}

	if err := ctx.writeTestList(testList); err != nil {
		ctx.Close()
		return nil, err
	}

	return ctx, nil
}

func (dts *BenchmarkDataset) Close() error {
	if err := dts.db.Close(); err != nil {
		return err
	}

	if dts.remove {
		if err := os.RemoveAll(dts.tdir); err != nil {
			return err
		}
	}
	return nil
}

func (dts *BenchmarkDataset) writeContest() error {
	dts.contest = &models.Contest{
		ContestType:          "weighted",
		StartTime:            time.Now().AddDate(0, 0, -1),
		EndTime:              time.Now().AddDate(0, 0, 1),
		ID:                   0,
		Name:                 "Performance Testing",
		ScoreboardViewStatus: models.ScoreboardViewStatusPublic,
	}
	return errors.Wrapf(dts.contest.Write(dts.db), "creating contest")
}

func (dts *BenchmarkDataset) writeTestList(testList []*suites.PerfTestSet) error {
	for _, testset := range testList {
		log.Printf("creating problem %v", testset.Name)
		problem, err := testset.AddToDB(dts.db, 2403, len(dts.problems)+1, dts.contest.ID)
		if err != nil {
			return errors.Wrapf(err, "creating testset %v's problem", testset.Name)
		}

		dts.problems[testset.Name] = problem
	}
	return nil
}

func (dts *BenchmarkDataset) NewContext(name string) (*BenchmarkContext, error) {
	ctx := &BenchmarkContext{dataset: dts}
	if err := ctx.writeUser(name); err != nil {
		return nil, err
	}
	return ctx, nil
}
