package parameterstore

import (
	"reflect"
	"strconv"
)

// The ParameterStore is the interface for
// representing a parameter store.
type ParameterStore interface {
	GetParams(map[string]string) error
}

// LoadParamsGroup sets each field of the struct paramsGroup to the value of
// the parameter. A parameter name is given by the value of the tag for each
// field in the paramsGroup and its value is fetched from the parameterStore.
// LoadParamsGroup returns an error in case there was one during the process.
//
// Returned Error Types:
//   * ErrParamsGroupInvalidType
//   The params argument is not a struct.
//
//   * ErrTagNotSetOrEmpty
//   A field does not have the parameter tag or its value is empty.
//
//   * ErrParameterStoreFailure
//   Parameter store failed when getting the parameters.
//
//   * ErrParsingParameter
//   Parsing error from parameter to its field.
//
//   * ErrFieldCannotBeSet
//   A field cannot be set. This is checked using CanSet function
//   from reflection package.
//
//	 * ErrTypeNotSupported
//	 The of a given field is not supported
//
func LoadParamsGroup(paramsGroup interface{}, ps ParameterStore, tag string) error {
	err := checkParamsGroupType(paramsGroup)
	if err != nil {
		return err
	}
	params := map[string]string{}
	err = setParamNames(params, paramsGroup, tag)
	if err != nil {
		return err
	}
	err = ps.GetParams(params)
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
// checkParamsGroupType returns ErrParamsGroupInvalidType in case paramsGroup type
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

// setParamNames sets the parameter name for each field in the paramsGroup struct
// to a key in params. Parameter names are taken from the given field tag.
// setParamNames returns ErrTagNotSetOrEmpty in case one of the fields
// does not have the parameter tag or it's value is empty.
func setParamNames(params map[string]string, paramsGroup interface{}, tag string) error {
	pGroupFields := reflect.ValueOf(paramsGroup).Elem()
	pGroupType := pGroupFields.Type()
	for i := 0; i < pGroupFields.NumField(); i++ {
		field := pGroupType.Field(i)
		paramName := field.Tag.Get(tag)
		if paramName == "" {
			return newErrTagNotSetOrEmpty(field.Name)
		}
		params[paramName] = ""
	}
	return nil
}

// parseParams takes each element in params and parses its value
// to the type of the corresponding field in paramsGroup, using
// the key on params and the tag in the paramsGroup field to
// match such correspondence.
//
// Returned Error Types:
//   * ErrParsingParameter
//   There was an error parsing the parameter to its field.
//
//   * ErrFieldCannotBeSet
//   A field cannot be set. This is checked using CanSet function
//   from reflection package.
//
//	 * ErrTypeNotSupported
//	 The of a given field is not supported
func parseParams(params map[string]string, paramsGroup interface{}, tag string) error {
	pGroupFields := reflect.ValueOf(paramsGroup).Elem()
	pGroupType := pGroupFields.Type()
	for i := 0; i < pGroupFields.NumField(); i++ {
		fieldType := pGroupType.Field(i).Type
		paramName := pGroupType.Field(i).Tag.Get(tag)
		if !pGroupFields.Field(i).CanSet() {
			return newErrFieldCannotBeSet(fieldType.Name())
		}
		switch k := fieldType.Kind(); k {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			parsedInt, err := strconv.ParseInt(params[paramName], 0, 0)
			if err != nil {
				return newErrParsingParameter(fieldType.Name(), paramName, k.String())
			}
			pGroupFields.Field(i).SetInt(parsedInt)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			parsedUInt, err := strconv.ParseUint(params[paramName], 0, 0)
			if err != nil {
				return newErrParsingParameter(fieldType.Name(), paramName, fieldType.Kind().String())
			}
			pGroupFields.Field(i).SetUint(parsedUInt)
		case reflect.Float32, reflect.Float64:
			parsedFloat, err := strconv.ParseFloat(params[paramName], 64)
			if err != nil {
				return newErrParsingParameter(fieldType.Name(), paramName, fieldType.Kind().String())
			}
			pGroupFields.Field(i).SetFloat(parsedFloat)
		case reflect.String:
			pGroupFields.Field(i).SetString(params[paramName])
		case reflect.Bool:
			parsedBool, err := strconv.ParseBool(params[paramName])
			if err != nil {
				return newErrParsingParameter(fieldType.Name(), paramName, fieldType.Kind().String())
			}
			pGroupFields.Field(i).SetBool(parsedBool)
		default:
			return newErrTypeNotSupported(fieldType.Name(), fieldType.Kind().String())
		}
	}
	return nil
}
