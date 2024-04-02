// Copyright 2024 Bill Nixon. All rights reserved. Use of this source code
// is governed by the license found in the LICENSE file.

// Package required offers utilities to check for required fields in a struct.
// A "required" struct tag identifies which fields should be non-zero. It
// provides a straightforward way to programmatically ensure that necessary
// data fields are populated, enhancing reliability and maintainability of
// code dealing with complex data structures.
package required

import (
	"errors"
	"reflect"
)

// ErrNotStructOrPtr indicates that the input argument is neither a struct
// nor a pointer to a struct.
var ErrNotStructOrPtr = errors.New("not a struct or pointer to a struct")

// isRequiredAndZero checks if the field is required and zero.
func isRequiredAndZero(field reflect.StructField, fieldValue reflect.Value) bool {
	requiredTag, hasRequired := field.Tag.Lookup("required")
	return hasRequired && requiredTag == "true" && fieldValue.IsZero()
}

// checkNested determines if the field value is a struct or a non-zero
// pointer that should be recursively checked for required unset fields.
func checkNested(fieldValue reflect.Value) bool {
	return fieldValue.Kind() == reflect.Struct || (fieldValue.Kind() == reflect.Ptr && !fieldValue.IsZero())
}
