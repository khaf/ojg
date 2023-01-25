// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"
	"time"

	"github.com/khaf/ojg/asm"
	"github.com/khaf/ojg/sen"
	"github.com/khaf/ojg/tt"
)

var sopt = sen.Options{Sort: true, TimeFormat: time.RFC3339Nano}

func testPlan(t *testing.T, plan, root string) map[string]any {
	parser := sen.Parser{}
	val, err := parser.Parse([]byte(plan))
	tt.Nil(t, err)
	list, _ := val.([]any)
	p := asm.NewPlan(list)

	val, err = parser.Parse([]byte(root))
	tt.Nil(t, err)
	r, _ := val.(map[string]any)
	err = p.Execute(r)
	tt.Nil(t, err)

	return r
}
