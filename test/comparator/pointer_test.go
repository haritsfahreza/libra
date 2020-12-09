package comparator_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/haritsfahreza/libra/pkg/comparator"
	"github.com/haritsfahreza/libra/pkg/diff"
)

func TestPointerComparator_Compare(t *testing.T) {
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
			"succeed when compare the pointer",
			args{
				ctx: nil,
				old: &person{
					ID:   10,
					Name: "test1",
				},
				new: &person{
					ID:   10,
					Name: "test2",
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &comparator.PointerComparator{}
			got, err := c.Compare(tt.args.ctx, reflect.ValueOf(tt.args.old), reflect.ValueOf(tt.args.new))
			if (err != nil) != tt.wantErr {
				t.Errorf("PointerComparator.Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PointerComparator.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
