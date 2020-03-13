package models

import (
	"git.nkagami.me/natsukagami/kjudge/db"
	"git.nkagami.me/natsukagami/kjudge/models/verify"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// TestGroupWithTests are wrapped test groups with tests included.
type TestGroupWithTests struct {
	*TestGroup
	Tests []*Test
}

// GetProblemTests collects test groups and tests from a problem.
func GetProblemTests(db db.DBContext, problemID int) ([]*TestGroupWithTests, error) {
	return getProblemTests(db, problemID, "*")
}

// GetProblemTestsMeta is like GetProblemTests, but inputs and outputs are not included.
func GetProblemTestsMeta(db db.DBContext, problemID int) ([]*TestGroupWithTests, error) {
	return getProblemTests(db, problemID, "id, name, test_group_id")
}

// GetProblemTests but allow us to omit rows (input, output)
func getProblemTests(db db.DBContext, problemID int, rows string) ([]*TestGroupWithTests, error) {
	testGroups, err := GetProblemTestGroups(db, problemID)
	if err != nil {
		return nil, err
	}
	// Collect the ID list
	var (
		IDs   []interface{}
		tgMap = make(map[int]*TestGroupWithTests)
	)
	for _, tg := range testGroups {
		IDs = append(IDs, tg.ID)
		tgMap[tg.ID] = &TestGroupWithTests{TestGroup: tg}
	}
	if len(IDs) == 0 {
		// Empty
		return nil, nil
	}
	// Query the tests
	query, params, err := sqlx.In("SELECT "+rows+" FROM tests WHERE test_group_id IN (?) ORDER BY name", IDs)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var tests []*Test
	if err := db.Select(&tests, query, params...); err != nil {
		return nil, errors.WithStack(err)
	}
	for _, test := range tests {
		tg := tgMap[test.TestGroupID]
		tg.Tests = append(tg.Tests, test)
	}
	// Collect the map into a slice.
	var res []*TestGroupWithTests
	for _, tg := range testGroups {
		res = append(res, tgMap[tg.ID])
	}
	return res, nil
}

// Verify verifies Test's contents.
func (r *Test) Verify() error {
	if r.Input == nil {
		return errors.New("input must not be null")
	}
	if r.Output == nil {
		return errors.New("output must not be null")
	}
	return errors.Wrapf(verify.Names(r.Name), "field name")
}