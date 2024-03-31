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
)

// ErrNotStructOrPtr is used if input is not a struct or a pointer to a struct.
var ErrNotStructOrPtr = errors.New("not a struct or pointer to a struct")
