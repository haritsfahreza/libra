package tampah_test

import (
	"testing"

	"github.com/uwuh/tampah"
)

type testStruct struct {
	Name      string
	Age       int
	Weight    float64
	IsMarried bool
}

type anotherTestStruct struct {
	Name string
}

func TestCompare(t *testing.T) {
	t.Run("failed when all nil", func(t *testing.T) {
		if _, err := tampah.Compare(nil, nil, nil); err == nil {
			t.Errorf("Error must not be nil. expected: %s actual: %v", "all values cannot be nil", err)
		}
	})

	t.Run("failed when different type", func(t *testing.T) {
		if _, err := tampah.Compare(nil, "", 1); err == nil {
			t.Errorf("Error must not be nil. expected: %s actual: %v", "different values type", err)
		}

		if _, err := tampah.Compare(nil, testStruct{
			Name: "test1",
		}, anotherTestStruct{
			Name: "test1",
		}); err == nil {
			t.Errorf("Error must not be nil. expected: %s actual: %v", "different values type", err)
		}
	})

	t.Run("succeed when create new object", func(t *testing.T) {
		diffs, err := tampah.Compare(nil, nil, testStruct{
			Name: "test1",
			Age:  22,
		})
		if err != nil {
			t.Errorf("Error must be nil. Got: %v", err)
		}
		if len(diffs) != 1 {
			t.Errorf("Compare result length must be only one. expected: %d actual: %d", 1, len(diffs))
		}

		if diffs[0].ChangeType != tampah.New {
			t.Errorf("Invalid ChangeType. expected: %s actual: %s", tampah.New, diffs[0].ChangeType)
		}

		if diffs[0].ObjectType != "tampah_test.testStruct" {
			t.Errorf("Invalid ObjectType. expected: %s actual: %s", "tampah_test.testStruct", diffs[0].ObjectType)
		}

		if diffs[0].New == nil {
			t.Errorf("Invalid New value. expected: %s actual: %v", "not nil", diffs[0].New)
		}

		if diffs[0].Old != nil {
			t.Errorf("Invalid Old value. expected: %s actual: %v", "nil", diffs[0].Old)
		}
	})

	t.Run("succeed when removed an object", func(t *testing.T) {
		diffs, err := tampah.Compare(nil, testStruct{
			Name: "test1",
			Age:  22,
		}, nil)
		if err != nil {
			t.Errorf("Error must be nil. Got: %v", err)
		}
		if len(diffs) != 1 {
			t.Errorf("Compare result length must be only one. expected: %d actual: %d", 1, len(diffs))
		}

		if diffs[0].ChangeType != tampah.Removed {
			t.Errorf("Invalid ChangeType. expected: %s actual: %s", tampah.Removed, diffs[0].ChangeType)
		}

		if diffs[0].ObjectType != "tampah_test.testStruct" {
			t.Errorf("Invalid ObjectType. expected: %s actual: %s", "tampah_test.testStruct", diffs[0].ObjectType)
		}

		if diffs[0].New != nil {
			t.Errorf("Invalid New value. expected: %s actual: %v", "nil", diffs[0].New)
		}

		if diffs[0].Old == nil {
			t.Errorf("Invalid Old value. expected: %s actual: %v", "not nil", diffs[0].Old)
		}
	})

	t.Run("succeed when changed two object", func(t *testing.T) {
		diffs, err := tampah.Compare(nil, testStruct{
			Name:      "test1",
			Age:       22,
			Weight:    float64(80),
			IsMarried: true,
		}, testStruct{
			Name:      "test1",
			Age:       23,
			Weight:    float64(85),
			IsMarried: true,
		})
		if err != nil {
			t.Errorf("Error must be nil. Got: %v", err)
		}
		if len(diffs) != 2 {
			t.Errorf("Compare result length must be two. expected: %d actual: %d", 2, len(diffs))
		} else {
			expectedFields := []string{"Age", "Weight"}
			for i := 0; i < len(diffs); i++ {
				if diffs[i].ChangeType != tampah.Changed {
					t.Errorf("Invalid diffs[%d].ChangeType. expected: %s actual: %s", i, tampah.Changed, diffs[i].ChangeType)
				}

				if diffs[i].ObjectType != "tampah_test.testStruct" {
					t.Errorf("Invalid diffs[%d].ObjectType. expected: %s actual: %s", i, "tampah_test.testStruct", diffs[i].ObjectType)
				}

				if diffs[i].Field != expectedFields[i] {
					t.Errorf("Invalid diffs[%d].Field. expected: %s actual: %s", i, expectedFields[i], diffs[i].Field)
				}

				if diffs[i].Old == diffs[i].New {
					t.Errorf("diffs[i].Old must be different with diffs[%d].New. old: %s new: %v", i, diffs[i].Old, diffs[i].New)
				}
			}
		}
	})
}
