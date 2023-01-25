// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/khaf/ojg/asm"
	"github.com/khaf/ojg/sen"
	"github.com/khaf/ojg/tt"
)

func TestReverse(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [reverse [a b c]]]
           [set $.asm.b [reverse [1 b 3]]]
         ]`,
		"{src: []}",
	)
	tt.Equal(t, `{a:[c b a] b:[3 b 1]}`, sen.String(root["asm"], &sopt))
}

func TestReverseArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"reverse", []any{}, 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestReverseArgType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"reverse", 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
