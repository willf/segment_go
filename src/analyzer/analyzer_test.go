package analyzer_test

import (
	. "analyzer"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"testing"
)

// From https://github.com/benbjohnson/testing. slightly modified.
// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[33m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ShouldNotError fails the test if an err is not nil.
func ShouldNotError(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[33m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// ShouldError fails the test if an err is nil.
func ShoulError(tb testing.TB, err error) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[33m%s:%d: should have errored. \033[39m\n\n", filepath.Base(file), line)
		tb.FailNow()
	}
}

// ShouldEqual fails the test if exp is not equal to act.
func ShouldEqual(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[33m%s:%d:\tExpected: %#v; got: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// ShouldNotEqual fails the test if exp is equal to act.
func ShouldNotEqual(tb testing.TB, exp, act interface{}) {
	if reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[33m%s:%d:\tShould not equal: got: %#v\033[39m\n\n", filepath.Base(file), line, exp)
		tb.FailNow()
	}
}

func makePath() string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, "..", "..", "data")
}

func TestNew(t *testing.T) {
	model_path := makePath()
	model_name := "test_bigram"
	analyzer := New(model_path, model_name, 20)
	lp, _ := analyzer.LogProb("the")
	ShouldNotEqual(t, lp, math.Inf(-1))
}

func TestUnigramFound(t *testing.T) {
	model_path := makePath()
	model_name := "test_bigram"
	analyzer := New(model_path, model_name, 20)
	lp, _ := analyzer.LogProb("the")
	ShouldNotEqual(t, lp, math.Inf(-1))
}

func TestUnigramNotFound(t *testing.T) {
	model_path := makePath()
	model_name := "test_bigram"
	analyzer := New(model_path, model_name, 20)
	lp, _ := analyzer.LogProb("NOTFOUND")
	ShouldEqual(t, lp, math.Inf(-1))
}

func TestBigramFound(t *testing.T) {
	model_path := makePath()
	model_name := "test_bigram"
	analyzer := New(model_path, model_name, 20)
	lp, _ := analyzer.LogProbTextGivenPrevious("the", "of")
	lp2, _ := analyzer.LogProb("the")
	ShouldNotEqual(t, lp, lp2)
}

func TestBigramNotFound(t *testing.T) {
	model_path := makePath()
	model_name := "test_bigram"
	analyzer := New(model_path, model_name, 20)
	lp, _ := analyzer.LogProbTextGivenPrevious("NOTFOUND", "NOTFOUND")
	ShouldEqual(t, lp, math.Inf(-1))
}

func TestSplit(t *testing.T) {
	a := New("", "", 20) // valid for this test
	got := a.Split("abc")
	// this assumes the order, but that's not right
	expected := Splits{TwoStrings{"a", "bc"}, TwoStrings{"ab", "c"}, TwoStrings{"abc", ""}}
	ShouldEqual(t, expected, got)
}

func TestSplitWithMultiByte(t *testing.T) {
	a := New("", "", 20) // valid for this test
	got := a.Split("日本語")
	// this assumes the order, but that's not right
	expected := Splits{TwoStrings{"日", "本語"}, TwoStrings{"日本", "語"}, TwoStrings{"日本語", ""}}
	ShouldEqual(t, expected, got)
}

func TestSplitWithMax(t *testing.T) {
	a := New("", "", 2) // valid for this test
	got := a.Split("abc")
	// this assumes the order, but that's not right
	expected := Splits{TwoStrings{"a", "bc"}, TwoStrings{"ab", "c"}}
	ShouldEqual(t, expected, got)
}

func TestCombine(t *testing.T) {
	var previous ProbTuple
	previous.LogProb = 1.0
	previous.Tokens = []string{"b", "c"}
	pfirst := 2.0
	first := "a"
	var expected ProbTuple
	expected.LogProb = 3.0
	expected.Tokens = []string{"a", "b", "c"}
	got := previous.Combine(pfirst, first)
	ShouldEqual(t, expected, got)
}

func TestProbTuplesSort(t *testing.T) {
	tups := ProbTuples{
		ProbTuple{LogProb: 1.0},
		ProbTuple{LogProb: 4.0},
		ProbTuple{LogProb: 1.0}}
	sort.Sort(tups)
	ShouldEqual(t, ProbTuple{LogProb: 1.0}, tups[0])
}

func TestProbTuplesMax(t *testing.T) {
	tups := ProbTuples{
		ProbTuple{LogProb: -49.49},
		ProbTuple{LogProb: -33.83},
		ProbTuple{LogProb: -13.29}}
	max, err := tups.Max()
	ShouldEqual(t, max, tups[2])
	ShouldNotError(t, err)
}

func TestProbTuplesMaxEmpty(t *testing.T) {
	tups := ProbTuples{}
	_, err := tups.Max()
	ShoulError(t, err)
}

func TestSegmentRecurseWhenEmpty(t *testing.T) {
	model_path := makePath()
	model_name := "test_bigram"
	analyzer := New(model_path, model_name, 20)
	text := ""
	previous := "<S>"
	n := 0
	memo := make(map[string]ProbTuple)
	s := analyzer.SegmentRecurse(text, previous, n, memo)
	ShouldEqual(t, s.Tokens, []string{})
	ShouldEqual(t, 0.0, s.LogProb)
}

func TestSegmentRecurseWhenMemoized(t *testing.T) {
	model_path := makePath()
	model_name := "test_bigram"
	analyzer := New(model_path, model_name, 20)
	text := "the"
	previous := "<S>"
	n := 0
	memo := make(map[string]ProbTuple)
	expected := ProbTuple{Tokens: []string{"the"}, LogProb: -6.0}
	memo[text] = expected
	got := analyzer.SegmentRecurse(text, previous, n, memo)
	ShouldEqual(t, expected, got)

}

func TestSegment(t *testing.T) {
	model_path := makePath()
	model_name := "small"
	analyzer := New(model_path, model_name, 20)
	text := "theboywholived"
	// we really don't care about this LogProb here ...
	expected := ProbTuple{Tokens: []string{"the", "boy", "who", "lived"}, LogProb: 0.0}
	got := analyzer.Segment(text)
	ShouldEqual(t, expected.Tokens, got)
}
