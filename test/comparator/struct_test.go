package comparator_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/haritsfahreza/libra/pkg/comparator"
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
	Address   address
}

type address struct {
	Street    string
	City      string
	Interface interface{}
}

type AddressToBeEmbedded struct {
	ID    int    `libra:"id"`
	City  string `libra:"ignore"`
	State string
}

type embeddedAddress struct {
	AddressToBeEmbedded
	Street string
}

type anotherPerson struct {
	ID      int `libra:"id"`
	IDAgain int `libra:"id"`
	Name    string
}

func TestStructComparator_Compare(t *testing.T) {
	oldEmbeddedAddress := embeddedAddress{}
	oldEmbeddedAddress.ID = 10
	oldEmbeddedAddress.Street = "Jalan 123"
	oldEmbeddedAddress.City = "Malang"
	oldEmbeddedAddress.State = "Jawa Timur"

	oldEmbeddedAddressWithoutChange := oldEmbeddedAddress
	oldEmbeddedAddressWithoutChange.State = ""

	newEmbeddedAddress := embeddedAddress{}
	newEmbeddedAddress.ID = 10
	newEmbeddedAddress.Street = "Jalan ABC"
	newEmbeddedAddress.City = "Ngalam"
	newEmbeddedAddress.State = "Jatim"

	newEmbeddedAddressWithoutChange := newEmbeddedAddress
	newEmbeddedAddressWithoutChange.State = ""

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
			"succeed when compare the structs",
			args{
				ctx: nil,
				old: person{
					ID:      10,
					Name:    "test1",
					Numbers: []int{1, 2, 3},
				},
				new: person{
					ID:      10,
					Name:    "test2",
					Numbers: []int{1, 2, 4},
				},
			},
			[]diff.Diff{{
				ChangeType: diff.Changed,
				ObjectType: "comparator_test.person",
				Field:      "Name",
				ObjectID:   "10",
				Old:        "test1",
				New:        "test2",
			}, {
				ChangeType: diff.Changed,
				ObjectType: "comparator_test.person",
				Field:      "Numbers",
				ObjectID:   "10",
				Old:        "1,2,3",
				New:        "1,2,4",
			}},
			false,
		}, {
			"failed when compare the structs with different value type",
			args{
				ctx: nil,
				old: person{
					ID:        10,
					Name:      "test1",
					Numbers:   []int{1, 2, 3},
					Interface: "A",
				},
				new: person{
					ID:        10,
					Name:      "test2",
					Numbers:   []int{1, 2, 4},
					Interface: 1,
				},
			},
			nil,
			true,
		}, {
			"success when ignore the field",
			args{
				ctx: nil,
				old: person{
					ID:     10,
					Name:   "test1",
					Ignore: "Should not compared",
				},
				new: person{
					ID:     10,
					Name:   "test2",
					Ignore: "Should not compared 2",
				},
			},
			[]diff.Diff{{
				ChangeType: diff.Changed,
				ObjectType: "comparator_test.person",
				Field:      "Name",
				ObjectID:   "10",
				Old:        "test1",
				New:        "test2",
			}},
			false,
		}, {
			"success when compare nested struct",
			args{
				ctx: nil,
				old: person{
					ID:   10,
					Name: "test1",
					Address: address{
						Street: "jalan 123",
					},
				},
				new: person{
					ID:   10,
					Name: "test1",
					Address: address{
						Street: "jalan ABC",
					},
				},
			},
			[]diff.Diff{{
				ChangeType: diff.Changed,
				ObjectType: "comparator_test.person",
				Field:      "Address.Street",
				ObjectID:   "10",
				Old:        "jalan 123",
				New:        "jalan ABC",
			}},
			false,
		}, {
			"failed when compare different field type in nested struct",
			args{
				ctx: nil,
				old: person{
					ID:   10,
					Name: "test1",
					Address: address{
						Interface: "A",
					},
				},
				new: person{
					ID:   10,
					Name: "test1",
					Address: address{
						Interface: 10,
					},
				},
			},
			nil,
			true,
		}, {
			"failed when the objects have multiple tag id",
			args{
				ctx: nil,
				old: anotherPerson{
					ID:   10,
					Name: "test1",
				},
				new: anotherPerson{
					ID:   10,
					Name: "test2",
				},
			},
			nil,
			true,
		}, {
			"success when compare embedded struct with any change inside of it",
			args{
				ctx: nil,
				old: oldEmbeddedAddress,
				new: newEmbeddedAddress,
			},
			[]diff.Diff{{
				ChangeType: diff.Changed,
				ObjectType: "comparator_test.embeddedAddress",
				Field:      "AddressToBeEmbedded.State",
				ObjectID:   "10",
				Old:        "Jawa Timur",
				New:        "Jatim",
			}, {
				ChangeType: diff.Changed,
				ObjectType: "comparator_test.embeddedAddress",
				Field:      "Street",
				ObjectID:   "10",
				Old:        "Jalan 123",
				New:        "Jalan ABC",
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &comparator.StructComparator{}
			got, err := c.Compare(tt.args.ctx, reflect.ValueOf(tt.args.old), reflect.ValueOf(tt.args.new))
			if (err != nil) != tt.wantErr {
				t.Errorf("StructComparator.Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StructComparator.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
