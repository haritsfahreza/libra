package libra_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/haritsfahreza/libra"
	"github.com/haritsfahreza/libra/pkg/diff"
)

type person struct {
	ID        int `libra:"id"`
	Name      string
	Age       int
	Weight    float64
	IsMarried bool
	Hobbies   []string
	Numbers   []int
	Ignore    string `libra:"ignore"`
	Interface interface{}
}

type anotherPerson struct {
	ID      int `libra:"id"`
	IDAgain int `libra:"id"`
	Name    string
}

func TestCompare(t *testing.T) {
	type args struct {
		ctx context.Context
		old interface{}
		new interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []diff.Diff
		wantErr bool
	}{
		{
			"failed when all nil",
			args{
				ctx: nil,
				old: nil,
				new: nil,
			},
			nil,
			true,
		}, {
			"failed when different type - 1",
			args{
				ctx: nil,
				old: "",
				new: 1,
			},
			nil,
			true,
		}, {
			"failed when different type - 2",
			args{
				ctx: nil,
				old: person{
					Name: "test1",
				},
				new: anotherPerson{
					Name: "test1",
				},
			},
			nil,
			true,
		}, {
			"succeed when create new object",
			args{
				ctx: nil,
				old: nil,
				new: person{
					ID:   1,
					Name: "test1",
				},
			},
			[]diff.Diff{{
				ChangeType: diff.New,
				ObjectType: "libra_test.person",
				ObjectID:   "1",
				Old:        nil,
				New: person{
					ID:   1,
					Name: "test1",
				},
			}},
			false,
		}, {
			"succeed when remove an object",
			args{
				ctx: nil,
				old: person{
					ID:   1,
					Name: "test1",
				},
				new: nil,
			},
			[]diff.Diff{{
				ChangeType: diff.Removed,
				ObjectType: "libra_test.person",
				ObjectID:   "1",
				Old: person{
					ID:   1,
					Name: "test1",
				},
				New: nil,
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := libra.Compare(tt.args.ctx, tt.args.old, tt.args.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

var bDiffs []diff.Diff

func benchmarkCompare(old, new interface{}, b *testing.B) {
	var diffs []diff.Diff
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
