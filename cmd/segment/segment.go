package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/willf/segment_go/lib/analyzer"
)

var model_name = flag.String("model", "twitter", "Model name")
var path_name = flag.String("path", "data", "Path name")
var max = flag.Int("max", 20, "Maximum word length")

func main() {
	flag.Parse()
	a := analyzer.New(*path_name, *model_name, *max)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("%s\t%s\n", text, strings.Join(a.Segment(text), " "))
	}
}
