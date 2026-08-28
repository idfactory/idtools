// Package idconfig loads application configuration from environment variables.
package idconfig

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

const (
	envTag      = "env"
	defaultTag  = "default"
	requiredTag = "required"
)

var durationType = reflect.TypeOf(time.Duration(0))

// Load populates the tagged fields of target from environment variables.
//
// Target must be a non-nil pointer to a struct. Supported field types are
// string, bool, int, int64, and time.Duration. Fields without an env tag are
// left unchanged.
//
// The default tag is used only when the environment variable is absent. A
// field tagged required:"true" must have a non-empty environment value.
// When any field is invalid, Load returns all field errors and leaves target
// unchanged.
func Load(target any) error {
	value := reflect.ValueOf(target)
	if value.Kind() != reflect.Pointer || value.IsNil() || value.Elem().Kind() != reflect.Struct {
		return errors.New("idconfig: target must be a non-nil pointer to a struct")
	}

	working := reflect.New(value.Elem().Type()).Elem()
	working.Set(value.Elem())

	var loadErrors []error
	structType := working.Type()
	for i := 0; i < working.NumField(); i++ {
		fieldType := structType.Field(i)
		key, tagged := fieldType.Tag.Lookup(envTag)
		if !tagged || key == "-" {
			continue
		}
		if key == "" {
			loadErrors = append(loadErrors, fieldError(fieldType, key, "env tag cannot be empty"))
			continue
		}

		field := working.Field(i)
		if !field.CanSet() {
			loadErrors = append(loadErrors, fieldError(fieldType, key, "field must be exported"))
			continue
		}

		required, err := requiredValue(fieldType)
		if err != nil {
			loadErrors = append(loadErrors, fieldError(fieldType, key, err.Error()))
			continue
		}

		raw, exists := os.LookupEnv(key)
		if required && (!exists || raw == "") {
			loadErrors = append(loadErrors, fieldError(fieldType, key, "required environment variable is not set"))
			continue
		}
		if !exists {
			var hasDefault bool
			raw, hasDefault = fieldType.Tag.Lookup(defaultTag)
			if !hasDefault {
				continue
			}
		}

		if err := setField(field, raw); err != nil {
			loadErrors = append(loadErrors, fieldError(fieldType, key, err.Error()))
		}
	}

	if err := errors.Join(loadErrors...); err != nil {
		return err
	}
	value.Elem().Set(working)
	return nil
}

func requiredValue(field reflect.StructField) (bool, error) {
	raw, exists := field.Tag.Lookup(requiredTag)
	if !exists {
		return false, nil
	}
	required, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("invalid required tag %q", raw)
	}
	return required, nil
}

func setField(field reflect.Value, raw string) error {
	if field.Type() == durationType {
		value, err := time.ParseDuration(raw)
		if err != nil {
			return errors.New("must be a duration such as 30s or 5m")
		}
		field.SetInt(int64(value))
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
	case reflect.Bool:
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return errors.New("must be a boolean")
		}
		field.SetBool(value)
	case reflect.Int, reflect.Int64:
		value, err := strconv.ParseInt(raw, 10, field.Type().Bits())
		if err != nil {
			return errors.New("must be an integer")
		}
		field.SetInt(value)
	default:
		return fmt.Errorf("unsupported field type %s", field.Type())
	}
	return nil
}

func fieldError(field reflect.StructField, key, message string) error {
	if key == "" {
		return fmt.Errorf("idconfig: field %s: %s", field.Name, message)
	}
	return fmt.Errorf("idconfig: field %s (%s): %s", field.Name, key, message)
}
