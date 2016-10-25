package probability_distribution_test

import (
	"math"
	"os"
	"path/filepath"
	. "probability_distribution"
	"testing"
)

func makePath(model string, filename string) string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, "..", "..", "data", model, filename)
}

func TestPath(t *testing.T) {
	data_path := makePath("test_unigram", "total.tsv")
	_, err := os.Stat(data_path)

	if err != nil {
		t.Errorf("Test file %s does not exist", data_path)
	}
}

func TestNew(t *testing.T) {
	total_path := makePath("test_unigram", "total.tsv")
	data_path := makePath("test_unigram", "frequencies.tsv")
	pd := New(total_path, data_path)
	if pd.LogTotal == 0.0 {
		t.Errorf("Unable to create pd from %s and %s", total_path, data_path)
	}
}

func TestNewNoSuchFile(t *testing.T) {
	total_path := makePath("NOTFOUND", "total.tsv")
	data_path := makePath("NOTFOUND", "frequencies.tsv")
	pd := New(total_path, data_path)
	if pd.LogTotal != 32.0 {
		t.Errorf("Failed to create empty pd from %s and %s, logTotal is %s", total_path, data_path, pd.LogTotal)
	}
}

func TestLogProbFound(t *testing.T) {
	total_path := makePath("test_unigram", "total.tsv")
	data_path := makePath("test_unigram", "frequencies.tsv")
	pd := New(total_path, data_path)
	lp, _ := pd.LogProb("the")
	if lp == math.Inf(-1) {
		t.Errorf("Expected to find a log prob for %s, but didn't", "the")
	}
}

func TestBigramLogProbFound(t *testing.T) {
	total_path := makePath("test_bigram", "2_total.tsv")
	data_path := makePath("test_bigram", "2_frequencies.tsv")
	pd := New(total_path, data_path)
	lp, _ := pd.LogProb("of the")
	if lp == math.Inf(-1) {
		t.Errorf("Expected to find a log prob for %s, but didn't", "of the")
	}
}

func TestLogProbNotFound(t *testing.T) {
	total_path := makePath("test_unigram", "total.tsv")
	data_path := makePath("test_unigram", "frequencies.tsv")
	pd := New(total_path, data_path)
	lp, _ := pd.LogProb("NOTFOUND")
	if lp != math.Inf(-1) {
		t.Errorf("Expected to not find a log prob for %s, but did", "NOTFOUND")
	}
}
