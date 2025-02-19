# [![{}j](assets/ojg_comet.svg)](https://github.com/khaf/ojg)

[![Build Status](https://github.com/khaf/ojg/actions/workflows/CI.yml/badge.svg)](https://github.com/khaf/ojg/actions)
[![Coverage Status](https://coveralls.io/repos/github/ohler55/ojg/badge.svg?branch=master)](https://coveralls.io/github/ohler55/ojg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/khaf/ojg)](https://goreportcard.com/report/github.com/khaf/ojg)

Optimized JSON for Go is a high performance parser with a variety of
additional JSON tools. OjG is optimized to processing huge data sets
where data does not necessarily conform to a fixed structure.

## Features

 - Fast JSON parser. Check out the cmd/benchmarks app in this repo.
 - Full JSONPath implemenation that operates on simple types as well as structs.
 - Generic types. Not the proposed golang generics but type safe JSON elements.
 - Fast JSON validator (7 times faster with io.Reader).
 - Fast JSON writer with a sort option (4 times faster).
 - JSON builder from JSON sources using a simple assembly plan.
 - Simple data builders using a push and pop approach.
 - Object encoding and decoding using an approach similar to that used with Oj for Ruby.
 - [Simple Encoding Notation](sen.md), a lazy way to write JSON omitting commas and quotes.

## Using

A basic Parse:

```golang
    obj, err := oj.ParseString(`{
        "a":[
            {"x":1,"y":2,"z":3},
            {"x":2,"y":4,"z":6}
        ]
    }`)
```

Using JSONPath expressions:

```golang
    x, err := jp.ParseString("a[?(@.x > 1)].y")
    ys := x.Get(obj)
    // returns [4]
```

The **oj** command (cmd/oj) uses JSON path for filtering and
extracting JSON elements. It also includes sorting, reformatting, and
colorizing options.

```
$ oj -m "(@.name == 'Pete')" myfile.json

```

More complete examples are available in the go docs for most
functions. The example for [Unmarshalling
interfaces](oj/example_interface_test.go) demonstrates a feature that
allows interfaces to be marshalled and unmarshalled.

## Installation
```
go get github.com/khaf/ojg
go get github.com/khaf/ojg/cmd/oj

```

or just import in your `.go` files.

```
import (
    "github.com/khaf/ojg/alt"
    "github.com/khaf/ojg/asm"
    "github.com/khaf/ojg/gen"
    "github.com/khaf/ojg/jp"
    "github.com/khaf/ojg/oj"
    "github.com/khaf/ojg/sen"
)
```

To build and install the `oj` application:

```
go install ./...
```

## Benchmarks

Higher numbers (longer bars) are better.

```
Parse string/[]byte
       json.Unmarshal           55916 ns/op    17776 B/op    334 allocs/op
         oj.Parse               39570 ns/op    18488 B/op    429 allocs/op
   oj-reuse.Parse               17881 ns/op     5691 B/op    364 allocs/op

   oj-reuse.Parse        █████████████████████▉ 3.13
         oj.Parse        █████████▉ 1.41
       json.Unmarshal    ▓▓▓▓▓▓▓ 1.00

Parse io.Reader
       json.Decode              63029 ns/op    32449 B/op    344 allocs/op
         oj.ParseReader         34289 ns/op    22583 B/op    430 allocs/op
   oj-reuse.ParseReader         25094 ns/op     9788 B/op    365 allocs/op
         oj.TokenizeLoad        13610 ns/op     6072 B/op    157 allocs/op

         oj.TokenizeLoad ████████████████████████████████▍ 4.63
   oj-reuse.ParseReader  █████████████████▌ 2.51
         oj.ParseReader  ████████████▊ 1.84
       json.Decode       ▓▓▓▓▓▓▓ 1.00

to JSON with indentation
       json.Marshal             78762 ns/op    26978 B/op    352 allocs/op
         oj.JSON                 7662 ns/op        0 B/op      0 allocs/op
        sen.Bytes                9053 ns/op        0 B/op      0 allocs/op

         oj.JSON         ███████████████████████████████████████████████████████████████████████▉ 10.28
        sen.Bytes        ████████████████████████████████████████████████████████████▉ 8.70
       json.Marshal      ▓▓▓▓▓▓▓ 1.00
```

See [all benchmarks](benchmarks.md)

[Compare Go JSON parsers](https://github.com/ohler55/compare-go-json)

## Releases

See [CHANGELOG.md](CHANGELOG.md)

## Links

- *Documentation*: [https://pkg.go.dev/github.com/khaf/ojg](https://pkg.go.dev/github.com/khaf/ojg)

- *GitHub* *repo*: https://github.com/khaf/ojg

- *JSONPath* description: https://goessner.net/articles/JsonPath

- *JSONPath Comparisons*: https://cburgmer.github.io/json-path-comparison

- *Go Report Card*: https://goreportcard.com/report/github.com/khaf/ojg

#### Links of Interest

 - *Oj, a Ruby JSON parser*: http://www.ohler.com/oj/doc/index.html also at https://github.com/ohler55/oj

 - *OjC, a C JSON parser*: http://www.ohler.com/ojc/doc/index.html also at https://github.com/ohler55/ojc

 - *Fast XML parser and marshaller on GitHub*: https://github.com/ohler55/ox

 - *Agoo, a high performance Ruby web server supporting GraphQL on GitHub*: https://github.com/ohler55/agoo

 - *Agoo-C, a high performance C web server supporting GraphQL on GitHub*: https://github.com/ohler55/agoo-c

#### Contributing

+ Provide a Pull Request off the `develop` branch.
+ Report a bug
+ Suggest an idea
