// Paackage analyzer contains the code for segmentation analysis
package analyzer

import (
	"errors"
	"path/filepath"
	"strings"

	pd "github.com/willf/segment_go/lib/probability_distribution"
)

// An Analyzer contains two data tables of unigram and
// bigram estimated log probabilities.
// It assummes data is in a whitespace delimted files
// frequencies.tsv -- unigram raw frequencies
// 2_frequencies.tsv -- bigram raw frequences
// Here is an exmaple of bigram data:
// >> of	the	4090128330
// >> ,	and	2800837645
// >> in	the	2053960926
// >> ,	the	1364701650
// >> to	the	1358602774
// total.tsv and 2_total.tsv contains a single (perhaps large) integer
// containing the total number of tokens.
// Here is an example:
// >> 431675447550
// The log probabilty of "of the", given the above data, is estimated
// as log2(4090128330/431675447550), or about -6.72.
// ModelPath points to a path with bigram data,
// and ModelName the directory with the frequency data.

func anaMakePath(model_path string, model_name string, prefix string, filename string) string {
	return filepath.Join(model_path, model_name, prefix+filename)
}

type Analyzer struct {
	ModelPath                      string
	ModelName                      string
	MaxTokenLength                 int
	UnigramProbabilityDistribution *pd.ProbabilityDistribution
	BigramProbabilityDistribution  *pd.ProbabilityDistribution
}

// New creates a new Analyzer with data at model_path/model_name,
// cutting off possible token length at max_token_length
func New(model_path string, model_name string, max_token_length int) *Analyzer {
	var a Analyzer
	a.ModelPath = model_path
	a.ModelName = model_name
	a.MaxTokenLength = max_token_length
	a.UnigramProbabilityDistribution =
		pd.New(
			anaMakePath(model_path, model_name, "", "total.tsv"),
			anaMakePath(model_path, model_name, "", "frequencies.tsv"))
	a.BigramProbabilityDistribution =
		pd.New(
			anaMakePath(model_path, model_name, "2_", "total.tsv"),
			anaMakePath(model_path, model_name, "2_", "frequencies.tsv"))
	return &a
}

// LogProb returns the log probabilty of a token string (treated as
// a unigram). Ok is false if the token is not in the table, in which
// case the log probabilty is undefined
func (a *Analyzer) LogProb(token string) (lp float64, ok bool) {
	lp, ok = a.UnigramProbabilityDistribution.LogProb(token)
	return
}

// LogProbTextGivenPrevious returns the log probabilty of a token string
// given another token. If the bigram pair is not found, then the results
// are the same as calling LogProb on the token.
func (a *Analyzer) LogProbTextGivenPrevious(token string, previous string) (lp float64, ok bool) {
	key := previous + " " + token
	lp, ok = a.BigramProbabilityDistribution.LogProb(key)
	if ok {
		return
	}
	lp, ok = a.UnigramProbabilityDistribution.LogProb(token)
	return
}

// TwoStrings is a tuple of two token strings
type TwoStrings [2]string

// Splits are list of TwoStrings
type Splits []TwoStrings

// Split eturns all the splits of a string up to a given length
// Unicode aware
func (a *Analyzer) Split(text string) (results Splits) {

	runes := []rune(text)
	max := len(runes)
	if max > a.MaxTokenLength {
		max = a.MaxTokenLength
	}
	results = make(Splits, max)
	for i := 0; i < max; i++ {
		var t TwoStrings
		t[0] = string(runes[0 : i+1])
		t[1] = string(runes[i+1 : len(runes)])
		results[i] = t
	}
	return
}

// ProbTuple is a log probability, array of tokens tuple
type ProbTuple struct {
	LogProb float64
	Tokens  []string
}

// ProbTuples is a list of ProbTuples
type ProbTuples []ProbTuple

// Len returns the number of ProbTuples. For sort interface.
func (slice ProbTuples) Len() int {
	return len(slice)
}

// Less returns whether a ProbTuple at i is less than one at j.
func (slice ProbTuples) Less(i, j int) bool {
	if slice[i].LogProb == slice[j].LogProb {
		return strings.Join(slice[i].Tokens, "") < strings.Join(slice[j].Tokens, "")
	}
	return slice[i].LogProb < slice[j].LogProb

}

// Swap swaps two ProbTuple instances at i and j
func (slice ProbTuples) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Max returns the ProbTuple with the highest log probability, or errors/
// for empty lists. In case of ties, the result is not necessarily
// deterministic.
func (slice ProbTuples) Max() (pt ProbTuple, err error) {
	max_i := 0
	if len(slice) == 0 {
		err = errors.New("Index out of range")
		return
	}
	for i, _ := range slice[0:] {
		if slice.Less(max_i, i) {
			max_i = i
		}
	}
	pt = slice[max_i]
	return

}

// Combine combines a log probability and a token with a previously
// found tuple. It add the log probabiliities, and adds the token to the
// front of the array
func (prevTuple *ProbTuple) Combine(pfirst float64, first string) (result ProbTuple) {
	result.LogProb = pfirst + prevTuple.LogProb
	result.Tokens = append([]string{first}, prevTuple.Tokens...)
	return
}

// SegmentRecurse recurses over the segmentation
func (a *Analyzer) SegmentRecurse(text string, previous string, n int, memo map[string]ProbTuple) (probTuple ProbTuple) {
	// case 1: empty string
	if len(text) == 0 {
		probTuple.LogProb = 0.0 // Hmm.
		probTuple.Tokens = []string{}
		return
	}
	// case 2: memoized segmentation
	val, found := memo[text]
	if found {
		probTuple = val
		return
	}
	// case 3: we have work to do!
	pts := make(ProbTuples, 0)
	splits := a.Split(text)
	for _, twoStrings := range splits {
		first := twoStrings[0]
		rem := twoStrings[1]
		log_p, _ := a.LogProbTextGivenPrevious(first, previous)
		recurse := a.SegmentRecurse(rem, first, n+1, memo)
		current := recurse.Combine(log_p, first)
		pts = append(pts, current)
	}
	probTuple, _ = pts.Max()
	memo[text] = probTuple

	return
}

// Segment returns the best segmentation of a string
func (a *Analyzer) Segment(text string) (tokens []string) {
	previous := "<S>"
	n := 0
	memo := make(map[string]ProbTuple)
	probTuple := a.SegmentRecurse(text, previous, n, memo)
	tokens = probTuple.Tokens
	return
}
