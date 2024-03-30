// Copyright 2024 Bill Nixon. All rights reserved.
// Use of this source code is governed by the license found in the LICENSE file.

// Package required provides a function to check for required struct fields.
package required

import (
	"errors"
	"reflect"
)

// ErrNotStructOrPtr is used if input is not a struct or a pointer to a struct.
var ErrNotStructOrPtr = errors.New("not a struct or pointer to a struct")

// Check returns names of struct fields tagged as `required:"true"` that
// are non zero.
//
// It accepts any type and returns ErrNotStructOrPtr if the input is not a
// struct or a pointer to a struct.
//
// Otherwise, it returns a slice of strings, each representing the path to an
// unset required field. An empty slice indicates all required fields are set.
//
// This function is useful to valid if structs are initialized.
func Check(s any) ([]string, error) {
	if reflect.Indirect(reflect.ValueOf(s)).Kind() != reflect.Struct {
		return nil, ErrNotStructOrPtr
	}

	return check_internal(s, ""), nil
}

func check_internal(s any, parentPath string) []string {
	var missing []string

	val := reflect.Indirect(reflect.ValueOf(s))

	// Ignore non-structs.
	if val.Kind() != reflect.Struct {
		return missing
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		// Construct fieldPath, prefixing with parentPath if necessary.
		fieldPath := field.Name
		if parentPath != "" {
			fieldPath = parentPath + "." + field.Name
		}

		// Determine if the field is required.
		requiredTag, hasRequired := field.Tag.Lookup("required")
		if hasRequired && requiredTag == "true" && value.IsZero() {
			missing = append(missing, fieldPath)
			continue // Skip further checks since required.
		}

		// Recursively check nested struct or non-zero pointer
		if value.Kind() == reflect.Struct || (value.Kind() == reflect.Ptr && !value.IsZero()) {
			nestedMissing := check_internal(value.Interface(), fieldPath)
			missing = append(missing, nestedMissing...)
		}
	}

	return missing
}
