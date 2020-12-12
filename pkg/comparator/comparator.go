package comparator

import (
	"context"
	"fmt"
	"reflect"

	"github.com/haritsfahreza/libra/pkg/diff"
)

type Comparator interface {
	Compare(ctx context.Context, oldVal, newVal reflect.Value) ([]diff.Diff, error)
}

func Validate(ctx context.Context, oldVal, newVal reflect.Value) error {
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

func compareNestedField(ctx context.Context, baseFieldName string, oldField, newField reflect.Value) ([]diff.Diff, error) {
	comparator := GetComparator(oldField.Kind())
	nestedDiffs, err := comparator.Compare(ctx, oldField, newField)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(nestedDiffs); i++ {
		nestedDiffs[i].Field = fmt.Sprintf("%s.%s", baseFieldName, nestedDiffs[i].Field)
	}

	return nestedDiffs, nil
}

func isNestedKind(kind reflect.Kind) bool {
	return kind == reflect.Struct ||
		kind == reflect.Map ||
		kind == reflect.Ptr ||
		kind == reflect.Func
}

func filterValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		return reflect.ValueOf(v.Interface())
	}

	if m := v.MethodByName("String"); m.IsValid() {
		return m.Call([]reflect.Value{})[0]
	}

	return v
}
