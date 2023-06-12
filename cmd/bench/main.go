// Run benchmarks a limited number of times (probably once for each test)
// to test kjudge's judging system.
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/natsukagami/kjudge/test/performance"
	"github.com/natsukagami/kjudge/test/performance/suites"
)

var (
	tmpDir       = flag.String("file", "", "Path to database. Will be created and auto-removed if not specified.")
	sandboxNames = flag.String("sandbox", "", "Sandbox implementations to put in the test. Comma seperated.")
	testNameList = flag.String("test", "", "Tests to use. Comma seperated.")
)

func main() {
	flag.Parse()
	sandboxList := parseSandboxList()
	testList := parseTestList()

	log.Println("creating test DB")
	dts, err := performance.NewDataset(*tmpDir, testList)
	if err != nil {
		log.Fatal(err)
	}
	defer dts.Close()

	for _, testset := range testList {
		for _, sandboxName := range sandboxList {
			testName := fmt.Sprintf("%v %v", testset.Name, sandboxName)
			result := testing.Benchmark(func(b *testing.B) {
				performance.RunSingleTest(b, dts, testset, sandboxName, testName)
			})
			fmt.Printf("%v: %v\n", testName, result)
		}
	}
}

func parseSandboxList() []string {
	var sandboxList []string
	if *sandboxNames == "" {
		sandboxList = performance.AllSandbox
	} else {
		sandboxList = strings.Split(*sandboxNames, ",")
	}
	return sandboxList
}

func parseTestList() []*suites.PerfTestSet {
	var testList []*suites.PerfTestSet
	if *testNameList == "" {
		testList = performance.AllTests
	} else {
		testNames := make(map[string]bool)
		for _, testName := range strings.Split(*testNameList, ",") {
			testNames[testName] = true
		}

		for _, test := range performance.AllTests {
			value, exists := testNames[test.Name]
			if exists {
				if !value {
					log.Panicf("duplicated name %v", test.Name)
				}
				testList = append(testList, test)
			}
		}
		if len(testList) != len(testNames) {
			log.Panicf("cannot find all tests. found %v/%v", len(testList), len(testNames))
		}
	}
	return testList
}
