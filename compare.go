package libra

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

//Compare is used to compare two different values and spot the differences from them
func Compare(ctx context.Context, old, new interface{}) ([]Diff, error) {
	oldVal := reflect.ValueOf(old)
	newVal := reflect.ValueOf(new)

	if err := validate(ctx, oldVal, newVal); err != nil {
		return nil, err
	}

	if !oldVal.IsValid() && newVal.IsValid() {
		newDiff := generateNewDiff(ctx, newVal)
		return []Diff{newDiff}, nil
	}

	if oldVal.IsValid() && !newVal.IsValid() {
		oldDiff := generateRemovedDiff(ctx, oldVal)
		return []Diff{oldDiff}, nil
	}

	diffs := []Diff{}
	var objectID, objectType string
	switch oldVal.Kind() {
	case reflect.Struct:
		objectType = oldVal.Type().String()
		for i := 0; i < oldVal.NumField(); i++ {
			typeField := oldVal.Type().Field(i)
			oldField := oldVal.Field(i)
			newField := newVal.Field(i)

			tag := typeField.Tag.Get("libra")
			if tag == "ignore" {
				continue
			}

			if tag == "id" {
				if objectID != "" {
					return nil, fmt.Errorf("tag `id` should defined once")
				}
				objectID = fmt.Sprintf("%v", oldField.Interface())
			}

			if err := validate(ctx, oldField, newField); err != nil {
				return nil, fmt.Errorf("Error on validate key %s Error : %s", typeField.Name, err.Error())
			}

			nestedDiffs, err := compareBaseKind(ctx, typeField.Name, oldField, newField)
			if err != nil {
				return nil, err
			}

			if len(nestedDiffs) > 0 {
				if objectID == "" && nestedDiffs[0].ObjectID != "" {
					objectID = nestedDiffs[0].ObjectID
				}

				diffs = append(diffs, nestedDiffs...)
				continue
			}

			if diff := generateChangedDiff(ctx, typeField.Name, oldField, newField); diff != nil {
				diffs = append(diffs, *diff)
			}
		}
	case reflect.Map:
		objectType = oldVal.Type().String()
		for _, key := range oldVal.MapKeys() {
			oldField := oldVal.MapIndex(key)
			newField := newVal.MapIndex(key)

			if err := validate(ctx, oldField, newField); err != nil {
				return nil, fmt.Errorf("Error on validate key %s Error : %s", key.String(), err.Error())
			}

			nestedDiffs, err := compareBaseKind(ctx, key.String(), oldField, newField)
			if err != nil {
				return nil, err
			}

			if len(nestedDiffs) > 0 {
				diffs = append(diffs, nestedDiffs...)
				continue
			}

			if diff := generateChangedDiff(ctx, key.String(), oldField, newField); diff != nil {
				diffs = append(diffs, *diff)
			}
		}
	case reflect.Ptr:
		return Compare(ctx, oldVal.Elem().Interface(), newVal.Elem().Interface())
	case reflect.Func:
		return nil, fmt.Errorf("Unsupported comparable values")
	default:
		objectType = oldVal.Type().String()
		if diff := generateChangedDiff(ctx, "", oldVal, newVal); diff != nil {
			diffs = append(diffs, *diff)
		}
	}

	for i := 0; i < len(diffs); i++ {
		diffs[i].ObjectType = objectType
		diffs[i].ObjectID = objectID
	}

	return diffs, nil
}

func generateNewDiff(ctx context.Context, obj reflect.Value) Diff {
	return Diff{
		ChangeType: New,
		ObjectType: obj.Type().String(),
		New:        obj.Interface(),
	}
}

func generateRemovedDiff(ctx context.Context, obj reflect.Value) Diff {
	return Diff{
		ChangeType: Removed,
		ObjectType: obj.Type().String(),
		Old:        obj.Interface(),
	}
}

func generateChangedDiff(ctx context.Context, fieldName string, oldVal, newVal reflect.Value) *Diff {
	var oldI, newI interface{}
	switch oldVal.Kind() {
	case reflect.Array, reflect.Slice:
		oldI = reflectArrayToString(ctx, oldVal)
		newI = reflectArrayToString(ctx, newVal)
	default:
		oldI = oldVal.Interface()
		newI = newVal.Interface()
	}

	if oldI != newI {
		return &Diff{
			ChangeType: Changed,
			Field:      fieldName,
			Old:        oldI,
			New:        newI,
		}
	}

	return nil
}

func compareBaseKind(ctx context.Context, baseFieldName string, oldField, newField reflect.Value) ([]Diff, error) {
	var oldF, newF reflect.Value
	if oldField.Kind() == reflect.Interface {
		oldF = reflect.ValueOf(oldField.Interface())
		newF = reflect.ValueOf(newField.Interface())
	} else {
		oldF = oldField
		newF = newField
	}

	if oldF.Kind() == reflect.Struct ||
		oldF.Kind() == reflect.Map ||
		oldF.Kind() == reflect.Ptr ||
		oldF.Kind() == reflect.Func {
		nestedDiffs, err := Compare(ctx, oldF.Interface(), newF.Interface())
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(nestedDiffs); i++ {
			nestedDiffs[i].Field = fmt.Sprintf("%s.%s", baseFieldName, nestedDiffs[i].Field)
		}

		return nestedDiffs, nil
	}
	return []Diff{}, nil
}

func validate(ctx context.Context, oldVal, newVal reflect.Value) error {
	if !oldVal.IsValid() && !newVal.IsValid() {
		return fmt.Errorf("all values cannot be nil")
	}

	if oldVal.IsValid() && newVal.IsValid() {
		ov := reflect.ValueOf(oldVal.Interface())
		nv := reflect.ValueOf(newVal.Interface())

		if ov.IsValid() && nv.IsValid() && ov.Type() != nv.Type() {
			return fmt.Errorf("different values type")
		}
	}

	return nil
}

func reflectArrayToString(ctx context.Context, value reflect.Value) string {
	var result string
	for i := 0; i < value.Len(); i++ {
		result += fmt.Sprintf("%v,", value.Index(i).Interface())
	}
	return strings.TrimSuffix(result, ",")
}
