package diff

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

func GenerateNewDiff(ctx context.Context, obj reflect.Value) Diff {
	objectID := ""
	if objID, err := GetObjectID(ctx, obj); err == nil {
		objectID = objID
	}

	return Diff{
		ChangeType: New,
		ObjectType: obj.Type().String(),
		ObjectID:   objectID,
		New:        obj.Interface(),
	}
}

func GenerateRemovedDiff(ctx context.Context, obj reflect.Value) Diff {
	objectID := ""
	if objID, err := GetObjectID(ctx, obj); err == nil {
		objectID = objID
	}

	return Diff{
		ChangeType: Removed,
		ObjectType: obj.Type().String(),
		ObjectID:   objectID,
		Old:        obj.Interface(),
	}
}

func GenerateChangedDiff(ctx context.Context, fieldName string, oldVal, newVal reflect.Value) *Diff {
	var oldI, newI interface{}

	if !oldVal.IsValid() || !newVal.IsValid() {
		return nil
	}

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

func GetObjectID(ctx context.Context, v reflect.Value) (string, error) {
	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("ObjectID is only available for Struct")
	}

	objectID := ""
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			objectIDInField, err := GetObjectID(ctx, field)
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
