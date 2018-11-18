[![GoDoc](https://godoc.org/github.com/uwuh/tampah?status.svg)](https://godoc.org/github.com/uwuh/tampah)
[![Build Status](https://travis-ci.org/uwuh/tampah.svg?branch=master)](https://travis-ci.org/uwuh/tampah)
[![codecov](https://codecov.io/gh/uwuh/tampah/branch/master/graph/badge.svg)](https://codecov.io/gh/uwuh/tampah)
[![Go Report Card](https://goreportcard.com/badge/github.com/uwuh/tampah)](https://goreportcard.com/report/github.com/uwuh/tampah)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

# tampah
Tampah is a Go library to compare the object and spot the difference between two of them

## Installation
```sh
$ go get -u github.com/uwuh/tampah
```

## Usage
```go
package main

import "github.com/uwuh/tampah"

type Person struct {
	Name      string
	Age       int
	Weight    float64
	IsMarried bool
}

func main() {
	//The object could be from your database or HTTP request
	oldPerson := Person{
		Name:      "Sudirman",
		Age:       22,
		Weight:    float64(80),
		IsMarried: true,
	}
	newPerson := Person{
		Name:      "Sudirman",
		Age:       23,
		Weight:    float64(85),
		IsMarried: true,
	}
	diffs, err := tampah.Compare(context.Background(), oldPerson, newPerson)
}

```

## License
See [LICENSE](LICENSE)
