package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func SetReflectValueFromString(value reflect.Value, str string) error {
	if value.Kind() == reflect.Ptr {
		// If the value is a pointer, dereference it
		if value.IsNil() {
			// If the pointer is nil, create a new instance and set it
			value.Set(reflect.New(value.Type().Elem()))
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		value.SetString(str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, err := strconv.ParseInt(str, 10, 64); err == nil {
			value.SetInt(intValue)
		} else {
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintValue, err := strconv.ParseUint(str, 10, 64); err == nil {
			value.SetUint(uintValue)
		} else {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, err := strconv.ParseFloat(str, 64); err == nil {
			value.SetFloat(floatValue)
		} else {
			return err
		}
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(str); err == nil {
			value.SetBool(boolValue)
		} else {
			return err
		}
	case reflect.Struct, reflect.Array:
		switch value.Type() {
		case reflect.TypeOf(uuid.UUID{}):
			if uuidValue, err := uuid.Parse(str); err == nil {
				value.Set(reflect.ValueOf(uuidValue))
			} else {
				return err
			}
		case reflect.TypeOf(time.Time{}):
			if timeValue, err := time.Parse(time.RFC3339, str); err == nil {
				value.Set(reflect.ValueOf(timeValue))
			} else {
				return err
			}
		}
	default:
		fmt.Println(value.Kind().String())
		return errors.New("invalid data type")
	}

	return nil
}

func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
	}

	return false
}

func Value(value interface{}) interface{} {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		} else {
			return v.Elem().Interface()
		}
	}
	return value
}
