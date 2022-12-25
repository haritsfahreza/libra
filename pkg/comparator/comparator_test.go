package comparator_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/haritsfahreza/libra/pkg/comparator"
)

func TestValidate(t *testing.T) {
	type args struct {
		ctx    context.Context
		oldVal reflect.Value
		newVal reflect.Value
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"return err when the all object is nil",
			args{
				ctx:    nil,
				oldVal: reflect.ValueOf(nil),
				newVal: reflect.ValueOf(nil),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := comparator.Validate(tt.args.ctx, tt.args.oldVal, tt.args.newVal); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
