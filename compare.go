package libra

import (
	"context"
	"errors"
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

	diffs := []Diff{}

	if !oldVal.IsValid() && newVal.IsValid() {
		//New object
		diffs = append(diffs, Diff{
			ChangeType: New,
			ObjectType: newVal.Type().String(),
			New:        newVal.Interface(),
		})
		return diffs, nil
	} else if oldVal.IsValid() && !newVal.IsValid() {
		//Removed object
		diffs = append(diffs, Diff{
			ChangeType: Removed,
			ObjectType: oldVal.Type().String(),
			Old:        oldVal.Interface(),
		})
		return diffs, nil
	} else {
		objectID := ""
		switch oldVal.Kind() {
		case reflect.Struct:
			objectType := oldVal.Type().String()
			for i := 0; i < oldVal.NumField(); i++ {
				typeField := oldVal.Type().Field(i)
				oldField := oldVal.Field(i)
				newField := newVal.Field(i)

				tag := typeField.Tag.Get("libra")
				if tag == "ignore" {
					continue
				} else if tag == "id" {
					if objectID != "" {
						return nil, fmt.Errorf("tag `id` should defined once")
					}
					objectID = fmt.Sprintf("%v", oldField.Interface())
				}

				if err := validate(ctx, oldField, newField); err != nil {
					return nil, fmt.Errorf("Error on validate key %s Error : %s", typeField.Name, err.Error())
				}

				if diff := generateDiff(ctx, Changed, objectType, typeField.Name, oldField, newField); diff != nil {
					diffs = append(diffs, *diff)
				}
			}
		case reflect.Map:
			objectType := oldVal.Type().String()
			for _, key := range oldVal.MapKeys() {
				oldField := oldVal.MapIndex(key)
				newField := newVal.MapIndex(key)

				if err := validate(ctx, oldField, newField); err != nil {
					return nil, fmt.Errorf("Error on validate key %s Error : %s", key.String(), err.Error())
				}

				if diff := generateDiff(ctx, Changed, objectType, key.String(), oldField, newField); diff != nil {
					diffs = append(diffs, *diff)
				}
			}
		case reflect.Ptr:
			return Compare(ctx, oldVal.Elem().Interface(), newVal.Elem().Interface())
		case reflect.Func:
			return nil, fmt.Errorf("Unsupported comparable values")
		default:
			if diff := generateDiff(ctx, Changed, oldVal.Type().String(), "", oldVal, newVal); diff != nil {
				diffs = append(diffs, *diff)
			}
		}

		if objectID != "" {
			for i := 0; i < len(diffs); i++ {
				diffs[i].ObjectID = objectID
			}
		}

		return diffs, nil
	}
}

func generateDiff(ctx context.Context, changeType ChangeType, objectType, fieldName string, oldVal, newVal reflect.Value) *Diff {
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
			ObjectType: objectType,
			Field:      fieldName,
			Old:        oldI,
			New:        newI,
		}
	}

	return nil
}

func validate(ctx context.Context, oldVal, newVal reflect.Value) error {
	if !oldVal.IsValid() && !newVal.IsValid() {
		return errors.New("all values cannot be nil")
	}

	if oldVal.IsValid() && newVal.IsValid() {
		ov := reflect.ValueOf(oldVal.Interface())
		nv := reflect.ValueOf(newVal.Interface())

		if ov.IsValid() && nv.IsValid() && ov.Type() != nv.Type() {
			return errors.New("different values type")
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
