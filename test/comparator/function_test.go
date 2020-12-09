package comparator_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/haritsfahreza/libra/pkg/comparator"
	"github.com/haritsfahreza/libra/pkg/diff"
)

func TestFunctionComparator_Compare(t *testing.T) {
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
			"failed when unsupported function type",
			args{
				ctx: nil,
				old: func() {},
				new: func() {},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &comparator.FunctionComparator{}
			got, err := c.Compare(tt.args.ctx, reflect.ValueOf(tt.args.old), reflect.ValueOf(tt.args.new))
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionComparator.Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionComparator.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
