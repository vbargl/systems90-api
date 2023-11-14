package urlparams

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func Marshal(val any) (string, error) {
	var (
		values  = make(url.Values)
		reflval reflect.Value
		errs    []error
	)

	reflval, err := getStruct(reflect.ValueOf(val))
	if err != nil {
		return "", err
	}

	for i := 0; i < reflval.NumField(); i++ {
		field := reflval.Type().Field(i)
		tag := strings.TrimSpace(field.Tag.Get("urlparams"))
		if tag == "" {
			continue
		}

		name, flag, err := parseTag(tag)
		if err != nil {
			errs = append(errs, fmt.Errorf("urlparams: type %s, field %s: %w", reflval.Type().Name(), field.Name, err))
			continue
		}

		value, err := getValue(reflval.Field(i))
		if err != nil {
			errs = append(errs, fmt.Errorf("urlparams: type %s, field %s: %w", reflval.Type().Name(), field.Name, err))
			continue
		}

		if flag&OmitEmptyFlag != 0 && value == "" {
			continue
		}

		values.Add(name, value)
	}

	if len(errs) != 0 {
		return "", errors.Join(errs...)
	}

	return values.Encode(), nil
}

func getStruct(val reflect.Value) (reflect.Value, error) {
	switch val.Kind() {
	case reflect.Struct:
		return val, nil
	case reflect.Ptr:
		return getStruct(val.Elem())
	default:
		return reflect.Value{}, fmt.Errorf("urlparams: %s not supported, only struct is supported", val.Kind())
	}
}

func getValue(reflval reflect.Value) (string, error) {
	switch reflval.Kind() {
	case reflect.String:
		return reflval.String(), nil
	case reflect.Bool:
		return fmt.Sprintf("%t", reflval.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", reflval.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", reflval.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", reflval.Float()), nil
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%f", reflval.Complex()), nil
	case reflect.Ptr:
		return getValue(reflval.Elem())
	default:
		return "", fmt.Errorf("unsupported type %s", reflval.Type().Name())
	}
}
