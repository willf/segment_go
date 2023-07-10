# Segment

A Go-based segmenter.


# To install

go install github.com/willf/segment_go/...@latest

# To run

```bash
❯ segment -h
Usage of segment:
  -max int
        Maximum word length (default 20)
  -model string
        Model name (default "twitter")
  -path string
        Path name (default "data")
```
Note that `-path` needs to point to the directory containing the model files.


```bash
echo "helloworld" | segment -path segment_go/data -model 'google_books'
helloworld  hello world
```
