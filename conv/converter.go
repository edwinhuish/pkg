package conv

import (
	"fmt"
	"reflect"
)

func StructToMap(structObj any, mapObj map[string]any) error {
	if mapObj == nil {
		return fmt.Errorf("mapObj must not be nil")
	}

	val := reflect.ValueOf(structObj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("structObj must be a struct or a pointer to a struct")
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name
		mapObj[fieldName] = field.Interface()
	}

	return nil
}

func MapToStruct(mapObj map[string]any, structObj any) error {
	val := reflect.ValueOf(structObj)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("structObj must be a non-nil pointer")
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("structObj must be a pointer to a struct")
	}

	for fieldName, fieldValue := range mapObj {
		structField := elem.FieldByName(fieldName)
		if !structField.IsValid() || !structField.CanSet() {
			continue
		}

		mapFieldValue := reflect.ValueOf(fieldValue)
		if mapFieldValue.Type().AssignableTo(structField.Type()) {
			structField.Set(mapFieldValue)
		} else if mapFieldValue.Type().ConvertibleTo(structField.Type()) {
			structField.Set(mapFieldValue.Convert(structField.Type()))
		}
	}

	return nil
}
