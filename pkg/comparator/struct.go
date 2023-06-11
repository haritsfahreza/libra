package comparator

import (
	"context"
	"fmt"
	"reflect"

	"github.com/haritsfahreza/libra/pkg/diff"
)

type StructComparator struct{}

var _ Comparator = (*StructComparator)(nil)

func (c *StructComparator) Compare(ctx context.Context, oldVal, newVal reflect.Value) ([]diff.Diff, error) {
	diffs := []diff.Diff{}
	objectType := oldVal.Type().String()
	for i := 0; i < oldVal.NumField(); i++ {
		typeField := oldVal.Type().Field(i)
		oldField := oldVal.Field(i)
		newField := newVal.Field(i)

		if !typeField.IsExported() {
			continue
		}

		tag := typeField.Tag.Get("libra")
		if tag == "ignore" || tag == "id" {
			continue
		}

		if err := Validate(ctx, oldField, newField); err != nil {
			return nil, fmt.Errorf("error on validate key %s Error : %s", typeField.Name, err.Error())
		}

		filteredOldValue := filterValue(oldField)
		filteredNewValue := filterValue(newField)
		if isNestedKind(filteredOldValue.Kind()) {
			nestedDiffs, err := compareNestedField(
				ctx,
				typeField.Name,
				filteredOldValue,
				filteredNewValue,
			)
			if err != nil {
				return nil, err
			}

			diffs = append(diffs, nestedDiffs...)
			continue
		}

		if diff := diff.GenerateChangedDiff(ctx, typeField.Name, filteredOldValue, filteredNewValue); diff != nil {
			diffs = append(diffs, *diff)
		}
	}

	objectID, err := diff.GetObjectID(ctx, oldVal)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(diffs); i++ {
		diffs[i].ObjectType = objectType
		diffs[i].ObjectID = objectID
	}

	return diffs, nil
}
