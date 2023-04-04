Allow trailing commas and comments in JSON in Go
================================================

[![Go reference](https://pkg.go.dev/badge/github.com/andreyvit/jsonfix.svg)](https://pkg.go.dev/github.com/andreyvit/jsonfix) ![Zero dependencies](https://img.shields.io/badge/deps-zero-brightgreen) ![under 100 LOC](https://img.shields.io/badge/size-%3C100%20LOC-green) ![100% coverage](https://img.shields.io/badge/coverage-100%25-green) [![Go Report Card](https://goreportcard.com/badge/github.com/andreyvit/jsonfix)](https://goreportcard.com/report/github.com/andreyvit/jsonfix)


Why?
----

Fixes annoying problems with manually-written JSON in Go:

* lack of trailing commas,
* lack of comments.

Why this library: it's a tiny JSON preprocessor, allowing to use standard `encoding/json`. No code duplication, no missing features (e.g. .DisallowUnknownFields()), no surprises.


Usage
-----

Install:

    go get github.com/andreyvit/jsonfix

Use:

```go
json.Unmarshal(jsonfix.Bytes(source), whatever)
```

We preserve line numbers when preprocessing.


TODO
----

- [ ] allow unquoted property names
- [ ] maybe allow block comments `/* ... */`


Contributing
------------

We accept contributions that:

* add better documentation and examples;
* add more tests;
* fix bugs;
* implement TODOs with full test coverage. 

Out of scope (unless you convince us otherwise):

* multiline strings
* single-quoted strings
* unquoted string values
* the rest of JSON5

We recommend [modd](https://github.com/cortesi/modd) (`go install github.com/cortesi/modd/cmd/modd@latest`) for continuous testing during development.

Maintain 100% coverage. It's not often the right choice, but it is for this library.


MIT license
-----------

Copyright (c) 2023 Andrey Tarantsov. Published under the terms of the [MIT license](LICENSE).
