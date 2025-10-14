package require

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// AssertionTesterInterface defines an interface to be used for testing assertion methods
type AssertionTesterInterface interface {
	TestMethod()
}

// AssertionTesterConformingObject is an object that conforms to the AssertionTesterInterface interface
type AssertionTesterConformingObject struct {
}

func (a *AssertionTesterConformingObject) TestMethod() {
}

// AssertionTesterNonConformingObject is an object that does not conform to the AssertionTesterInterface interface
type AssertionTesterNonConformingObject struct {
}

type MockT struct {
	// Failed marks the test as failed.
	Failed bool
	// finished marks the test as finished, indicating that FailNow was called
	// and no further code should be executed after that.
	finished bool
}

// Helper is like [testing.T.Helper] but does nothing.
func (MockT) Helper() {}

func (t *MockT) FailNow() {
	t.Failed = true
	t.finished = true
}

func (t *MockT) Finished() bool { return t.finished }

func (t *MockT) Errorf(format string, args ...interface{}) {
	_, _ = format, args
}

func TestImplements(t *testing.T) {
	t.Parallel()

	Implements(t, (*AssertionTesterInterface)(nil), new(AssertionTesterConformingObject))

	mockT := new(MockT)
	Implements(mockT, (*AssertionTesterInterface)(nil), new(AssertionTesterNonConformingObject))
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestIsType(t *testing.T) {
	t.Parallel()

	IsType(t, new(AssertionTesterConformingObject), new(AssertionTesterConformingObject))

	mockT := new(MockT)
	IsType(mockT, new(AssertionTesterConformingObject), new(AssertionTesterNonConformingObject))
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestEqual(t *testing.T) {
	t.Parallel()

	Equal(t, 1, 1)

	mockT := new(MockT)
	Equal(mockT, 1, 2)
	if !mockT.Failed {
		t.Error("Check should fail")
	}

}

func TestNotEqual(t *testing.T) {
	t.Parallel()

	NotEqual(t, 1, 2)
	mockT := new(MockT)
	NotEqual(mockT, 2, 2)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestExactly(t *testing.T) {
	t.Parallel()

	a := float32(1)
	b := float32(1)
	c := float64(1)

	Exactly(t, a, b)

	mockT := new(MockT)
	Exactly(mockT, a, c)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestNotNil(t *testing.T) {
	t.Parallel()

	NotNil(t, new(AssertionTesterConformingObject))

	mockT := new(MockT)
	NotNil(mockT, nil)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestNil(t *testing.T) {
	t.Parallel()

	Nil(t, nil)

	mockT := new(MockT)
	Nil(mockT, new(AssertionTesterConformingObject))
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestTrue(t *testing.T) {
	t.Parallel()

	True(t, true)

	mockT := new(MockT)
	True(mockT, false)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestFalse(t *testing.T) {
	t.Parallel()

	False(t, false)

	mockT := new(MockT)
	False(mockT, true)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	Contains(t, "Hello World", "Hello")

	mockT := new(MockT)
	Contains(mockT, "Hello World", "Salut")
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestNotContains(t *testing.T) {
	t.Parallel()

	NotContains(t, "Hello World", "Hello!")

	mockT := new(MockT)
	NotContains(mockT, "Hello World", "Hello")
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestPanics(t *testing.T) {
	t.Parallel()

	Panics(t, func() {
		panic("Panic!")
	})

	mockT := new(MockT)
	Panics(mockT, func() {})
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestNotPanics(t *testing.T) {
	t.Parallel()

	NotPanics(t, func() {})

	mockT := new(MockT)
	NotPanics(mockT, func() {
		panic("Panic!")
	})
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestNoError(t *testing.T) {
	t.Parallel()

	NoError(t, nil)

	mockT := new(MockT)
	NoError(mockT, errors.New("some error"))
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	Error(t, errors.New("some error"))

	mockT := new(MockT)
	Error(mockT, nil)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestErrorContains(t *testing.T) {
	t.Parallel()

	ErrorContains(t, errors.New("some error: another error"), "some error")

	mockT := new(MockT)
	ErrorContains(mockT, errors.New("some error"), "different error")
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestEqualError(t *testing.T) {
	t.Parallel()

	EqualError(t, errors.New("some error"), "some error")

	mockT := new(MockT)
	EqualError(mockT, errors.New("some error"), "Not some error")
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestEmpty(t *testing.T) {
	t.Parallel()

	Empty(t, "")

	mockT := new(MockT)
	Empty(mockT, "x")
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestNotEmpty(t *testing.T) {
	t.Parallel()

	NotEmpty(t, "x")

	mockT := new(MockT)
	NotEmpty(mockT, "")
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestWithinDuration(t *testing.T) {
	t.Parallel()

	a := time.Now()
	b := a.Add(10 * time.Second)

	WithinDuration(t, a, b, 15*time.Second)

	mockT := new(MockT)
	WithinDuration(mockT, a, b, 5*time.Second)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestInDelta(t *testing.T) {
	t.Parallel()

	InDelta(t, 1.001, 1, 0.01)

	mockT := new(MockT)
	InDelta(mockT, 1, 2, 0.5)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestZero(t *testing.T) {
	t.Parallel()

	Zero(t, "")

	mockT := new(MockT)
	Zero(mockT, "x")
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestNotZero(t *testing.T) {
	t.Parallel()

	NotZero(t, "x")

	mockT := new(MockT)
	NotZero(mockT, "")
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestJSONEq_EqualSONString(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	JSONEq(mockT, `{"hello": "world", "foo": "bar"}`, `{"hello": "world", "foo": "bar"}`)
	if mockT.Failed {
		t.Error("Check should pass")
	}
}

func TestJSONEq_EquivalentButNotEqual(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	JSONEq(mockT, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
	if mockT.Failed {
		t.Error("Check should pass")
	}
}

func TestJSONEq_HashOfArraysAndHashes(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	JSONEq(mockT, "{\r\n\t\"numeric\": 1.5,\r\n\t\"array\": [{\"foo\": \"bar\"}, 1, \"string\", [\"nested\", \"array\", 5.5]],\r\n\t\"hash\": {\"nested\": \"hash\", \"nested_slice\": [\"this\", \"is\", \"nested\"]},\r\n\t\"string\": \"foo\"\r\n}",
		"{\r\n\t\"numeric\": 1.5,\r\n\t\"hash\": {\"nested\": \"hash\", \"nested_slice\": [\"this\", \"is\", \"nested\"]},\r\n\t\"string\": \"foo\",\r\n\t\"array\": [{\"foo\": \"bar\"}, 1, \"string\", [\"nested\", \"array\", 5.5]]\r\n}")
	if mockT.Failed {
		t.Error("Check should pass")
	}
}

func TestJSONEq_Array(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	JSONEq(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `["foo", {"nested": "hash", "hello": "world"}]`)
	if mockT.Failed {
		t.Error("Check should pass")
	}
}

func TestJSONEq_HashAndArrayNotEquivalent(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	JSONEq(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `{"foo": "bar", {"nested": "hash", "hello": "world"}}`)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestJSONEq_HashesNotEquivalent(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	JSONEq(mockT, `{"foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestJSONEq_ActualIsNotJSON(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	JSONEq(mockT, `{"foo": "bar"}`, "Not JSON")
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestJSONEq_ExpectedIsNotJSON(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	JSONEq(mockT, "Not JSON", `{"foo": "bar", "hello": "world"}`)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestJSONEq_ExpectedAndActualNotJSON(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	JSONEq(mockT, "Not JSON", "Not JSON")
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestJSONEq_ArraysOfDifferentOrder(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	JSONEq(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `[{ "hello": "world", "nested": "hash"}, "foo"]`)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestYAMLEq_EqualYAMLString(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	YAMLEq(mockT, `{"hello": "world", "foo": "bar"}`, `{"hello": "world", "foo": "bar"}`)
	if mockT.Failed {
		t.Error("Check should pass")
	}
}

func TestYAMLEq_EquivalentButNotEqual(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	YAMLEq(mockT, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
	if mockT.Failed {
		t.Error("Check should pass")
	}
}

func TestYAMLEq_HashOfArraysAndHashes(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	expected := `
numeric: 1.5
array:
  - foo: bar
  - 1
  - "string"
  - ["nested", "array", 5.5]
hash:
  nested: hash
  nested_slice: [this, is, nested]
string: "foo"
`

	actual := `
numeric: 1.5
hash:
  nested: hash
  nested_slice: [this, is, nested]
string: "foo"
array:
  - foo: bar
  - 1
  - "string"
  - ["nested", "array", 5.5]
`
	YAMLEq(mockT, expected, actual)
	if mockT.Failed {
		t.Error("Check should pass")
	}
}

func TestYAMLEq_Array(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	YAMLEq(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `["foo", {"nested": "hash", "hello": "world"}]`)
	if mockT.Failed {
		t.Error("Check should pass")
	}
}

func TestYAMLEq_HashAndArrayNotEquivalent(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	YAMLEq(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `{"foo": "bar", {"nested": "hash", "hello": "world"}}`)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestYAMLEq_HashesNotEquivalent(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	YAMLEq(mockT, `{"foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestYAMLEq_ActualIsSimpleString(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	YAMLEq(mockT, `{"foo": "bar"}`, "Simple String")
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestYAMLEq_ExpectedIsSimpleString(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	YAMLEq(mockT, "Simple String", `{"foo": "bar", "hello": "world"}`)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func TestYAMLEq_ExpectedAndActualSimpleString(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	YAMLEq(mockT, "Simple String", "Simple String")
	if mockT.Failed {
		t.Error("Check should pass")
	}
}

func TestYAMLEq_ArraysOfDifferentOrder(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)
	YAMLEq(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `[{ "hello": "world", "nested": "hash"}, "foo"]`)
	if !mockT.Failed {
		t.Error("Check should fail")
	}
}

func ExampleComparisonAssertionFunc() {
	t := &testing.T{} // provided by test

	adder := func(x, y int) int {
		return x + y
	}

	type args struct {
		x int
		y int
	}

	tests := []struct {
		name      string
		args      args
		expect    int
		assertion ComparisonAssertionFunc
	}{
		{"2+2=4", args{2, 2}, 4, Equal},
		{"2+2!=5", args{2, 2}, 5, NotEqual},
		{"2+3==5", args{2, 3}, 5, Exactly},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.expect, adder(tt.args.x, tt.args.y))
		})
	}
}

func TestComparisonAssertionFunc(t *testing.T) {
	t.Parallel()

	type iface interface {
		Name() string
	}

	tests := []struct {
		name      string
		expect    interface{}
		got       interface{}
		assertion ComparisonAssertionFunc
	}{
		{"implements", (*iface)(nil), t, Implements},
		{"isType", (*testing.T)(nil), t, IsType},
		{"equal", t, t, Equal},
		{"equalValues", t, t, EqualValues},
		{"exactly", t, t, Exactly},
		{"notEqual", t, nil, NotEqual},
		{"NotEqualValues", t, nil, NotEqualValues},
		{"notContains", []int{1, 2, 3}, 4, NotContains},
		{"subset", []int{1, 2, 3, 4}, []int{2, 3}, Subset},
		{"notSubset", []int{1, 2, 3, 4}, []int{0, 3}, NotSubset},
		{"elementsMatch", []byte("abc"), []byte("bac"), ElementsMatch},
		{"regexp", "^t.*y$", "testify", Regexp},
		{"notRegexp", "^t.*y$", "Testify", NotRegexp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.expect, tt.got)
		})
	}
}

func ExampleValueAssertionFunc() {
	t := &testing.T{} // provided by test

	dumbParse := func(input string) interface{} {
		var x interface{}
		json.Unmarshal([]byte(input), &x)
		return x
	}

	tests := []struct {
		name      string
		arg       string
		assertion ValueAssertionFunc
	}{
		{"true is not nil", "true", NotNil},
		{"empty string is nil", "", Nil},
		{"zero is not nil", "0", NotNil},
		{"zero is zero", "0", Zero},
		{"false is zero", "false", Zero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, dumbParse(tt.arg))
		})
	}
}

func TestValueAssertionFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     interface{}
		assertion ValueAssertionFunc
	}{
		{"notNil", true, NotNil},
		{"nil", nil, Nil},
		{"empty", []int{}, Empty},
		{"notEmpty", []int{1}, NotEmpty},
		{"zero", false, Zero},
		{"notZero", 42, NotZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.value)
		})
	}
}

func ExampleBoolAssertionFunc() {
	t := &testing.T{} // provided by test

	isOkay := func(x int) bool {
		return x >= 42
	}

	tests := []struct {
		name      string
		arg       int
		assertion BoolAssertionFunc
	}{
		{"-1 is bad", -1, False},
		{"42 is good", 42, True},
		{"41 is bad", 41, False},
		{"45 is cool", 45, True},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, isOkay(tt.arg))
		})
	}
}

func TestBoolAssertionFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     bool
		assertion BoolAssertionFunc
	}{
		{"true", true, True},
		{"false", false, False},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.value)
		})
	}
}

func ExampleErrorAssertionFunc() {
	t := &testing.T{} // provided by test

	dumbParseNum := func(input string, v interface{}) error {
		return json.Unmarshal([]byte(input), v)
	}

	tests := []struct {
		name      string
		arg       string
		assertion ErrorAssertionFunc
	}{
		{"1.2 is number", "1.2", NoError},
		{"1.2.3 not number", "1.2.3", Error},
		{"true is not number", "true", Error},
		{"3 is number", "3", NoError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x float64
			tt.assertion(t, dumbParseNum(tt.arg, &x))
		})
	}
}

func TestErrorAssertionFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		err       error
		assertion ErrorAssertionFunc
	}{
		{"noError", nil, NoError},
		{"error", errors.New("whoops"), Error},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.err)
		})
	}
}

func TestEventuallyWithTFalse(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)

	condition := func(collect *assert.CollectT) {
		True(collect, false)
	}

	EventuallyWithT(mockT, condition, 100*time.Millisecond, 20*time.Millisecond)
	True(t, mockT.Failed, "Check should fail")
}

func TestEventuallyWithTTrue(t *testing.T) {
	t.Parallel()

	mockT := new(MockT)

	counter := 0
	condition := func(collect *assert.CollectT) {
		defer func() {
			counter += 1
		}()
		True(collect, counter == 1)
	}

	EventuallyWithT(mockT, condition, 100*time.Millisecond, 20*time.Millisecond)
	False(t, mockT.Failed, "Check should pass")
	Equal(t, 2, counter, "Condition is expected to be called 2 times")
}

func TestFailInsideEventuallyViaCommandLine(t *testing.T) {
	t.Setenv("TestFailInsideEventually", "1")
	cmd := exec.Command("go", "test", "-v", "-race", "-count=1", "-run", "^TestFailInsideEventually$")
	out, err := cmd.CombinedOutput()
	assert.Error(t, err, "go test for TestFailInsideEventually must fail")
	finishedTests := 0
	observedErrors := 0
	observedConditionFailures := 0
	observedPanics := 0
	extractTestNameExp := regexp.MustCompile(`TestFailInsideEventuallyViaCommandLine[^ ]*`)
	name := ""
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "=== RUN   ") {
			name = extractTestNameExp.FindString(line)
			fmt.Println(line)
		}

		if strings.Contains(line, "❌") {
			fmt.Println(name, line, "<-- unexpected error")
			observedErrors++
		}
		if strings.Contains(line, "✅ FINISHED") {
			fmt.Println(name, line, "<-- expected successful test")
			finishedTests++
		}
		if strings.Contains(line, "Error:") && strings.Contains(line, "Condition") {
			fmt.Println(name, line, "<-- expected 'Condition' message")
			observedConditionFailures++
		}
		if strings.Contains(line, "Panic in condition:") {
			fmt.Println(name, line, "<-- expected 'Condition' panic message")
			observedPanics++
		}
	}

	assert.Equal(t, 0, observedErrors, "Unexpected errors detected, see output")

	// There are 6 tests that are expected to fail.
	// If you change the number of tests in TestFailInsideEventually, please update this number accordingly.
	assert.Equal(t, 6, finishedTests, "Expected number of FINISHED tests not found")

	// There are 5 tests where the condition is never satisfied.
	// One test uses assert.Fail but eventually returns true and thus satisfies the condition.
	// If you change the number of tests in TestFailInsideEventually, please update this number accordingly.
	assert.Equal(t, 5, observedConditionFailures, "Expected number of 'Condition' messages not found, see output")

	// There are 2 tests that panic, so we expect 2 panic messages.
	assert.Equal(t, 2, observedPanics, "Missed expected panics, see output")
}

func TestFailInsideEventually(t *testing.T) {
	if os.Getenv("TestFailInsideEventually") == "" {
		t.Skip("Skipping test, run via TestFailInsideEventuallyViaCommandLine")
	}

	// Using testing.T:
	// Enable this test temporarily to manually check that the issue is fixed.
	// Note that MockT does not reproduce the issue, so we have to use the real *testing.T.
	// TODO: Enhance MockT to reproduce the issue, then we can remove the t.Skip above.
	// UPDATE: I tried to enhance MockT, but it still does not reproduce the issue or fails
	//         in a different way. Therefore I keep using *testing.T for now.
	//         In general, require.MockT does not play well when assert.Fail is called.
	//         The test will not be marked as failed, even though Fail is called.

	// The Bug:
	// Calling require.Fail (or similar) inside require.Eventually will prevent the 'condition'
	// to exit with a result. The channel assignment in the assert.Eventually will block
	// and hang the test until the timeout is reached. There was is no other way to wait
	// for the unclean exit. The changes to assert.Eventually committed with this test
	// fix this issue, by also waiting for an unclean exit of the condition and moreover
	// by handling panics inside the condition gracefully.

	// How to read the test results:
	// - See [TestFailInsideEventuallyViaCommandLine], which automates this
	// - The test will always fail, because it calls Fail or assert.Fail or panics
	// - The test is "successful" if it fails quickly and cleanly, i.e. without
	//   multiple calls to the eventually function, except if expected.
	// - The test logs should contain specific messages, see below.
	// - The test must not log any UNCLEAN EXIT or MISSED ASSERTIONS messages or any errors ❌.

	type test struct {
		Name     string
		Return   bool
		FailFunc func(t TestingT)
		MustStop bool // after the FailFunc is called
	}

	const returnStop = true
	const returnNoStop = false
	const mustStop = true
	const mustNotStop = false

	RequireFail := func(t TestingT) { Fail(t, "💥 fail now") }
	AssertFail := func(t TestingT) { assert.Fail(t, "💥 mark as failed") }
	Panic := func(_ TestingT) { panic("💥 panicking now") }

	for _, tt := range []test{
		// Test cases that must exit immediately after the first call to the condition.
		{"require.Fail must stop", returnStop, RequireFail, mustStop},
		{"require.Fail must stop even if told not to", returnNoStop, RequireFail, mustStop},
		{"assert.Fail must stop if told to", returnStop, AssertFail, mustStop},

		// The following test case is the only one that must not stop and where multiple calls
		// to the condition are expected, because assert.Fail does not stop the execution of the condition.
		{"assert.Fail must not stop if told not to", returnNoStop, AssertFail, mustNotStop},

		// Panics must always stop, because they are not expected and indicate a bug in the code.
		{"panic must stop", returnStop, Panic, mustStop},
		{"panic must stop even if told not to", returnNoStop, Panic, mustStop},

		// Make sure to update the assertions in TestFailInsideEventuallyViaCommandLine
		// accordingly if you change the number of tests here.
	} {
		count := 0
		start := time.Now()
		timeout := time.Second * 1
		tick := time.Second / 3
		ok := false

		ok = t.Run(tt.Name, func(t *testing.T) {
			// Cannot use a MockT here, because it does reproduce the issue.
			Eventually(t, func() bool {
				count++
				t.Log("🪲 eventually call number:", count, "calling FailNow!")
				tt.FailFunc(t)   // any case that calls Fail or assert.Fail should stop retrying
				return tt.Return // indicate whether to stop retrying
			}, timeout, tick)
		})

		dur := time.Since(start)
		t.Log("🪲 test duration:", dur)

		// TODO: Replace with plain t with MockT once it can reproduce the issue.
		// Until can only indicate that the test should have failed using stdout.

		c := new(assert.CollectT)
		assert.True(c, !ok, "❌ UNCLEAN EXIT: test was expected to fail, but passed")

		if tt.MustStop {
			assert.Equal(c, 1, count, "❌ UNCLEAN EXIT: eventually func should be called exactly once")
			assert.Less(c, dur, tick, "❌ UNCLEAN EXIT: eventually func should be called only once, but took too long")
		} else {
			assert.Greater(c, count, 1, "❌ MISSED ASSERTIONS: eventually func should be called multiple times, but was called only once")
			assert.Greater(c, dur, tick, "❌ MISSED ASSERTIONS: eventually func should be called multiple times over time, but total duration was too short")
		}

		if c.Failed() {
			t.Log("❌ TEST FAILED")
		} else {
			t.Log("✅ FINISHED", tt.Name)
		}
	}
}
