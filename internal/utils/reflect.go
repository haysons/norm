package utils

import (
	"errors"
	"reflect"
)

func SliceSetElem(destValue reflect.Value, size int, callback func(i int, elem reflect.Value) (bool, error)) error {
	destType := destValue.Type()
	isArray := destType.Kind() == reflect.Array
	elemType := destType.Elem()
	isPtr := elemType.Kind() == reflect.Ptr
	if isPtr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() == reflect.Ptr {
		return errors.New("slice element cannot be a multilevel pointer")
	}
	if !isArray && destValue.Cap() == 0 {
		destValue.Set(reflect.MakeSlice(destValue.Type(), 0, size))
	}
	for i := 0; i < size; i++ {
		elem := reflect.New(elemType)
		set, err := callback(i, elem.Elem())
		if err != nil {
			return err
		}
		if !set {
			break
		}
		if !isPtr {
			elem = elem.Elem()
		}
		if isArray {
			if i <= destValue.Len()-1 {
				destValue.Index(i).Set(elem)
			}
		} else {
			destValue.Set(reflect.Append(destValue, elem))
		}
	}
	return nil
}

func PtrValue(destValue reflect.Value) reflect.Value {
	for destValue.Kind() == reflect.Ptr {
		if destValue.IsNil() && destValue.CanSet() {
			destValue.Set(reflect.New(destValue.Type().Elem()))
		}
		destValue = destValue.Elem()
	}
	return destValue
}

func StructFields(t reflect.Type) []reflect.StructField {
	return structFields(t, nil)
}

func structFields(t reflect.Type, parentIndex []int) []reflect.StructField {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fields := make([]reflect.StructField, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		field.Index = append(parentIndex, field.Index...)
		if field.Anonymous {
			fields = append(fields, structFields(field.Type, field.Index)...)
		} else {
			fields = append(fields, field)
		}
	}
	return fields
}
