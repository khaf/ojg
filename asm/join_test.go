// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/khaf/ojg/asm"
	"github.com/khaf/ojg/sen"
	"github.com/khaf/ojg/tt"
)

func TestJoin(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [join [a b c] '+']]
           [set $.asm.b [join [a b c]]]
         ]`,
		"{src: []}",
	)
	opt := sopt
	opt.Indent = 2
	tt.Equal(t,
		`{
  a: a+b+c
  b: abc
}`, sen.String(root["asm"], &opt))
}

func TestJoinArgCount(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"join", []any{}, "x", 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestJoinArgType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"join", 1, "x"},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestJoinArgType2(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"join", []any{}, 1},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}

func TestJoinArgType3(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"join", []any{"x", 3}},
	})
	err := p.Execute(map[string]any{})
	tt.NotNil(t, err)
}
