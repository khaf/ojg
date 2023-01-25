// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import "github.com/khaf/ojg/oj"

type data struct {
	src string
	// Empty means no error expected while non empty should be compared
	// err.Error().
	expect  string
	value   any
	onlyOne bool
	options *oj.Options
	indent  int
}
