package main

import "strings"

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
