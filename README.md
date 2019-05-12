[![GoDoc](https://godoc.org/github.com/haritsfahreza/libra?status.svg)](https://godoc.org/github.com/haritsfahreza/libra)
[![Build Status](https://travis-ci.org/haritsfahreza/libra.svg?branch=master)](https://travis-ci.org/haritsfahreza/libra)
[![codecov](https://codecov.io/gh/haritsfahreza/libra/branch/master/graph/badge.svg)](https://codecov.io/gh/haritsfahreza/libra)
[![Go Report Card](https://goreportcard.com/badge/github.com/haritsfahreza/libra)](https://goreportcard.com/report/github.com/haritsfahreza/libra)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

# libra
Libra is a Go library to compare the object and spot the difference between two of them

## Installation
```sh
$ go get -u github.com/haritsfahreza/libra
```

## Usage
```go
package main

import "github.com/haritsfahreza/libra"

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
	diffs, err := libra.Compare(context.Background(), oldPerson, newPerson)
}

```

## License
See [LICENSE](LICENSE)
