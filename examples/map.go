package main

import (
	"fmt"

	"github.com/haritsfahreza/libra"
)

func compareMapExample() {
	oldPerson := map[string]interface{}{
		"Name":           "Gopher",
		"Age":            10,
		"Weight":         50.0,
		"IsMarried":      false,
		"Hobbies":        []string{"Coding"},
		"Numbers":        []int{0, 1, 2},
		"AdditionalInfo": "I love Golang",
	}

	newPerson := map[string]interface{}{
		"Name":           "Gopher",
		"Age":            10,
		"Weight":         60.0,
		"IsMarried":      false,
		"Hobbies":        []string{"Hacking"},
		"Numbers":        []int{1, 2, 3},
		"AdditionalInfo": "I love Golang so much",
	}

	diffs, err := libra.Compare(nil, oldPerson, newPerson)
	if err != nil {
		panic(err)
	}

	for i, diff := range diffs {
		fmt.Printf("#%d : ChangeType=%s Field=%s ObjectType=%s Old='%v' New='%v'\n", i, diff.ChangeType, diff.Field, diff.ObjectType, diff.Old, diff.New)
	}
}
