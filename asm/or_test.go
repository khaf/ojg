// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/khaf/ojg/asm"
	"github.com/khaf/ojg/sen"
	"github.com/khaf/ojg/tt"
)

func TestOrTrue(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [or false "$.src[0]"]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "true", sen.String(root["asm"]))
}

func TestOrFalse(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [or false "$.src[1]" false]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "false", sen.String(root["asm"]))
}

func TestOrNull(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm [or "$.src[2]"]]
         ]`,
		"{src: [true false]}",
	)
	tt.Equal(t, "false", sen.String(root["asm"]))
}

func TestOrNotBool(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"or", 1, 2},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
