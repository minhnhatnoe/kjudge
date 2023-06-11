package performance

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/natsukagami/kjudge/db"
	"github.com/natsukagami/kjudge/models"
	"github.com/natsukagami/kjudge/server/auth"
	"github.com/natsukagami/kjudge/test/performance/suites"
	"github.com/pkg/errors"
)

type BenchmarkContext struct {
	tdir     string
	remove   bool
	db       *db.DB
	user     *models.User
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

func NewBenchmarkContext(tmpDir string, testList []*suites.PerfTestSet) (*BenchmarkContext, error) {
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

	ctx := &BenchmarkContext{tdir: tmpDir, remove: removeDB, db: benchDB, problems: make(map[string]*models.Problem)}

	if err := ctx.writeContest(); err != nil {
		ctx.Close()
		return nil, err
	}

	if err := ctx.writeUser(); err != nil {
		ctx.Close()
		return nil, err
	}

	if err := ctx.writeTestList(testList); err != nil {
		ctx.Close()
		return nil, err
	}

	return ctx, nil
}

func (ctx *BenchmarkContext) Close() error {
	if err := ctx.db.Close(); err != nil {
		return err
	}

	if ctx.remove {
		if err := os.RemoveAll(ctx.tdir); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *BenchmarkContext) writeContest() error {
	ctx.contest = &models.Contest{
		ContestType:          "weighted",
		StartTime:            time.Now().AddDate(0, 0, -1),
		EndTime:              time.Now().AddDate(0, 0, 1),
		ID:                   0,
		Name:                 "Performance Testing",
		ScoreboardViewStatus: models.ScoreboardViewStatusPublic,
	}
	return errors.Wrapf(ctx.contest.Write(ctx.db), "creating contest")
}

func (ctx *BenchmarkContext) writeUser() error {
	password, err := auth.PasswordHash("bigquestions")
	if err != nil {
		return errors.Wrap(err, "hashing password")
	}
	ctx.user = &models.User{
		ID:           "Iroh",
		Password:     string(password),
		DisplayName:  "The Dragon of the West",
		Organization: "Order of the White Lotus",
	}
	return errors.Wrap(ctx.user.Write(ctx.db), "creating user")
}

func (ctx *BenchmarkContext) writeTestList(testList []*suites.PerfTestSet) error {
	for _, testset := range testList {
		log.Printf("creating problem %v", testset.Name)
		problem, err := testset.AddToDB(ctx.db, 2403, len(ctx.problems)+1, ctx.contest.ID)
		if err != nil {
			return errors.Wrapf(err, "creating testset %v's problem", testset.Name)
		}
		
		ctx.problems[testset.Name] = problem
	}
	return nil
}

const testSolution = `#include "solution.hpp"`

func (ctx *BenchmarkContext) writeSolutions(problemName string, N int) error {
	for i := 0; i < N; i++ {
		problem := ctx.problems[problemName]
		sub := models.Submission{
			ProblemID:   problem.ID,
			UserID:      ctx.user.ID,
			Source:      []byte(testSolution),
			Language:    models.LanguageCpp,
			SubmittedAt: time.Now(),
			Verdict:     models.VerdictIsInQueue,
		}
		if err := sub.Write(ctx.db); err != nil {
			return err
		}

		job := models.NewJobScore(sub.ID)
		if err := job.Write(ctx.db); err != nil {
			return err
		}
	}
	return nil
}

// TODO
func (ctx *BenchmarkContext) assertRunComplete(testset *suites.PerfTestSet) error {
	return nil
}
