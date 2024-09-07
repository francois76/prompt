package reflection

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

func Assign(ref reflect.Value, fieldName string, input func(fieldStruct reflect.StructField) (string, error)) error {
	path := strings.Split(fieldName, ".")
	return assignSliceArrayStruct(ref, path[1:], input)
}

func assignSliceArrayStruct(ref reflect.Value, path []string, input func(fieldStruct reflect.StructField) (string, error)) error {
	if len(path) == 0 {
		return nil
	}
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	log.Debug().Any("path", path).Any("kind", ref.Kind().String()).Msg("assignSliceArrayStruct")

	nameAndPosition := strings.Split(path[0], "[")
	field := ref.FieldByName(nameAndPosition[0])

	if !field.IsValid() || !field.CanSet() {
		log.Warn().Msgf("field %s is not valid or cannot be set", nameAndPosition[0])
		return fmt.Errorf("field %s is not valid or cannot be set", nameAndPosition[0])
	}

	if len(nameAndPosition) > 1 { // It's a slice or array
		position, err := strconv.Atoi(strings.TrimRight(nameAndPosition[1], "]"))
		if err != nil {
			return fmt.Errorf("invalid index: %s", nameAndPosition[1])
		}

		if field.Kind() == reflect.Slice {
			if field.Len() <= position {
				newSlice := reflect.MakeSlice(field.Type(), position+1, position+1)
				reflect.Copy(newSlice, field)
				field.Set(newSlice)
			}
		} else if field.Kind() == reflect.Array {
			if field.Len() <= position {
				return fmt.Errorf("index out of bounds for array: %d", position)
			}
		} else {
			return fmt.Errorf("field %s is not a slice or array", nameAndPosition[0])
		}

		item := field.Index(position)
		if len(path) == 1 {
			return assignBase(item, "", input)
		}
		return assignSliceArrayStruct(item, path[1:], input)
	} else if field.Kind() == reflect.Struct {
		if len(path) == 1 {
			return fmt.Errorf("unabled to prompt a missing structure. Please remove the field validation on the field %s", nameAndPosition[0])
		}
		return assignSliceArrayStruct(field, path[1:], input)
	}
	return assignBase(ref, path[0], input)
}

func assignBase(ref reflect.Value, fieldName string, input func(fieldStruct reflect.StructField) (string, error)) error {
	log.Debug().Any("field", fieldName).Msg("assignBase")
	var field reflect.Value
	if fieldName == "" {
		field = ref
	} else {
		field = ref.FieldByName(fieldName)
	}

	if field.IsValid() && field.CanSet() {
		fieldStruct, _ := ref.Type().FieldByName(fieldName)
		value, err := input(fieldStruct)
		if err != nil {
			return err
		}
		val := castStringToType(value, field.Type())
		if val != nil {
			field.Set(reflect.ValueOf(val))
		}
	}
	return nil
}

func castStringToType(s string, t reflect.Type) any {
	switch t.Kind() {
	case reflect.String:
		return s
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil
		}
		return reflect.ValueOf(val).Convert(t).Interface()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil
		}
		return reflect.ValueOf(val).Convert(t).Interface()
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil
		}
		return reflect.ValueOf(val).Convert(t).Interface()
	case reflect.Bool:
		val, err := strconv.ParseBool(s)
		if err != nil {
			return nil
		}
		return val
	default:
		return nil
	}
}
