// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm_test

import (
	"testing"

	"github.com/khaf/ojg/asm"
	"github.com/khaf/ojg/sen"
	"github.com/khaf/ojg/tt"
)

func TestLteInt(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [lte 2 "$.src[2]" 3]]
           [set $.asm.b [lte 2 1]]
         ]`,
		"{src: [1 2 3]}",
	)
	tt.Equal(t, "{a:true b:false}", sen.String(root["asm"], &sopt))
}

func TestLteFloat(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a ["<=" 1.1 "$.src[1]" 2.2]]
           [set $.asm.b ["<=" 2.0 1.0]]
         ]`,
		"{src: [1.1 2.2]}",
	)
	tt.Equal(t, "{a:true b:false}", sen.String(root["asm"], &sopt))
}

func TestLteString(t *testing.T) {
	root := testPlan(t,
		`[
           [set $.asm.a [lte abc "$.src[1]" xyz]]
           [set $.asm.b [lte def abc]]
         ]`,
		"{src: [abc xyz]}",
	)
	tt.Equal(t, "{a:true b:false}", sen.String(root["asm"], &sopt))
}

func TestLteWrongType(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"set", "$.asm.i", []any{"lte", true, false}},
	})
	root := map[string]any{}
	err := p.Execute(root)
	tt.NotNil(t, err)
}

func TestLteWrongType2(t *testing.T) {
	p := asm.NewPlan([]any{
		[]any{"set", "$.asm.i", []any{"lte", 1, false}},
	})
	root := map[string]any{}
	err := p.Execute(root)
	tt.NotNil(t, err)
}
