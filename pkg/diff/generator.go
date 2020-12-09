package diff

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

func GenerateNewDiff(ctx context.Context, obj reflect.Value) Diff {
	return Diff{
		ChangeType: New,
		ObjectType: obj.Type().String(),
		New:        obj.Interface(),
	}
}

func GenerateRemovedDiff(ctx context.Context, obj reflect.Value) Diff {
	return Diff{
		ChangeType: Removed,
		ObjectType: obj.Type().String(),
		Old:        obj.Interface(),
	}
}

func GenerateChangedDiff(ctx context.Context, fieldName string, oldVal, newVal reflect.Value) *Diff {
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

func reflectArrayToString(ctx context.Context, value reflect.Value) string {
	var result string
	for i := 0; i < value.Len(); i++ {
		result += fmt.Sprintf("%v,", value.Index(i).Interface())
	}

	return strings.TrimSuffix(result, ",")
}
