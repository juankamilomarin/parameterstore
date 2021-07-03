package parameterstore

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const tag = "psTag"
const errorTag = "psErrorTag"
const privateTag = "psPrivateTag"
const intParsingErrorTag = "psIntParsingErrorTag"
const uintParsingErrorTag = "psUintParsingErrorTag"
const floatParsingErrorTag = "psFloatParsingErrorTag"
const boolParsingErrorTag = "psBoolParsingErrorTag"
const notSupportedTypeTag = "psNotSupportedTypeTag"

type paramsGroup struct {
	IntParam               int               `psTag:"int_param"`
	IntParamBase2          int               `psTag:"int_param_base2"`
	IntParamBase8          int               `psTag:"int_param_base8"`
	IntParamBase16         int               `psTag:"int_param_base16"`
	Int8Param              int8              `psTag:"int8_param"`
	Int16Param             int16             `psTag:"int16_param"`
	Int32Param             int32             `psTag:"int32_param"`
	Int64Param             int64             `psTag:"int64_param"`
	UintParam              uint              `psTag:"uint_param"`
	Uint8Param             uint8             `psTag:"uint8_param"`
	Uint16Param            uint16            `psTag:"uint16_param"`
	Uint32Param            uint32            `psTag:"uint32_param"`
	Uint64Param            uint64            `psTag:"uint64_param"`
	Float32Param           float32           `psTag:"float32_param"`
	Float64Param           float64           `psTag:"float64_param"`
	StringParam            string            `psTag:"string_param"`
	BoolParam              bool              `psTag:"bool_param"`
	ErrorParam             int               `psErrorTag:"error"`
	priveParam             int               `psPrivateTag:"private_param"`
	IntParsingErrorParam   int               `psIntParsingErrorTag:"int_param_parsing_error"`
	UintParsingErrorParam  uint              `psUintParsingErrorTag:"uint_param_parsing_error"`
	FloatParsingErrorParam float32           `psFloatParsingErrorTag:"float_param_parsing_error"`
	BoolParsingErrorParam  bool              `psBoolParsingErrorTag:"bool_param_parsing_error"`
	MapParam               map[string]string `psNotSupportedTypeTag:"type_not_supported"`
}

type parameterStore struct{}

func (ps parameterStore) GetParams(m map[string]string) error {
	p := map[string]string{
		"int_param":                 "-2147483648",
		"int_param_base2":           "-0b10000000000000000000000000000000",
		"int_param_base8":           "-0o20000000000",
		"int_param_base16":          "-0x80000000",
		"int8_param":                "-127",
		"int16_param":               "-32768",
		"int32_param":               "-2147483648",
		"int64_param":               "-9223372036854775808",
		"uint_param":                "18446744073709551615",
		"uint8_param":               "255",
		"uint16_param":              "65535",
		"uint32_param":              "4294967295",
		"uint64_param":              "18446744073709551615",
		"float32_param":             "0.123456789121212",
		"float64_param":             "0.123456789121212121212",
		"string_param":              "hello world!",
		"bool_param":                "true",
		"int_param_parsing_error":   "error",
		"uint_param_parsing_error":  "error",
		"float_param_parsing_error": "error",
		"bool_param_parsing_error":  "error",
	}

	for key := range m {
		if key == "error" {
			return errors.New("cannot get parameters")
		}
		m[key] = p[key]
	}
	return nil
}

func TestParameterStore(t *testing.T) {
	ps := parameterStore{}
	pg := paramsGroup{priveParam: 0}
	err := LoadParamsGroup(&pg, ps, tag)
	assert.Nil(t, err)
	assert.Equal(t, pg.IntParam, -2147483648)
	assert.Equal(t, pg.IntParamBase2, -0b10000000000000000000000000000000)
	assert.Equal(t, pg.IntParamBase8, -0o20000000000)
	assert.Equal(t, pg.IntParamBase16, -0x80000000)
	assert.Equal(t, pg.Int8Param, int8(-127))
	assert.Equal(t, pg.Int16Param, int16(-32768))
	assert.Equal(t, pg.Int32Param, int32(-2147483648))
	assert.Equal(t, pg.Int64Param, int64(-9223372036854775808))
	assert.Equal(t, pg.UintParam, uint(18446744073709551615))
	assert.Equal(t, pg.Uint8Param, uint8(255))
	assert.Equal(t, pg.Uint16Param, uint16(65535))
	assert.Equal(t, pg.Uint32Param, uint32(4294967295))
	assert.Equal(t, pg.Uint64Param, uint64(18446744073709551615))
	assert.Equal(t, pg.Float32Param, float32(0.123456789121212))
	assert.Equal(t, pg.Float64Param, float64(0.123456789121212121212))
	assert.Equal(t, pg.StringParam, "hello world!")
	assert.Equal(t, pg.BoolParam, true)

	var strPointer *string
	str := "hello world!"
	strPointer = &str
	err = LoadParamsGroup(strPointer, ps, tag)
	assert.Error(t, err)
	assert.Equal(t, ErrParamsGroupInvalidType.Error(), err.Error())

	err = LoadParamsGroup(*strPointer, ps, tag)
	assert.Error(t, err)
	assert.Equal(t, ErrParamsGroupInvalidType.Error(), err.Error())

	err = LoadParamsGroup(&pg, ps, errorTag)
	assert.Error(t, err)
	assert.Equal(t, "error executing ParameterStore.GetParams: cannot get parameters", err.Error())

	err = LoadParamsGroup(&pg, ps, privateTag)
	assert.Error(t, err)
	assert.Equal(t, "field priveParam cannot be set", err.Error())

	err = LoadParamsGroup(&pg, ps, intParsingErrorTag)
	assert.Error(t, err)
	assert.Equal(t, "cannot parse parameter int_param_parsing_error to field IntParsingErrorParam of type int", err.Error())

	err = LoadParamsGroup(&pg, ps, uintParsingErrorTag)
	assert.Error(t, err)
	assert.Equal(t, "cannot parse parameter uint_param_parsing_error to field UintParsingErrorParam of type uint", err.Error())

	err = LoadParamsGroup(&pg, ps, floatParsingErrorTag)
	assert.Error(t, err)
	assert.Equal(t, "cannot parse parameter float_param_parsing_error to field FloatParsingErrorParam of type float32", err.Error())

	err = LoadParamsGroup(&pg, ps, boolParsingErrorTag)
	assert.Error(t, err)
	assert.Equal(t, "cannot parse parameter bool_param_parsing_error to field BoolParsingErrorParam of type bool", err.Error())

	err = LoadParamsGroup(&pg, ps, notSupportedTypeTag)
	assert.Error(t, err)
	assert.Equal(t, "type map for field MapParam not supported", err.Error())
}
