package parameterstore

import (
	"reflect"
	"strconv"
)

// The ParameterStore is the interface for
// representing a parameter store.
type ParameterStore interface {
	// Receives a slice of parameter names to lookup and returns a key-value map
	// of those parameters
	GetParams([]string) (map[string]string, error)
}

// LoadParamsGroup loads each field of the struct paramsGroup with the value of
// its matching parameter.
// A parameter name is given by the value of the tag for each field in the
// paramsGroup struct and its value is fetched from the parameterStore.
// If the tag is not set on a given field in paramsGroup, the loading process
// is ignored for that field.
// LoadParamsGroup returns an error in each of these cases:
//
// * The paramsGroup argument is not a struct.
//
// * A field in paramsGroup cannot be set. This is checked using
// CanSet function from reflection package.
//
// * A field type in paramsGroup is not supported.
//
// * A field in paramsGroup does not have the parameter tag or its value is empty.
//
// * A parameter value cannot be parsed to is matching field value.
//
// * The parameter store failed when getting the parameters.
func LoadParamsGroup(paramsGroup interface{}, ps ParameterStore, tag string) error {
	err := checkParamsGroupType(paramsGroup)
	if err != nil {
		return err
	}
	paramNames := getParamNames(paramsGroup, tag)
	params, err := ps.GetParams(paramNames)
	if err != nil {
		return newErrParameterStoreFailure(err)
	}
	err = parseParams(params, paramsGroup, tag)
	if err != nil {
		return err
	}
	return nil
}

// checkParamsGroupType checks that the type of paramsGroup is a pointer and
// its value type is a struct.
// checkParamsGroupType returns an error in case paramsGroup type
// is not valid.
func checkParamsGroupType(paramsGroup interface{}) error {
	if reflect.TypeOf(paramsGroup).Kind() != reflect.Ptr {
		return ErrParamsGroupInvalidType
	} else {
		if reflect.ValueOf(paramsGroup).Elem().Kind() != reflect.Struct {
			return ErrParamsGroupInvalidType
		}
	}
	return nil
}

// getParamNames takes the tag from each field in the paramsGroup struct and adds
// the value of that key to a slice of keys which is returned.
func getParamNames(paramsGroup interface{}, tag string) []string {
	params := []string{}
	pGroupFields := reflect.ValueOf(paramsGroup).Elem()
	pGroupType := pGroupFields.Type()
	for i := 0; i < pGroupFields.NumField(); i++ {
		field := pGroupType.Field(i)
		paramName := field.Tag.Get(tag)
		if paramName != "" {
			params = append(params, paramName)
		}
	}
	return params
}

// parseParams takes each element in params and parses its value
// to the type of the corresponding field in paramsGroup, using
// the key on params and the tag in the paramsGroup field to
// match such correspondence. ParseParams returns errors in
// each of these cases:
//
// * The parameter value cannot be parsed to the field value.
//
// * A field cannot be set. This is checked using CanSet function
// from reflection package.
//
// * A field type is not supported
func parseParams(params map[string]string, paramsGroup interface{}, tag string) error {
	pGroupFields := reflect.ValueOf(paramsGroup).Elem()
	pGroupType := pGroupFields.Type()
	for i := 0; i < pGroupFields.NumField(); i++ {
		field := pGroupType.Field(i)
		fieldValue := pGroupFields.Field(i)
		fieldName := field.Name
		fieldKind := field.Type.Kind()
		paramName := field.Tag.Get(tag)
		if paramName != "" {
			if !fieldValue.CanSet() {
				return newErrFieldCannotBeSet(fieldName)
			}
			switch fieldKind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				parsedInt, err := strconv.ParseInt(params[paramName], 0, 0)
				if err != nil {
					return newErrParsingParameter(fieldName, paramName, fieldKind.String())
				}
				fieldValue.SetInt(parsedInt)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				parsedUInt, err := strconv.ParseUint(params[paramName], 0, 0)
				if err != nil {
					return newErrParsingParameter(fieldName, paramName, fieldKind.String())
				}
				fieldValue.SetUint(parsedUInt)
			case reflect.Float32, reflect.Float64:
				parsedFloat, err := strconv.ParseFloat(params[paramName], 64)
				if err != nil {
					return newErrParsingParameter(fieldName, paramName, fieldKind.String())
				}
				fieldValue.SetFloat(parsedFloat)
			case reflect.String:
				fieldValue.SetString(params[paramName])
			case reflect.Bool:
				parsedBool, err := strconv.ParseBool(params[paramName])
				if err != nil {
					return newErrParsingParameter(fieldName, paramName, fieldKind.String())
				}
				fieldValue.SetBool(parsedBool)
			default:
				return newErrTypeNotSupported(fieldName, fieldKind.String())
			}
		}
	}
	return nil
}
