package comparator

import (
	"context"
	"fmt"
	"reflect"

	"github.com/haritsfahreza/libra/pkg/diff"
)

type FunctionComparator struct{}

var _ Comparator = (*FunctionComparator)(nil)

func (c *FunctionComparator) Compare(ctx context.Context, oldVal, newVal reflect.Value) ([]diff.Diff, error) {
	return nil, fmt.Errorf("Unsupported comparable values")
}
