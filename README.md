# Lazytest 

[![Go Report Card Badge](http://goreportcard.com/badge/gophergala2016/lazytest)](http://goreportcard.com/report/gophergala2016/lazytest)

A continuous test runner for Go.

Once started, it will listen for file changes in a given directory. If a file change is detected, only the tests affected by that file change will be re-run. 

### Usage:
````
  -exclude string
      exclude paths (default "/vendor/")
  -extensions string
      file extensions to watch (default "go,tpl,html")
  -root string
      watch root (default ".")
````
