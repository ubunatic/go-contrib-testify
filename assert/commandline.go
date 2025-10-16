package assert

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

// CommandLineTestSpec specifies a test to be run via 'go test' command line.
//
// ⚠️ PREVIEW: This API might change in the future! ⚠️
type CommandLineTestSpec struct {
	Name        string   // test name, e.g. "TestFailInsideEventually"
	PackagePath string   // package path, e.g. "./require", default: "./..."
	Args        []string // command line arguments, if empty use default args

	ExpectFailure       bool // if true, the test is expected to fail (i.e. return a non-zero exit code)
	ExpectedErrorLogs   int  // number of unexpected errors
	ExpectedSuccessLogs int  // number of expected successful tests

	ExpectedLineMatches map[string]int // strings that must be found in the output lines (each line is counted at most once per match)

	ErrorMarker   string // string that marks an unexpected error in the output, default: "❌"
	SuccessMarker string // string that marks a successful test in the output, default: "✅"
}

func (s CommandLineTestSpec) withDefaults(t TestingT) CommandLineTestSpec {
	if s.Name == "" {
		FailNow(t, "testSpec.name must be set")
	}

	if s.PackagePath == "" {
		s.PackagePath = "./..."
	}
	if len(s.Args) == 0 {
		s.Args = []string{"-v", "-race", "-count=1", "-run", fmt.Sprintf("^%s$", s.Name)}
	}
	if s.ErrorMarker == "" {
		s.ErrorMarker = "❌"
	}
	if s.SuccessMarker == "" {
		s.SuccessMarker = "✅"
	}
	return s
}

func (s CommandLineTestSpec) assertLineMatches(t *testing.T, lines []string) bool {
	t.Helper()
	if s.ExpectedLineMatches == nil {
		return t.Failed()
	}
	for match, expCount := range s.ExpectedLineMatches {
		found := 0
		for _, line := range lines {
			if strings.Contains(line, match) {
				found++
			}
		}
		if Equal(t, expCount, found, "Expected line match %q not found the expected number of times", match) {
			t.Log("Found expected line match", fmt.Sprintf("%q", match), "the expected number of times:", expCount)
		}
	}
	return !t.Failed()
}

// CommandLineTest runs 'go test' with the specified arguments and checks the output.
//
// ⚠️ PREVIEW: This API might change in the future! ⚠️
//
// Use it run tests that must fail, panic, or behave in a special way that would break the
// current test process. It returns true if the test passed, false if it failed.
// Usage example:
//
//	func TestSomethingViaCommandLine(t *testing.T) {
//		spec := assert.CommandLineTestSpec{
//			Name:                "TestSomething",
//			ExpectedErrorLogs:   1,
//			ExpectedSuccessLogs: 1,
//		}
//		assert.CommandLineTest(t, spec)
//	}
//
//	func TestSomething(t *testing.T) {
//		if os.Getenv("TestSomething") != "1" {
//			t.Skip("skipping test that must be run via CommandLineTest")
//		}
//		t.Log("✅ <-- default success marker")
//		t.Error("❌ <-- default failure marker")
//	}
//
// Also see "TestFailInsideEventuallyViaCommandLine" example in require/requirements_test.go
func CommandLineTest(t *testing.T, spec CommandLineTestSpec) bool {
	t.Helper()
	spec = spec.withDefaults(t)

	t.Run(spec.Name, func(t *testing.T) {
		t.Helper()
		t.Setenv(spec.Name, "1") // signal to the test to run
		args := append([]string{"test"}, spec.Args...)
		args = append(args, spec.PackagePath)
		cmd := exec.Command("go", args...)

		out, err := cmd.CombinedOutput()
		if spec.ExpectFailure {
			Error(t, err, "'go test' command must return an error due to expected failure")
		} else {
			NoError(t, err, "'go test' command must not return an error")
		}

		extractTestNameExp := regexp.MustCompile(spec.Name + `[^ ]*`)

		observedErrors := 0
		observedSuccesses := 0
		name := ""
		lines := strings.Split(string(out), "\n")

		for _, line := range lines {
			// line = strings.TrimSpace(line)
			if strings.HasPrefix(strings.TrimSpace(line), "=== RUN   ") {
				name = extractTestNameExp.FindString(line)
				t.Log("Running test", name)
			}

			if strings.Contains(line, spec.ErrorMarker) {
				// fmt.Println(name, line, fmt.Sprintf("<-- matches error marker %q", spec.ErrorMarker))
				t.Logf("Test %s logged unexpected error:", name)
				fmt.Println(line)
				observedErrors++
			}
			if strings.Contains(line, spec.SuccessMarker) {
				// fmt.Println(name, line, fmt.Sprintf("<-- matches success marker %q", spec.SuccessMarker))
				t.Logf("Test %s logged expected success:", name)
				fmt.Println(line)
				observedSuccesses++
			}
		}

		spec.assertLineMatches(t, lines)

		Equal(t, spec.ExpectedErrorLogs, observedErrors, "Unexpected number of error markers, see output")
		Equal(t, spec.ExpectedSuccessLogs, observedSuccesses, "Unexpected number of success markers, see output")
	})

	return !t.Failed()
}
