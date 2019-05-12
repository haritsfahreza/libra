package libra

import (
	"context"
	"errors"
	"reflect"
)

var zeroValue = reflect.Value{}

//Compare is used to compare two different values and spot the differents from them
func Compare(ctx context.Context, old, new interface{}) ([]Diff, error) {
	oldValue := reflect.ValueOf(old)
	newValue := reflect.ValueOf(new)

	if err := validate(ctx, oldValue, newValue); err != nil {
		return nil, err
	}

	diff := []Diff{}

	if oldValue == zeroValue && newValue != zeroValue {
		//New object
		diff = append(diff, Diff{
			ChangeType: New,
			ObjectType: newValue.Type().String(),
			New:        newValue.Interface(),
		})
		return diff, nil
	} else if oldValue != zeroValue && newValue == zeroValue {
		//Removed object
		diff = append(diff, Diff{
			ChangeType: Removed,
			ObjectType: oldValue.Type().String(),
			Old:        oldValue,
		})
		return diff, nil
	} else {
		//Changed object
		objectType := oldValue.Type().String()
		for i := 0; i < oldValue.NumField(); i++ {
			oldFieldValue := oldValue.Field(i).Interface()
			newFieldValue := newValue.Field(i).Interface()
			typeField := oldValue.Type().Field(i)

			if oldFieldValue != newFieldValue {
				diff = append(diff, Diff{
					ChangeType: Changed,
					ObjectType: objectType,
					Field:      typeField.Name,
					Old:        oldFieldValue,
					New:        newFieldValue,
				})
			}
		}
		return diff, nil
	}
}

func validate(ctx context.Context, oldValue, newValue reflect.Value) error {
	if oldValue == zeroValue && newValue == zeroValue {
		return errors.New("all values cannot be nil")
	}

	if (oldValue != reflect.Value{} && newValue != reflect.Value{}) && oldValue.Type() != newValue.Type() {
		return errors.New("different values type")
	}
	return nil
}
