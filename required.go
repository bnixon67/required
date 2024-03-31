// Copyright 2024 Bill Nixon. All rights reserved. Use of this source code
// is governed by the license found in the LICENSE file.

// Package required provides functions to check for required struct fields.
package required

import (
	"errors"
)

// ErrNotStructOrPtr is used if input is not a struct or a pointer to a struct.
var ErrNotStructOrPtr = errors.New("not a struct or pointer to a struct")
