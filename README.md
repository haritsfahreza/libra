[![Go Reference](https://pkg.go.dev/badge/github.com/haritsfahreza/libra.svg)](https://pkg.go.dev/github.com/haritsfahreza/libra)
[![Build Status](https://travis-ci.org/haritsfahreza/libra.svg?branch=master)](https://travis-ci.org/haritsfahreza/libra)
[![codecov](https://codecov.io/gh/haritsfahreza/libra/branch/master/graph/badge.svg)](https://codecov.io/gh/haritsfahreza/libra)
[![Go Report Card](https://goreportcard.com/badge/github.com/haritsfahreza/libra)](https://goreportcard.com/report/github.com/haritsfahreza/libra)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

# libra

Libra is a Go library to compare the interface (struct, maps, etc) and spot the differences between two of them.

## Installing

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

Please see [examples](https://pkg.go.dev/github.com/haritsfahreza/libra#ex-Compare--Struct) for the other usage references

## Contributing

Please read [CONTRIBUTING.md](https://github.com/haritsfahreza/libra/blob/master/CODE_OF_CONDUCT.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/haritsfahreza/libra/tags).

## Authors

- **Harits Fahreza Christyonotoputra** - _Initial work_ - [haritsfahreza](https://github.com/haritsfahreza)

See also the list of [contributors](https://github.com/haritsfahreza/libra/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

- With `libra`, you can speed up your Go software development, especially when you build the object auditing and data versioning system.
- This project is inspired by [`Javers`](https://javers.org).
