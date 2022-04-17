package main

import (
	"strconv"
	"strings"
	"sync"
)

var (
	sm             sync.Map
	FilterSuffixes = []string{
		".gif",
		".css",
		".jpg",
		".jpeg",
		".png",
		".zip",
		".svg",
	}
)

func isUnique(url string) bool {
	_, present := sm.Load(url)
	if present {
		return false
	}
	sm.Store(url, true)
	return true
}

func filterTypes(res string) bool {
	for _, suff := range FilterSuffixes {
		if strings.Contains(res, suff) {
			return false
		}
	}
	return true
}

func filterNums(str string, n int) bool {
	c := 0
	for _, character := range str {
		_, err := strconv.Atoi(string(character))
		if err == nil {
			c++
		}
	}
	return (c < n)
}
