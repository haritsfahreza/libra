package libra_test

import (
	"testing"

	"github.com/haritsfahreza/libra"
)

type testStruct struct {
	Name      string
	Age       int
	Weight    float64
	IsMarried bool
	Hobbies   []string
	Numbers   []int
}

type anotherTestStruct struct {
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

		if _, err := libra.Compare(nil, testStruct{
			Name: "test1",
		}, anotherTestStruct{
			Name: "test1",
		}); err == nil {
			t.Errorf("Error must not be nil. expected: %s actual: %v", "different values type", err)
		}
	})

	t.Run("succeed when create new object", func(t *testing.T) {
		diffs, err := libra.Compare(nil, nil, testStruct{
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

		if diffs[0].ObjectType != "libra_test.testStruct" {
			t.Errorf("Invalid ObjectType. expected: %s actual: %s", "libra_test.testStruct", diffs[0].ObjectType)
		}

		if diffs[0].New == nil {
			t.Errorf("Invalid New value. expected: %s actual: %v", "not nil", diffs[0].New)
		}

		if diffs[0].Old != nil {
			t.Errorf("Invalid Old value. expected: %s actual: %v", "nil", diffs[0].Old)
		}
	})

	t.Run("succeed when removed an object", func(t *testing.T) {
		diffs, err := libra.Compare(nil, testStruct{
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

		if diffs[0].ObjectType != "libra_test.testStruct" {
			t.Errorf("Invalid ObjectType. expected: %s actual: %s", "libra_test.testStruct", diffs[0].ObjectType)
		}

		if diffs[0].New != nil {
			t.Errorf("Invalid New value. expected: %s actual: %v", "nil", diffs[0].New)
		}

		if diffs[0].Old == nil {
			t.Errorf("Invalid Old value. expected: %s actual: %v", "not nil", diffs[0].Old)
		}
	})

	t.Run("succeed when changed two objects", func(t *testing.T) {
		diffs, err := libra.Compare(nil, testStruct{
			Name:      "test1",
			Age:       22,
			Weight:    float64(80),
			IsMarried: true,
			Hobbies:   []string{"Swimming"},
			Numbers:   []int{1, 2},
		}, testStruct{
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

				if diffs[i].ObjectType != "libra_test.testStruct" {
					t.Errorf("Invalid diffs[%d].ObjectType. expected: %s actual: %s", i, "libra_test.testStruct", diffs[i].ObjectType)
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

	t.Run("failed when the value is nil", func(t *testing.T) {
		_, err := libra.Compare(nil, map[string]interface{}{"Age": "A", "Weight": 80}, map[string]interface{}{"Age": 12, "Weight": 80})
		if err == nil {
			t.Errorf("Error must not be nil. expected: %s actual: %v", "different values type", err)
		}
	})
}
