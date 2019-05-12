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
//Prepare the objects that you want to compare
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
if err != nil {
	panic(err)
}

//Use the diffs array with your own purpose e.g. printout
for i, diff := range diffs {
	fmt.Printf("#%d : ChangeType=%s Field=%s ObjectType=%s Old='%v' New='%v'\n", i, diff.ChangeType, diff.Field, diff.ObjectType, diff.Old, diff.New)
}
```
Please see [examples](https://github.com/haritsfahreza/libra/tree/master/examples) for the other usage references

## License
See [LICENSE](LICENSE)
