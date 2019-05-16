package libra_test

import (
	"fmt"
	"testing"

	"github.com/haritsfahreza/libra"
)

type person struct {
	Name      string
	Age       int
	Weight    float64
	IsMarried bool
	Hobbies   []string
	Numbers   []int
	Ignore    string `libra:"ignore"`
}

type anotherPerson struct {
	Name string
}

func TestCompare(t *testing.T) {
	t.Run("failed when all nil", func(t *testing.T) {
		if _, err := libra.Compare(nil, nil, nil); err == nil {
			t.Errorf("Error must not be nil. expected: %s actual: %v", "all values cannot be nil", err)
		}
	})

	t.Run("failed when different type", func(t *testing.T) {
		if _, err := libra.Compare(nil, "", 1); err == nil {
			t.Errorf("Error must not be nil. expected: %s actual: %v", "different values type", err)
		}

		if _, err := libra.Compare(nil, person{
			Name: "test1",
		}, anotherPerson{
			Name: "test1",
		}); err == nil {
			t.Errorf("Error must not be nil. expected: %s actual: %v", "different values type", err)
		}
	})

	t.Run("succeed when create new object", func(t *testing.T) {
		diffs, err := libra.Compare(nil, nil, person{
			Name: "test1",
			Age:  22,
		})
		if err != nil {
			t.Errorf("Error must be nil. Got: %v", err)
		}
		if len(diffs) != 1 {
			t.Errorf("Compare result length must be only one. expected: %d actual: %d", 1, len(diffs))
		}

		if diffs[0].ChangeType != libra.New {
			t.Errorf("Invalid ChangeType. expected: %s actual: %s", libra.New, diffs[0].ChangeType)
		}

		if diffs[0].ObjectType != "libra_test.person" {
			t.Errorf("Invalid ObjectType. expected: %s actual: %s", "libra_test.person", diffs[0].ObjectType)
		}

		if diffs[0].New == nil {
			t.Errorf("Invalid New value. expected: %s actual: %v", "not nil", diffs[0].New)
		}

		if diffs[0].Old != nil {
			t.Errorf("Invalid Old value. expected: %s actual: %v", "nil", diffs[0].Old)
		}
	})

	t.Run("succeed when removed an object", func(t *testing.T) {
		diffs, err := libra.Compare(nil, person{
			Name: "test1",
			Age:  22,
		}, nil)
		if err != nil {
			t.Errorf("Error must be nil. Got: %v", err)
		}
		if len(diffs) != 1 {
			t.Errorf("Compare result length must be only one. expected: %d actual: %d", 1, len(diffs))
		}

		if diffs[0].ChangeType != libra.Removed {
			t.Errorf("Invalid ChangeType. expected: %s actual: %s", libra.Removed, diffs[0].ChangeType)
		}

		if diffs[0].ObjectType != "libra_test.person" {
			t.Errorf("Invalid ObjectType. expected: %s actual: %s", "libra_test.person", diffs[0].ObjectType)
		}

		if diffs[0].New != nil {
			t.Errorf("Invalid New value. expected: %s actual: %v", "nil", diffs[0].New)
		}

		if diffs[0].Old == nil {
			t.Errorf("Invalid Old value. expected: %s actual: %v", "not nil", diffs[0].Old)
		}
	})

	t.Run("succeed when changed two objects", func(t *testing.T) {
		diffs, err := libra.Compare(nil, person{
			Name:      "test1",
			Age:       22,
			Weight:    float64(80),
			IsMarried: true,
			Hobbies:   []string{"Swimming"},
			Numbers:   []int{1, 2},
		}, person{
			Name:      "test1",
			Age:       23,
			Weight:    float64(85),
			IsMarried: true,
			Hobbies:   []string{"Swimming", "Hiking"},
			Numbers:   []int{2},
		})
		if err != nil {
			t.Errorf("Error must be nil. Got: %v", err)
		}

		expectedFields := []string{"Age", "Weight", "Hobbies", "Numbers"}
		if len(diffs) != len(expectedFields) {
			t.Errorf("Invalid result length. expected: %d actual: %d", len(expectedFields), len(diffs))
		} else {
			for i := 0; i < len(diffs); i++ {
				if diffs[i].ChangeType != libra.Changed {
					t.Errorf("Invalid diffs[%d].ChangeType. expected: %s actual: %s", i, libra.Changed, diffs[i].ChangeType)
				}

				if diffs[i].ObjectType != "libra_test.person" {
					t.Errorf("Invalid diffs[%d].ObjectType. expected: %s actual: %s", i, "libra_test.person", diffs[i].ObjectType)
				}

				if diffs[i].Field != expectedFields[i] {
					t.Errorf("Invalid diffs[%d].Field. expected: %s actual: %s", i, expectedFields[i], diffs[i].Field)
				}

				if diffs[i].Old == diffs[i].New {
					t.Errorf("diffs[%d].Old must be different with diffs[%d].New. old: %s new: %v", i, i, diffs[i].Old, diffs[i].New)
				}
			}
		}
	})

	t.Run("succeed when changed two maps", func(t *testing.T) {
		diffs, err := libra.Compare(nil, map[string]interface{}{"Age": 22, "Weight": 80}, map[string]interface{}{"Age": 23, "Weight": 80})
		if err != nil {
			t.Errorf("Error must be nil. Got: %v", err)
		}

		expectedFields := []string{"Age"}
		if len(diffs) != len(expectedFields) {
			t.Errorf("Invalid result length. expected: %d actual: %d", len(expectedFields), len(diffs))
		} else {
			for i := 0; i < len(diffs); i++ {
				if diffs[i].ChangeType != libra.Changed {
					t.Errorf("Invalid diffs[%d].ChangeType. expected: %s actual: %s", i, libra.Changed, diffs[i].ChangeType)
				}

				if diffs[i].ObjectType != "map[string]interface {}" {
					t.Errorf("Invalid diffs[%d].ObjectType. expected: %s actual: %s", i, "map[string]interface {}", diffs[i].ObjectType)
				}

				if diffs[i].Field != expectedFields[i] {
					t.Errorf("Invalid diffs[%d].Field. expected: %s actual: %s", i, expectedFields[i], diffs[i].Field)
				}

				if diffs[i].Old == diffs[i].New {
					t.Errorf("diffs[%d].Old must be different with diffs[%d].New. old: %s new: %v", i, i, diffs[i].Old, diffs[i].New)
				}
			}
		}
	})

	t.Run("failed when the values has different type", func(t *testing.T) {
		_, err := libra.Compare(nil, map[string]interface{}{"Age": "A", "Weight": 80}, map[string]interface{}{"Age": 12, "Weight": 80})
		if err == nil {
			t.Errorf("Error must not be nil. expected: %s actual: %v", "different values type", err)
		}
	})
}

var bDiffs []libra.Diff

func benchmarkCompare(old, new interface{}, b *testing.B) {
	var diffs []libra.Diff
	var err error
	for i := 0; i < b.N; i++ {
		diffs, err = libra.Compare(nil, old, new)
		if err != nil {
			panic(err)
		}
	}
	bDiffs = diffs
}
func BenchmarkCompareStruct(b *testing.B) {
	benchmarkCompare(person{
		Name:      "test1",
		Age:       22,
		Weight:    float64(80),
		IsMarried: true,
		Hobbies:   []string{"Swimming"},
		Numbers:   []int{1, 2},
	}, person{
		Name:      "test1",
		Age:       23,
		Weight:    float64(85),
		IsMarried: true,
		Hobbies:   []string{"Swimming", "Hiking"},
		Numbers:   []int{2},
	}, b)
}

func BenchmarkCompareMap(b *testing.B) {
	benchmarkCompare(map[string]interface{}{
		"Name":           "Gopher",
		"Age":            10,
		"Weight":         50.0,
		"IsMarried":      false,
		"Hobbies":        []string{"Coding"},
		"Numbers":        []int{0, 1, 2},
		"AdditionalInfo": "I love Golang",
	}, map[string]interface{}{
		"Name":           "Gopher",
		"Age":            10,
		"Weight":         60.0,
		"IsMarried":      false,
		"Hobbies":        []string{"Hacking"},
		"Numbers":        []int{1, 2, 3},
		"AdditionalInfo": "I love Golang so much",
	}, b)
}

func ExampleCompare_struct() {
	oldPerson := person{
		Name:      "Gopher",
		Age:       10,
		Weight:    50.0,
		IsMarried: false,
		Hobbies:   []string{"Coding"},
		Numbers:   []int{0, 1, 2},
	}

	newPerson := person{
		Name:      "Gopher",
		Age:       10,
		Weight:    60.0,
		IsMarried: false,
		Hobbies:   []string{"Hacking"},
		Numbers:   []int{1, 2, 3},
	}

	diffs, err := libra.Compare(nil, oldPerson, newPerson)
	if err != nil {
		panic(err)
	}

	for i, diff := range diffs {
		fmt.Printf("#%d : ChangeType=%s Field=%s ObjectType=%s Old='%v' New='%v'\n", i, diff.ChangeType, diff.Field, diff.ObjectType, diff.Old, diff.New)
	}
	// Output:
	// #0 : ChangeType=changed Field=Weight ObjectType=libra_test.person Old='50' New='60'
	// #1 : ChangeType=changed Field=Hobbies ObjectType=libra_test.person Old='Coding' New='Hacking'
	// #2 : ChangeType=changed Field=Numbers ObjectType=libra_test.person Old='0,1,2' New='1,2,3'
}
