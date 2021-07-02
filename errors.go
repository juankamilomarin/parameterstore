package parameterstore

import (
	"errors"
	"fmt"
)

var ErrParamsGroupInvalidType = errors.New("params is not a pointer which value is a struct")

type ErrTagNotSetOrEmpty struct {
	// Satisfies the generic error interface.
	error
	fieldName string
}

// Satisfies the error interface.
func (e ErrTagNotSetOrEmpty) Error() string {
	return fmt.Sprintf(`tag not set or empty for field %s`, e.fieldName)
}

func newErrTagNotSetOrEmpty(fieldName string) *ErrTagNotSetOrEmpty {
	return &ErrTagNotSetOrEmpty{fieldName: fieldName}
}

type ErrParsingParameter struct {
	// Satisfies the generic error interface.
	error
	fieldName string
	tagName   string
	typeDesc  string
}

// Satisfies the error interface.
func (e ErrParsingParameter) Error() string {
	return fmt.Sprintf(`parsing parameter %s to field %s of type %s`, e.tagName, e.fieldName, e.typeDesc)
}

func newErrParsingParameter(fieldName string, tagName string, typeDesc string) *ErrParsingParameter {
	return &ErrParsingParameter{fieldName: fieldName, tagName: tagName, typeDesc: typeDesc}
}

type ErrTypeNotSupported struct {
	// Satisfies the generic error interface.
	error
	typeDesc  string
	fieldName string
}

// Satisfies the error interface.
func (e ErrTypeNotSupported) Error() string {
	return fmt.Sprintf(`type %s for field %s not supported`, e.typeDesc, e.fieldName)
}

func newErrTypeNotSupported(fieldName string, typeDesc string) *ErrTypeNotSupported {
	return &ErrTypeNotSupported{fieldName: fieldName, typeDesc: typeDesc}
}

type ErrFieldCannotBeSet struct {
	// Satisfies the generic error interface.
	error
	fieldName string
}

// Satisfies the error interface.
func (e ErrFieldCannotBeSet) Error() string {
	return fmt.Sprintf(`field %s cannot be set`, e.fieldName)
}

func newErrFieldCannotBeSet(fieldName string) *ErrFieldCannotBeSet {
	return &ErrFieldCannotBeSet{fieldName: fieldName}
}

type ErrParameterStoreFailure struct {
	// Satisfies the generic error interface.
	error
}

// Satisfies the error interface.
func (e ErrParameterStoreFailure) Error() string {
	return fmt.Sprintf(`error executing ParameterStore.GetParams: %s`, e.error.Error())
}

func newErrParameterStoreFailure(e error) *ErrParameterStoreFailure {
	return &ErrParameterStoreFailure{error: e}
}
