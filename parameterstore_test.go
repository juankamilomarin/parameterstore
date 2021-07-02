package parameterstore

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const tag = "psTag"

type paramsGroup struct {
	IntParam       int     `psTag:"int_param"`
	IntParamBase2  int     `psTag:"int_param_base2"`
	IntParamBase8  int     `psTag:"int_param_base8"`
	IntParamBase16 int     `psTag:"int_param_base16"`
	Int8Param      int8    `psTag:"int8_param"`
	Int16Param     int16   `psTag:"int16_param"`
	Int32Param     int32   `psTag:"int32_param"`
	Int64Param     int64   `psTag:"int64_param"`
	UintParam      uint    `psTag:"uint_param"`
	Uint8Param     uint8   `psTag:"uint8_param"`
	Uint16Param    uint16  `psTag:"uint16_param"`
	Uint32Param    uint32  `psTag:"uint32_param"`
	Uint64Param    uint64  `psTag:"uint64_param"`
	Float32Param   float32 `psTag:"float32_param"`
	Float64Param   float64 `psTag:"float64_param"`
	StringParam    string  `psTag:"string_param"`
	BoolParam      bool    `psTag:"bool_param"`
}

type paramsGroupWithError struct {
	IntParam int `psTag:"error"`
}

type paramsGroupWithTagEmpty struct {
	IntParam int `psTag:""`
}

type parameterStore struct{}

//base is implied by the string's prefix: 2 for "0b", 8 for "0" or "0o", 16 for "0x"
func (ps parameterStore) GetParams(m map[string]string) error {
	p := map[string]string{
		"int_param":        "-2147483648",
		"int_param_base2":  "-0b10000000000000000000000000000000",
		"int_param_base8":  "-0o20000000000",
		"int_param_base16": "-0x80000000",
		"int8_param":       "-127",
		"int16_param":      "-32768",
		"int32_param":      "-2147483648",
		"int64_param":      "-9223372036854775808",
		"uint_param":       "18446744073709551615",
		"uint8_param":      "255",
		"uint16_param":     "65535",
		"uint32_param":     "4294967295",
		"uint64_param":     "18446744073709551615",
		"float32_param":    "0.123456789121212",
		"float64_param":    "0.123456789121212121212",
		"string_param":     "hello world!",
		"bool_param":       "true",
	}

	for key := range m {
		if key == "error" {
			return errors.New("cannot get parameters")
		}
		m[key] = p[key]
	}
	return nil
}

func TestParameterStore_LoadParamsGroup(t *testing.T) {
	ps := parameterStore{}
	pg := paramsGroup{}
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

	pgError := paramsGroupWithError{}
	err = LoadParamsGroup(&pgError, ps, tag)
	assert.Error(t, err)
	assert.Equal(t, "error executing ParameterStore.GetParams: cannot get parameters", err.Error())
}

func TestParameterStore_checkParamsGroupType(t *testing.T) {
	var strPointer *string
	str := "hello world!"
	strPointer = &str
	err := checkParamsGroupType(strPointer)
	assert.Error(t, err)
	assert.Equal(t, ErrParamsGroupInvalidType.Error(), err.Error())
	err = checkParamsGroupType(*strPointer)
	assert.Error(t, err)
	assert.Equal(t, ErrParamsGroupInvalidType.Error(), err.Error())
}

func TestParameterStore_setParamNames(t *testing.T) {
	params := map[string]string{}
	pg := paramsGroup{}
	err := setParamNames(params, &pg, tag)
	assert.Nil(t, err)
	if _, found := params["int_param"]; !found {
		assert.Fail(t, "int_param not found")
	}
	if _, found := params["int_param_base2"]; !found {
		assert.Fail(t, "int_param_base2 not found")
	}
	if _, found := params["int_param_base8"]; !found {
		assert.Fail(t, "int_param_base8 not found")
	}
	if _, found := params["int_param_base16"]; !found {
		assert.Fail(t, "int_param_base16 not found")
	}
	if _, found := params["int8_param"]; !found {
		assert.Fail(t, "int8_param not found")
	}
	if _, found := params["int16_param"]; !found {
		assert.Fail(t, "int16_param not found")
	}
	if _, found := params["int32_param"]; !found {
		assert.Fail(t, "int32_param not found")
	}
	if _, found := params["int64_param"]; !found {
		assert.Fail(t, "int64_param not found")
	}
	if _, found := params["uint_param"]; !found {
		assert.Fail(t, "uint_param not found")
	}
	if _, found := params["uint8_param"]; !found {
		assert.Fail(t, "uint8_param not found")
	}
	if _, found := params["uint16_param"]; !found {
		assert.Fail(t, "uint16_param not found")
	}
	if _, found := params["uint32_param"]; !found {
		assert.Fail(t, "uint32_param not found")
	}
	if _, found := params["uint64_param"]; !found {
		assert.Fail(t, "uint64_param not found")
	}
	if _, found := params["float32_param"]; !found {
		assert.Fail(t, "float32_param not found")
	}
	if _, found := params["float64_param"]; !found {
		assert.Fail(t, "float64_param not found")
	}
	if _, found := params["string_param"]; !found {
		assert.Fail(t, "string_param not found")
	}
	if _, found := params["bool_param"]; !found {
		assert.Fail(t, "bool_param not found")
	}
	pgTagEmtpy := paramsGroupWithTagEmpty{}
	err = setParamNames(params, &pgTagEmtpy, tag)
	assert.Error(t, err)
	assert.Equal(t, "tag not set or empty for field IntParam", err.Error())
}
