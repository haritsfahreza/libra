package libra

import (
	"context"
	"reflect"

	"github.com/haritsfahreza/libra/pkg/comparator"
	"github.com/haritsfahreza/libra/pkg/diff"
)

//Compare is used to compare two different values and spot the differences from them
func Compare(ctx context.Context, old, new interface{}) ([]diff.Diff, error) {
	oldVal := reflect.ValueOf(old)
	newVal := reflect.ValueOf(new)

	if err := comparator.Validate(ctx, oldVal, newVal); err != nil {
		return nil, err
	}

	if !oldVal.IsValid() && newVal.IsValid() {
		newDiff := diff.GenerateNewDiff(ctx, newVal)
		return []diff.Diff{newDiff}, nil
	}

	if oldVal.IsValid() && !newVal.IsValid() {
		oldDiff := diff.GenerateRemovedDiff(ctx, oldVal)
		return []diff.Diff{oldDiff}, nil
	}

	return comparator.GetComparator(oldVal.Kind()).Compare(ctx, oldVal, newVal)
}
