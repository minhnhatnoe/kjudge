package performance

import (
	"time"

	"github.com/natsukagami/kjudge/models"
	"github.com/natsukagami/kjudge/server/auth"
	"github.com/natsukagami/kjudge/test/performance/suites"
	"github.com/pkg/errors"
)

type BenchmarkContext struct {
	dataset *BenchmarkDataset
	user    *models.User
}

func (ctx *BenchmarkContext) writeUser(name string) error {
	password, err := auth.PasswordHash("bigquestions")
	if err != nil {
		return errors.Wrap(err, "hashing password")
	}

	ctx.user = &models.User{
		ID:           name,
		Password:     string(password),
		DisplayName:  "The Dragon of the West",
		Organization: "Order of the White Lotus",
	}
	return errors.Wrapf(ctx.user.Write(ctx.dataset.db), "creating user %v", name)
}

const testSolution = `#include "solution.hpp"`

func (ctx *BenchmarkContext) writeSolutions(problemName string, N int) error {
	for i := 0; i < N; i++ {
		problem := ctx.dataset.problems[problemName]
		sub := models.Submission{
			ProblemID:   problem.ID,
			UserID:      ctx.user.ID,
			Source:      []byte(testSolution),
			Language:    models.LanguageCpp,
			SubmittedAt: time.Now(),
			Verdict:     models.VerdictIsInQueue,
		}
		if err := sub.Write(ctx.dataset.db); err != nil {
			return err
		}

		job := models.NewJobScore(sub.ID)
		if err := job.Write(ctx.dataset.db); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *BenchmarkContext) assertRunComplete(testset *suites.PerfTestSet) error {
	return nil
}
