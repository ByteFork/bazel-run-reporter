package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	r := []byte(`
		<testsuites time="15.682687">
			<testsuite name="Tests.Registration" time="6.605871">
				<testcase name="testCase1" classname="Tests.Registration" time="2.113871" />
				<testcase name="testCase2" classname="Tests.Registration" time="1.051" />
				<testcase name="testCase3" classname="Tests.Registration" time="3.441" />
			</testsuite>
			<testsuite name="Tests.Authentication" time="9.076816">
				<testsuite name="Tests.Authentication.Login" time="4.356">
					<testcase name="testCase4" classname="Tests.Authentication.Login" time="2.244" />
					<testcase name="testCase5" classname="Tests.Authentication.Login" time="0.781" />
					<testcase name="testCase6" classname="Tests.Authentication.Login" time="1.331" />
				</testsuite>
				<testcase name="testCase7" classname="Tests.Authentication" time="2.508" />
				<testcase name="testCase8" classname="Tests.Authentication" time="1.230816" />
				<testcase name="testCase9" classname="Tests.Authentication" time="0.982">
					<failure message="Assertion error message" type="AssertionError">
						<!-- Call stack printed here -->
					</failure>
				</testcase>
			</testsuite>
		</testsuites>
	`)

	parsed, err := Parse(r)
	if err != nil {
		t.Fatalf("Failed to parse valid XML: %v", err)
	}

	if len(parsed.TestSuites) != 2 {
		t.Errorf("Expected 1 test suite, got %d", len(parsed.TestSuites))
	}

	suite := parsed.TestSuites[1]
	if suite.Name != "Tests.Authentication" {
		t.Errorf("Expected suite name 'Tests.Authentication', got '%s'", suite.Name)
	}

	if len(suite.TestCases) != 3 {
		t.Errorf("Expected 3 test cases, got %d", len(suite.TestCases))
	}

	if !suite.TestCases[2].IsFailure() {
		t.Errorf("Expected second test case to be a failure")
	}
}

func TestTestSuiteCompute(t *testing.T) {
	suite := TestSuite{
		Name: "ComputeTest",
		TestCases: []TestCase{
			{Name: "PassTest", ClassName: "example.Test"},
			{Name: "FailTest", ClassName: "example.Test", Failure: &Failure{Message: "Failed"}},
			{Name: "ErrorTest", ClassName: "example.Test", Error: &Error{Message: "Errored"}},
			{Name: "SkipTest", ClassName: "example.Test", Skipped: &Skipped{Message: "Skipped"}},
		},
	}

	suite.Compute()

	if suite.Tests != 4 {
		t.Errorf("Expected 4 tests, got %d", suite.Tests)
	}

	if suite.Failures != 1 {
		t.Errorf("Expected 1 failure, got %d", suite.Failures)
	}

	if suite.Errors != 1 {
		t.Errorf("Expected 1 error, got %d", suite.Errors)
	}

	if suite.Skipped != 1 {
		t.Errorf("Expected 1 skipped, got %d", suite.Skipped)
	}
}

func TestTestSuitesCompute(t *testing.T) {
	suites := TestSuites{
		TestSuites: []TestSuite{
			{
				Name: "Suite1",
				TestCases: []TestCase{
					{Name: "Test1", ClassName: "example.Test"},
					{Name: "Test2", ClassName: "example.Test", Failure: &Failure{Message: "Failed"}},
				},
			},
			{
				Name: "Suite2",
				TestCases: []TestCase{
					{Name: "Test3", ClassName: "example.Test", Error: &Error{Message: "Errored"}},
					{Name: "Test4", ClassName: "example.Test", Skipped: &Skipped{Message: "Skipped"}},
				},
			},
		},
	}

	suites.Compute()

	if suites.Tests != 4 {
		t.Errorf("Expected 4 tests, got %d", suites.Tests)
	}

	if suites.Failures != 1 {
		t.Errorf("Expected 1 failure, got %d", suites.Failures)
	}

	if suites.Errors != 1 {
		t.Errorf("Expected 1 error, got %d", suites.Errors)
	}

	if suites.Skipped != 1 {
		t.Errorf("Expected 1 skipped, got %d", suites.Skipped)
	}
}

func TestMergeTestSuites(t *testing.T) {
	report1 := TestSuites{
		TestSuites: []TestSuite{
			{
				Name: "Suite1",
				TestCases: []TestCase{
					{Name: "Test1", ClassName: "example.Test"},
					{Name: "Test2", ClassName: "example.Test", Failure: &Failure{Message: "Failed"}},
				},
			},
		},
	}

	report2 := TestSuites{
		TestSuites: []TestSuite{
			{
				Name: "Suite1", // Same suite name as in report1
				TestCases: []TestCase{
					{Name: "Test3", ClassName: "example.Test"},
				},
			},
			{
				Name: "Suite2", // New suite
				TestCases: []TestCase{
					{Name: "Test4", ClassName: "example.Test"},
				},
			},
		},
	}

	merged := MergeTestSuites(report1, report2)

	if len(merged.TestSuites) != 2 {
		t.Errorf("Expected 2 test suites, got %d", len(merged.TestSuites))
	}

	// Find Suite1 and check it has 3 tests
	for _, suite := range merged.TestSuites {
		if suite.Name == "Suite1" {
			if len(suite.TestCases) != 3 {
				t.Errorf("Expected Suite1 to have 3 test cases, got %d", len(suite.TestCases))
			}
		}
	}
}

func TestAddTestCasesUnique(t *testing.T) {
	suite := TestSuite{
		Name: "TestAddTestCases",
		TestCases: []TestCase{
			{Name: "Test1", ClassName: "example.Test"},
			{Name: "Test2", ClassName: "example.Test"},
		},
	}

	// Add new test cases, including one with same name/classname
	suite.AddTestCases(true, []TestCase{
		{Name: "Test2", ClassName: "example.Test", Failure: &Failure{Message: "Failed"}}, // Should replace existing
		{Name: "Test3", ClassName: "example.Test"},                                       // Should add new
	}...)

	if len(suite.TestCases) != 3 {
		t.Errorf("Expected 3 test cases after adding, got %d", len(suite.TestCases))
	}

	// Check that Test2 was replaced and now has a failure
	for _, tc := range suite.TestCases {
		if tc.Name == "Test2" && tc.ClassName == "example.Test" {
			if !tc.IsFailure() {
				t.Errorf("Expected Test2 to be marked as failure after replacement")
			}
		}
	}
}

func TestAddTestCasesNonUnique(t *testing.T) {
	suite := TestSuite{
		Name: "TestAddTestCases",
		TestCases: []TestCase{
			{Name: "Test1", ClassName: "example.Test"},
			{Name: "Test2", ClassName: "example.Test"},
		},
	}

	// Add new test cases with unique=false
	suite.AddTestCases(false, []TestCase{
		{Name: "Test2", ClassName: "example.Test", Failure: &Failure{Message: "Failed"}}, // Should add, not replace
		{Name: "Test3", ClassName: "example.Test"},                                       // Should add
	}...)

	if len(suite.TestCases) != 4 {
		t.Errorf("Expected 4 test cases after adding, got %d", len(suite.TestCases))
	}

	count := 0

	for _, tc := range suite.TestCases {
		if tc.Name == "Test2" && tc.ClassName == "example.Test" {
			count++
		}
	}

	if count != 2 {
		t.Errorf("Expected 2 instances of Test2, got %d", count)
	}
}

func TestTestCasesMethod(t *testing.T) {
	suites := TestSuites{
		TestSuites: []TestSuite{
			{
				Name: "Suite1",
				TestCases: []TestCase{
					{Name: "Test1", ClassName: "example.Test"},
					{Name: "Test2", ClassName: "example.Test"},
				},
			},
			{
				Name: "Suite2",
				TestCases: []TestCase{
					{Name: "Test3", ClassName: "example.Test"},
				},
			},
		},
	}

	testCases := suites.TestCases()
	if len(testCases) != 3 {
		t.Errorf("Expected 3 test cases in total, got %d", len(testCases))
	}
}

// Integration test using actual files in testdata directory.
func TestParseActualFiles(t *testing.T) {
	testDataDir := "testdata"
	if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
		t.Skipf("Skipping test - testdata directory doesn't exist")
	}

	var files []string

	err := filepath.Walk(testDataDir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".xml" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk testdata directory: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No XML files found in testdata directory")
	}

	parsedSuites := make([]TestSuites, 0, len(files))

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", file, err)
			}

			parsed, err := Parse(data)
			if err != nil {
				t.Fatalf("Failed to parse file %s: %v", file, err)
			}

			// Simple validation - should have at least one test suite
			if len(parsed.TestSuites) == 0 {
				t.Errorf("No test suites found in %s", file)
			}

			parsedSuites = append(parsedSuites, parsed)
		})
	}

	m := MergeTestSuites(parsedSuites...)
	switch len(m.TestSuites) {
	case 0:
		t.Error("No test suites found in merged result")
	case 2:
		// This is the expected case, no error needed
	default:
		t.Errorf("Expected 2 test suites in merged result, got %d", len(m.TestSuites))
	}
}

func TestIsError(t *testing.T) {
	tests := []struct {
		name     string
		testCase TestCase
		want     bool
	}{
		{
			name:     "No error",
			testCase: TestCase{Name: "Test"},
			want:     false,
		},
		{
			name:     "With error object",
			testCase: TestCase{Name: "Test", Error: &Error{Message: "Error"}},
			want:     true,
		},
		{
			name:     "With status=error",
			testCase: TestCase{Name: "Test", Status: "error"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.testCase.IsError(); got != tt.want {
				t.Errorf("TestCase.IsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFailure(t *testing.T) {
	tests := []struct {
		name     string
		testCase TestCase
		want     bool
	}{
		{
			name:     "No failure",
			testCase: TestCase{Name: "Test"},
			want:     false,
		},
		{
			name:     "With failure object",
			testCase: TestCase{Name: "Test", Failure: &Failure{Message: "Failed"}},
			want:     true,
		},
		{
			name:     "With status=failure",
			testCase: TestCase{Name: "Test", Status: "failure"},
			want:     true,
		},
		{
			name:     "With status=failed",
			testCase: TestCase{Name: "Test", Status: "failed"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.testCase.IsFailure(); got != tt.want {
				t.Errorf("TestCase.IsFailure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSkipped(t *testing.T) {
	tests := []struct {
		name     string
		testCase TestCase
		want     bool
	}{
		{
			name:     "Not skipped",
			testCase: TestCase{Name: "Test"},
			want:     false,
		},
		{
			name:     "With skipped object",
			testCase: TestCase{Name: "Test", Skipped: &Skipped{Message: "Skipped"}},
			want:     true,
		},
		{
			name:     "With status=skipped",
			testCase: TestCase{Name: "Test", Status: "skipped"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.testCase.IsSkipped(); got != tt.want {
				t.Errorf("TestCase.IsSkipped() = %v, want %v", got, tt.want)
			}
		})
	}
}
