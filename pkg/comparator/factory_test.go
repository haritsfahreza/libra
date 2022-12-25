package comparator_test

import (
	"reflect"
	"testing"

	"github.com/haritsfahreza/libra/pkg/comparator"
)

func TestGetComparator(t *testing.T) {
	type args struct {
		kind reflect.Kind
	}
	tests := []struct {
		name string
		args args
		want comparator.Comparator
	}{
		{
			"return struct comparator when the object type is struct",
			args{
				kind: reflect.Struct,
			},
			&comparator.StructComparator{},
		},
		{
			"return map comparator when the object type is map",
			args{
				kind: reflect.Map,
			},
			&comparator.MapComparator{},
		},
		{
			"return pointer comparator when the object type is pointer",
			args{
				kind: reflect.Ptr,
			},
			&comparator.PointerComparator{},
		},
		{
			"return function comparator when the object type is function",
			args{
				kind: reflect.Func,
			},
			&comparator.FunctionComparator{},
		},
		{
			"return generic comparator when the object type is string",
			args{
				kind: reflect.String,
			},
			&comparator.GenericComparator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := comparator.GetComparator(tt.args.kind); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetComparator() = %v, want %v", got, tt.want)
			}
		})
	}
}
