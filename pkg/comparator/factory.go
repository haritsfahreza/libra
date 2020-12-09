package comparator

import "reflect"

func GetComparator(kind reflect.Kind) Comparator {
	switch kind {
	case reflect.Struct:
		return &StructComparator{}
	case reflect.Map:
		return &MapComparator{}
	case reflect.Ptr:
		return &PointerComparator{}
	case reflect.Func:
		return &FunctionComparator{}
	default:
		return &GenericComparator{}
	}
}
