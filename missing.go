// Copyright 2024 Bill Nixon. All rights reserved. Use of this source code
// is governed by the license found in the LICENSE file.

package required

import (
	"reflect"
)

// MissingFields returns names of struct fields tagged as `required:"true"`
// that are non zero.
//
// It accepts any type and returns ErrNotStructOrPtr if the input is not a
// struct or a pointer to a struct.
//
// Otherwise, it returns a slice of strings, each representing the path to an
// unset required field. An empty slice indicates all required fields are set.
//
// This function is useful to valid if structs are initialized.
func MissingFields(s any) ([]string, error) {
	v := reflect.ValueOf(s)

	if reflect.Indirect(v).Kind() != reflect.Struct {
		return nil, ErrNotStructOrPtr
	}

	// Return result of the recursive check
	return missingInternal(v, ""), nil
}

func missingInternal(v reflect.Value, parentPath string) []string {
	var missing []string

	v = reflect.Indirect(v)

	// Ignore non-structs.
	if v.Kind() != reflect.Struct {
		return missing
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typ.Field(i)
		value := v.Field(i)

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
			nestedMissing := missingInternal(value, fieldPath)
			missing = append(missing, nestedMissing...)
		}
	}

	return missing
}
