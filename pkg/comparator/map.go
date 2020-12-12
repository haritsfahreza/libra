package comparator

import (
	"context"
	"fmt"
	"reflect"

	"github.com/haritsfahreza/libra/pkg/diff"
)

type MapComparator struct{}

var _ Comparator = (*MapComparator)(nil)

func (c *MapComparator) Compare(ctx context.Context, oldVal, newVal reflect.Value) ([]diff.Diff, error) {
	diffs := []diff.Diff{}
	objectType := oldVal.Type().String()
	objectID := ""
	for _, key := range oldVal.MapKeys() {
		oldField := oldVal.MapIndex(key)
		newField := newVal.MapIndex(key)

		if err := Validate(ctx, oldField, newField); err != nil {
			return nil, fmt.Errorf("Error on validate key %s Error : %s", key.String(), err.Error())
		}

		filteredOldValue := filterValue(oldField)
		filteredNewValue := filterValue(newField)
		if isNestedKind(filteredOldValue.Kind()) {
			nestedDiffs, err := compareNestedField(
				ctx,
				key.String(),
				filteredOldValue,
				filteredNewValue,
			)
			if err != nil {
				return nil, err
			}

			if objectID == "" && len(nestedDiffs) > 0 && nestedDiffs[0].ObjectID != "" {
				objectID = nestedDiffs[0].ObjectID
			}

			diffs = append(diffs, nestedDiffs...)
			continue
		}

		if diff := diff.GenerateChangedDiff(ctx, key.String(), filteredOldValue, filteredNewValue); diff != nil {
			diffs = append(diffs, *diff)
		}
	}

	for i := 0; i < len(diffs); i++ {
		diffs[i].ObjectType = objectType
		diffs[i].ObjectID = objectID
	}

	return diffs, nil
}
