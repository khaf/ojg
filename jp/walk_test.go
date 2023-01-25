// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp_test

import (
	"sort"
	"testing"

	"github.com/khaf/ojg/gen"
	"github.com/khaf/ojg/jp"
	"github.com/khaf/ojg/sen"
	"github.com/khaf/ojg/tt"
)

func TestWalk(t *testing.T) {
	data := map[string]any{"a": []any{1, 2, 3}, "b": nil}
	var paths []string
	jp.Walk(data, func(path jp.Expr, value any) { paths = append(paths, path.String()) })
	sort.Strings(paths)
	tt.Equal(t, `[$ $.a "$.a[0]" "$.a[1]" "$.a[2]" $.b]`, string(sen.Bytes(paths)))
}

func TestWalkNode(t *testing.T) {
	data := gen.Object{"a": gen.Array{gen.Int(1), gen.Int(2), gen.Int(3)}, "b": nil}
	var paths []string
	jp.Walk(data, func(path jp.Expr, value any) { paths = append(paths, path.String()) })
	sort.Strings(paths)
	tt.Equal(t, `[$ $.a "$.a[0]" "$.a[1]" "$.a[2]" $.b]`, string(sen.Bytes(paths)))
}
