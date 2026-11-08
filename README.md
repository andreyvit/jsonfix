Allow trailing commas, comments and bare keys in JSON in Go
===========================================================

[![Go reference](https://pkg.go.dev/badge/github.com/andreyvit/jsonfix.svg)](https://pkg.go.dev/github.com/andreyvit/jsonfix) ![Zero dependencies](https://img.shields.io/badge/deps-zero-brightgreen) [![Go Report Card](https://goreportcard.com/badge/github.com/andreyvit/jsonfix)](https://goreportcard.com/report/github.com/andreyvit/jsonfix)

```json5
// Comments are allowed
{
    "FOO": [1, 2, 3,], // trailing comma in an array
    BAR: 42,           // unquoted property name
    "BOZ": 24,         // trailing comma in an object
}
/* block comments
   are supported as well */
```


Why?
----

Fixes annoying problems with manually-written JSON in Go:

* lack of trailing commas,
* lack of comments.

We also support bare property names. You can use those for convenience, but the primary motivation is that if you feed JSON through one of those newfangled modern JS formatters, they will often unquote the keys, which is very annoying in some cases.

Why this library: it's a tiny JSON preprocessor, allowing to use standard `encoding/json`. No code duplication, no missing features (e.g. `DisallowUnknownFields`), no surprises, smooth upgrade path to JSON v2 (won't even require any changes here).


Usage
-----

Install:

    go get github.com/andreyvit/jsonfix

This preprocesses comments, trailing commas and bare keys; it is not a full JSON5 parser. Single-quoted strings and unquoted string values are unsupported.

Use:

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/andreyvit/jsonfix"
)

func main() {
	source := []byte(`{answer: 42, // a comment
}`)
	var result struct {
		Answer int `json:"answer"`
	}
	if err := json.Unmarshal(jsonfix.Bytes(source), &result); err != nil {
		panic(err)
	}
	fmt.Println(result.Answer) // 42
}
```

We preserve line numbers when preprocessing.


Contributing
------------

We accept contributions that:

* add better documentation and examples;
* add more tests;
* improve handling of invalid JSON (we don't know any problems so far, but it hasn't been a huge focus so far);
* implement single-quoted strings (one could make a case similar to bare keys here, if widespread formatters start preferring single quotes);
* fix bugs.

Out of scope (unless you convince us otherwise):

* multiline strings
* unquoted string values
* the rest of JSON5

We recommend [modd](https://github.com/cortesi/modd) (`go install github.com/cortesi/modd/cmd/modd@latest`) for continuous testing during development.

Maintain 100% coverage. It's not often the right choice, but it is for this library.


Changelog
---------

* 1.1.0 (2025-01-27): Added bare keys and block comments.

* 1.0.0 (2023-04-04): Initial release supporting trailing commas and line comments.


MIT license
-----------

Copyright (c) 2023–2025 Andrey Tarantsov. Published under the terms of the [MIT license](LICENSE).
