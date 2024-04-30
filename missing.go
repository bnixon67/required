// Copyright 2024 Bill Nixon. All rights reserved. Use of this source code
// is governed by the license found in the LICENSE file.

package required

import (
	"reflect"
)

// MissingFields finds unset required fields in a struct, identified by the
// `required:"true"` tag.  It returns an empty slice if all required fields
// are set or a slice of paths to unset required fields otherwise.
func MissingFields(input any) ([]string, error) {
	inputValue := reflect.Indirect(reflect.ValueOf(input))
	if inputValue.Kind() != reflect.Struct {
		return nil, ErrNotStructOrPtr
	}

	return findMissing(inputValue, "", true), nil
}

// buildFieldPath constructs the path to a field, prefixed by its parent
// path if present. For top-level fields, only the field name is used.
func buildFieldPath(parentPath, fieldName string) string {
	if parentPath == "" {
		return fieldName
	}
	return parentPath + "." + fieldName
}

// findMissing recursively identifies unset fields tagged as required,
// returning their paths.
func findMissing(inputValue reflect.Value, parentPath string, isDirectEmbedded bool) []string {
	var missingFields []string
	inputValue = reflect.Indirect(inputValue)

	// Ignore non-structs.
	if inputValue.Kind() != reflect.Struct {
		return missingFields
	}

	typ := inputValue.Type()
	for i := 0; i < inputValue.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := inputValue.Field(i)

		// Check if the current field is embedded
		isEmbedded := field.Anonymous && isDirectEmbedded

		// Calculate the current field path
		fieldPath := parentPath
		if !isEmbedded {
			fieldPath = buildFieldPath(parentPath, field.Name)
		}

		// Check if the field is required and zero
		if isRequiredAndZero(field, fieldValue) {
			if fieldPath != "" {
				missingFields = append(missingFields, fieldPath)
			} else {
				// Handle cases where the field might be a top-level embedded struct
				missingFields = append(missingFields, field.Name)
			}
			continue

		}

		// Recurse into nested structs, passing true only if the current field is embedded
		if checkNested(fieldValue) {
			nestedMissing := findMissing(fieldValue, fieldPath, field.Anonymous)
			missingFields = append(missingFields, nestedMissing...)
		}
	}

	return missingFields
}
