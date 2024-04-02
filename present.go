// Copyright 2024 Bill Nixon. All rights reserved. Use of this source code
// is governed by the license found in the LICENSE file.

package required

import (
	"reflect"
)

// ArePresent verifies all struct fields tagged as `required:"true"`
// are non-zero.  Returns true if all required fields are set, false and
// ErrNotStructOrPtr if the input is not a struct or its pointer.
func ArePresent(input any) (bool, error) {
	inputValue := reflect.Indirect(reflect.ValueOf(input))
	if inputValue.Kind() != reflect.Struct {
		return false, ErrNotStructOrPtr
	}

	return checkPresent(inputValue), nil
}

// checkPresent recursively checks if required fields are present.
func checkPresent(inputValue reflect.Value) bool {
	inputValue = reflect.Indirect(inputValue)

	// Ignore non-structs.
	if inputValue.Kind() != reflect.Struct {
		return true
	}

	typ := inputValue.Type()
	for i := 0; i < inputValue.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := inputValue.Field(i)

		if isRequiredAndZero(field, fieldValue) {
			return false
		}

		if checkNested(fieldValue) && !checkPresent(fieldValue) {
			return false
		}
	}

	return true
}
