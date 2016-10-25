package probability_distribution

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

type ProbabilityDistribution struct {
	total_file_path string
	data_file_path  string
	LogTotal        float64
	table           map[string]float64
}

func New(total_file_path string, data_file_path string) *ProbabilityDistribution {
	var pd ProbabilityDistribution
	pd.total_file_path = total_file_path
	pd.data_file_path = data_file_path
	pd.table = make(map[string]float64)
	total, _ := readTotal(total_file_path)
	if total == 0.0 {
		pd.LogTotal = 32.0
	} else {
		pd.LogTotal = math.Log2(total)
	}
	pd.readFrequencies(data_file_path, pd.LogTotal)
	return &pd
}

func readTotal(total_file_path string) (total float64, err error) {
	file, err := os.Open(total_file_path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		total, err = strconv.ParseFloat(text, 64)
		if err == nil {
			return
		}
	}
	return
}

func (pd *ProbabilityDistribution) readFrequencies(data_file_path string, LogTotal float64) (err error) {
	file, err := os.Open(data_file_path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		parts := strings.Fields(text)
		if len(parts) >= 2 {
			key := strings.Join(parts[0:len(parts)-1], " ")
			freq, err := strconv.ParseFloat(parts[len(parts)-1], 64)
			if err != nil {
				break
			}
			pd.table[key] = math.Log2(freq) - LogTotal
		}
	}
	return
}

func (pd *ProbabilityDistribution) LogProb(text string) (log_prob float64, ok bool) {
	log_prob, ok = pd.table[text]
	if !ok {
		log_prob = math.Inf(-1)
	}
	return
}
