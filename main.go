package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Result struct {
	Type    string
	Message string
}

var (
	Queue   = make(chan *url.URL)
	Results = make(chan Result)
)

func reader() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		words := strings.Fields(s.Text())
		for _, word := range words {
			if strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "http://") {
				parsed, err := url.Parse(word)
				if err != nil {
					continue
				}
				Queue <- parsed
			}
		}
	}
	close(Queue)
}

func writer(filter *bool, verbose *bool) {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for res := range Results {
		res.Message = strings.TrimSpace(res.Message)
		if isUnique(res.Message) && res.Message != "" {
			if (filterTypes(res.Message) && filterNums(res.Message, 4)) || !(*filter) {
				if *verbose {
					fmt.Fprintf(w, "[%s] %s\n", res.Type, res.Message)
				} else {
					fmt.Fprintln(w, res.Message)
				}
			}
		}
	}
}

func main() {
	// options
	_ = flag.Bool("", false, "Uses all parts of URL by default.")
	filter := flag.Bool("filter", false, "Filter images and css.")
	mode := flag.String("mode", "wordlist", "Options: wordlist, endpoints")
	verbose := flag.Bool("s", false, "Show source of output.")
	keys := flag.Bool("keys", false, "Use parameter keys.")
	vals := flag.Bool("vals", false, "Use parameter values.")
	paths := flag.Bool("paths", false, "Use URL paths.")
	domains := flag.Bool("domains", false, "Use domain names.")
	flag.Parse()

	// if none specified, show all
	if !(*keys) && !(*vals) && !(*paths) && !(*domains) {
		*keys = true
		*vals = true
		*paths = true
		*domains = true
	}

	// check for stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "No input detected, use `cat urls.txt | words`")
		os.Exit(1)
	}

	go reader()

	switch {
	case *mode == "endpoints":
		go getEndpoints()
	case *mode == "wordlist":
		go buildWordlist(keys, vals, paths, domains)
	}

	writer(filter, verbose)
}
