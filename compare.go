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

				diff, err := generateDiff(ctx, Changed, objectType, typeField.Name, oldField, newField)
				if err != nil {
					return nil, err
				}

				if diff != nil {
					diffs = append(diffs, *diff)
				}
			}
			break

		case reflect.Map:
			objectType := oldVal.Type().String()
			for _, key := range oldVal.MapKeys() {
				oldField := oldVal.MapIndex(key)
				newField := newVal.MapIndex(key)

				diff, err := generateDiff(ctx, Changed, objectType, key.String(), oldField, newField)
				if err != nil {
					return nil, err
				}

				if diff != nil {
					diffs = append(diffs, *diff)
				}
			}
			break

		case reflect.Ptr:
			return Compare(ctx, oldVal.Elem().Interface(), newVal.Elem().Interface())
		default:
			return nil, fmt.Errorf("Unsupported comparable values")
		}

		if objectID != "" {
			for i := 0; i < len(diffs); i++ {
				diffs[i].ObjectID = objectID
			}
		}

		return diffs, nil
	}
}

func generateDiff(ctx context.Context, changeType ChangeType, objectType, fieldName string, oldField, newField reflect.Value) (*Diff, error) {
	if err := validate(ctx, oldField, newField); err != nil {
		return nil, fmt.Errorf("Error on validate key %s Error : %s", fieldName, err.Error())
	}

	var oldFieldValue, newFieldValue interface{}
	if reflect.ValueOf(oldField.Interface()).Kind() == reflect.Slice {
		oldFieldValue = reflectArrayToString(ctx, reflect.ValueOf(oldField.Interface()))
		newFieldValue = reflectArrayToString(ctx, reflect.ValueOf(newField.Interface()))
	} else {
		oldFieldValue = oldField.Interface()
		newFieldValue = newField.Interface()
	}

	if oldFieldValue != newFieldValue {
		return &Diff{
			ChangeType: Changed,
			ObjectType: objectType,
			Field:      fieldName,
			Old:        oldFieldValue,
			New:        newFieldValue,
		}, nil
	}

	return nil, nil
}

func validate(ctx context.Context, oldValue, newValue reflect.Value) error {
	if !oldValue.IsValid() && !newValue.IsValid() {
		return errors.New("all values cannot be nil")
	}

	if oldValue.IsValid() && newValue.IsValid() {
		oldValueType := reflect.ValueOf(oldValue.Interface()).Type()
		newValueType := reflect.ValueOf(newValue.Interface()).Type()

		if oldValueType != newValueType {
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
