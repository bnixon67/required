// Copyright 2024 Bill Nixon. All rights reserved. Use of this source code
// is governed by the license found in the LICENSE file.

package required

import (
	"reflect"
)

// ArePresent checks if all struct fields tagged as `required:"true"` are
// non-zero.
//
// It accepts any type and returns false and ErrNotStructOrPtr if the input
// is not a struct or a pointer to a struct. Otherwise, it returns true if
// all required fields are set, or false if any are missing.
//
// This function simplifies struct initialization validation.
func ArePresent(s any) (bool, error) {
	v := reflect.ValueOf(s)

	if reflect.Indirect(v).Kind() != reflect.Struct {
		return false, ErrNotStructOrPtr
	}

	// Return the result of the recursive check.
	return presentInternal(v), nil
}

func presentInternal(v reflect.Value) bool {
	v = reflect.Indirect(v) // Dereference pointer, if any

	// Ignore non-structs.
	if v.Kind() != reflect.Struct {
		return true
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typ.Field(i)
		value := v.Field(i)

		// Check if field is required and missing
		requiredTag, hasRequired := field.Tag.Lookup("required")
		if hasRequired && requiredTag == "true" && value.IsZero() {
			return false // Found a required field that is missing.
		}

		// Recursively check nested struct or non-zero pointer
		if value.Kind() == reflect.Struct || (value.Kind() == reflect.Ptr && !value.IsZero()) {
			if !presentInternal(value) {
				return false
			}
		}
	}

	return true
}
