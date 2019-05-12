package libra

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var defaultValue = reflect.Value{}

//Compare is used to compare two different values and spot the differences from them
func Compare(ctx context.Context, old, new interface{}) ([]Diff, error) {
	oldObj := reflect.ValueOf(old)
	newObj := reflect.ValueOf(new)

	if err := validate(ctx, oldObj, newObj); err != nil {
		return nil, err
	}

	diffs := []Diff{}

	if oldObj == defaultValue && newObj != defaultValue {
		//New object
		diffs = append(diffs, Diff{
			ChangeType: New,
			ObjectType: newObj.Type().String(),
			New:        newObj.Interface(),
		})
		return diffs, nil
	} else if oldObj != defaultValue && newObj == defaultValue {
		//Removed object
		diffs = append(diffs, Diff{
			ChangeType: Removed,
			ObjectType: oldObj.Type().String(),
			Old:        oldObj,
		})
		return diffs, nil
	} else {
		if oldObj.Kind() == reflect.Struct {
			objectType := oldObj.Type().String()
			for i := 0; i < oldObj.NumField(); i++ {
				typeField := oldObj.Type().Field(i)
				oldField := oldObj.Field(i)
				newField := newObj.Field(i)

				diff, err := generateDiff(ctx, Changed, objectType, typeField.Name, oldField, newField)
				if err != nil {
					return nil, err
				}

				if diff != nil {
					diffs = append(diffs, *diff)
				}
			}
		} else if oldObj.Kind() == reflect.Map {
			objectType := oldObj.Type().String()
			for _, key := range oldObj.MapKeys() {
				oldField := oldObj.MapIndex(key)
				newField := newObj.MapIndex(key)

				diff, err := generateDiff(ctx, Changed, objectType, key.String(), oldField, newField)
				if err != nil {
					return nil, err
				}

				if diff != nil {
					diffs = append(diffs, *diff)
				}
			}
		} else {
			return nil, fmt.Errorf("Unsupported comparable values")
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
	if oldValue == defaultValue && newValue == defaultValue {
		return errors.New("all values cannot be nil")
	}

	if oldValue != defaultValue && newValue != defaultValue {
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
