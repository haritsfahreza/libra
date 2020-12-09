package comparator_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/haritsfahreza/libra/pkg/comparator"
	"github.com/haritsfahreza/libra/pkg/diff"
)

func TestMapComparator_Compare(t *testing.T) {
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
			"succeed when compare the maps",
			args{
				ctx: nil,
				old: map[string]interface{}{"Age": 22, "Weight": 80},
				new: map[string]interface{}{"Age": 23, "Weight": 80},
			},
			[]diff.Diff{{
				ChangeType: diff.Changed,
				ObjectType: "map[string]interface {}",
				Field:      "Age",
				Old:        22,
				New:        23,
			}},
			false,
		}, {
			"succeed when compare maps with nested struct",
			args{
				ctx: nil,
				old: map[string]interface{}{"Age": 22, "Person": person{Name: "Rima"}},
				new: map[string]interface{}{"Age": 22, "Person": person{Name: "Reza"}},
			},
			[]diff.Diff{{
				ChangeType: diff.Changed,
				ObjectType: "map[string]interface {}",
				ObjectID:   "0",
				Field:      "Person.Name",
				Old:        "Rima",
				New:        "Reza",
			}},
			false,
		}, {
			"failed when compare the maps with different value type",
			args{
				ctx: nil,
				old: map[string]interface{}{"Age": "A", "Weight": 80},
				new: map[string]interface{}{"Age": 23, "Weight": 80},
			},
			nil,
			true,
		}, {
			"failed when compare the maps with different value in nested type",
			args{
				ctx: nil,
				old: map[string]interface{}{"Weight": 80, "Person": person{Interface: "A"}},
				new: map[string]interface{}{"Weight": 80, "Person": person{Interface: 1}},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &comparator.MapComparator{}
			got, err := c.Compare(tt.args.ctx, reflect.ValueOf(tt.args.old), reflect.ValueOf(tt.args.new))
			if (err != nil) != tt.wantErr {
				t.Errorf("MapComparator.Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapComparator.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
