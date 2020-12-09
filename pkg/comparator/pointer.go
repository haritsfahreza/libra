package comparator

import (
	"context"
	"reflect"

	"github.com/haritsfahreza/libra/pkg/diff"
)

type PointerComparator struct{}

var _ Comparator = (*PointerComparator)(nil)

func (c *PointerComparator) Compare(ctx context.Context, oldVal, newVal reflect.Value) ([]diff.Diff, error) {
	oldPointerValue := reflect.ValueOf(oldVal.Elem().Interface())
	newPointerValue := reflect.ValueOf(newVal.Elem().Interface())
	comparator := GetComparator(oldPointerValue.Kind())

	return comparator.Compare(ctx, oldPointerValue, newPointerValue)
}
