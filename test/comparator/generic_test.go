package comparator_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/haritsfahreza/libra/pkg/comparator"
	"github.com/haritsfahreza/libra/pkg/diff"
)

func TestGenericComparator_Compare(t *testing.T) {
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
			"succeed compare basic types",
			args{
				ctx: nil,
				old: "foo",
				new: "bar",
			},
			[]diff.Diff{{
				ChangeType: diff.Changed,
				ObjectType: "string",
				Old:        "foo",
				New:        "bar",
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &comparator.GenericComparator{}
			got, err := c.Compare(tt.args.ctx, reflect.ValueOf(tt.args.old), reflect.ValueOf(tt.args.new))
			if (err != nil) != tt.wantErr {
				t.Errorf("GenericComparator.Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenericComparator.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
