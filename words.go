/*
Takes dirty input with urls as input, generates a wordlist based on the components
*/

package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	sm sync.Map
)

func isUnique(url string) bool {
	_, present := sm.Load(url)
	if present {
		return false
	}
	sm.Store(url, true)
	return true
}

func main() {
	// check for stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "No input detected, use `cat scope | gau | words`")
		os.Exit(1)
	}

	results := make(chan string)
	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			words := strings.Fields(s.Text())
			for _, word := range words {
				parsed, err := url.Parse(word)
				if err != nil {
					continue
				}
				host := parsed.Host
				path := parsed.Path
				query := parsed.Query()

				for _, w := range strings.Split(host, ".") {
					results <- w
				}
				for _, w := range strings.Split(path, "/") {
					results <- w
				}
				for k, v := range query {
					results <- k
					for _, v1 := range v {
						for _, v2 := range strings.Split(v1, "/") {
							results <- v2
						}
						results <- v1
					}
				}
			}
		}
		close(results)
	}()

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for res := range results {
		if isUnique(res) && res != "" {
			fmt.Fprintln(w, res)
		}
	}
}
