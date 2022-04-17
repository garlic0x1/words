/*
Takes dirty input with urls as input, generates a wordlist based on the components
*/

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

func getEndpoints() {
	for parsed := range Queue {
		newPath := ""
		for _, str := range strings.Split(parsed.Path, "/") {
			if !strings.Contains(str, "=") && !strings.Contains(str, ";") {
				newPath += "/" + str
			} else {
				break
			}
		}
		parsed.Path = newPath

		Results <- Result{
			Type:    "endpoint",
			Message: parsed.Scheme + "://" + parsed.Host + parsed.Path,
		}
	}
	close(Results)
}

func buildWordlist(keys *bool, vals *bool, paths *bool, domains *bool) {
	for parsed := range Queue {
		host := parsed.Host
		path := parsed.Path
		query := parsed.Query()

		if *domains {
			for _, w := range strings.Split(host, ".") {
				Results <- Result{
					Type:    "domain",
					Message: w,
				}
			}
		}
		if *paths {
			for _, w := range strings.Split(path, "/") {
				Results <- Result{
					Type:    "path",
					Message: w,
				}
			}
		}
		for k, values := range query {
			if *keys {
				Results <- Result{
					Type:    "key",
					Message: k,
				}
			}
			if *vals {
				for _, v := range values {
					/*
						for _, v2 := range strings.Split(v1, "/") {
							Results <- v2
						}
					*/
					Results <- Result{
						Type:    "val",
						Message: v,
					}
				}
			}
		}
	}
	close(Results)
}

func main() {
	// options
	_ = flag.Bool("", false, "Uses all parts of URL by default.")
	filter := flag.Bool("filter", false, "Filter images and css.")
	endpoints := flag.Bool("endpoints", false, "Output unique endpoints from list of URLs.")
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

	if *endpoints {
		go getEndpoints()
	} else {
		go buildWordlist(keys, vals, paths, domains)
	}

	writer(filter, verbose)
}
