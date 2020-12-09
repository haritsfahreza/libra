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

		tag := typeField.Tag.Get("libra")
		if tag == "ignore" {
			continue
		}

		if err := Validate(ctx, oldField, newField); err != nil {
			return nil, fmt.Errorf("Error on validate key %s Error : %s", typeField.Name, err.Error())
		}

		filteredOldValue := getInterfaceValue(oldField)
		filteredNewValue := getInterfaceValue(newField)
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

		if diff := diff.GenerateChangedDiff(ctx, typeField.Name, oldField, newField); diff != nil {
			diffs = append(diffs, *diff)
		}
	}

	objectID, err := getObjectID(ctx, oldVal)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(diffs); i++ {
		diffs[i].ObjectType = objectType
		diffs[i].ObjectID = objectID
	}

	return diffs, nil
}

func getObjectID(ctx context.Context, v reflect.Value) (string, error) {
	objectID := ""
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			objectIDInField, err := getObjectID(ctx, field)
			if err != nil {
				return "", err
			}

			if objectID != "" && objectIDInField != "" {
				return "", fmt.Errorf("tag `id` should defined once")
			}

			if objectID == "" && objectIDInField != "" {
				objectID = objectIDInField
			}
			continue
		}

		tag := typeField.Tag.Get("libra")
		if tag == "id" {
			if objectID != "" {
				return "", fmt.Errorf("tag `id` should defined once")
			}
			objectID = fmt.Sprintf("%v", field.Interface())
		}
	}

	return objectID, nil
}
