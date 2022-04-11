# words
Generate wordlists from input containing URLs

# Examples
`cat urls.txt | words -s -filter`  
`cat urls.txt | hakrawler | words -keys -vals`

# Help
```
$ words -h
Usage of words:
  -	Uses all parts of URL by default
  -domains
    	Use domain names.
  -filter
    	Filter images and css.
  -keys
    	Use parameter keys.
  -paths
    	Use URL paths.
  -s	Show source of output.
  -vals
    	Use parameter values.
```
