package comparator

import (
	"context"
	"reflect"

	"github.com/haritsfahreza/libra/pkg/diff"
)

type GenericComparator struct{}

var _ Comparator = (*GenericComparator)(nil)

func (c *GenericComparator) Compare(ctx context.Context, oldVal, newVal reflect.Value) ([]diff.Diff, error) {
	if changedDiff := diff.GenerateChangedDiff(ctx, "", oldVal, newVal); changedDiff != nil {
		changedDiff.ObjectType = oldVal.Type().String()

		return []diff.Diff{*changedDiff}, nil
	}

	return []diff.Diff{}, nil
}
