[![Go Reference](https://pkg.go.dev/badge/github.com/gontainer/grouperror.svg)](https://pkg.go.dev/github.com/gontainer/grouperror)
[![Tests](https://github.com/gontainer/grouperror/actions/workflows/tests.yml/badge.svg)](https://github.com/gontainer/grouperror/actions/workflows/tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/gontainer/grouperror/badge.svg?branch=main)](https://coveralls.io/github/gontainer/grouperror?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/gontainer/grouperror)](https://goreportcard.com/report/github.com/gontainer/grouperror)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gontainer_grouperror&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gontainer_grouperror)

# Grouperror

This package provides a toolset to join and split errors.

```go
err := grouperror.Prefix("my group: ", errors.New("error1"), nil, errors.New("error2"))
errs := grouperror.Collection(err) // []error{error("my group: error1"), error("my group: error2")}
```

See [examples](examples_test.go).
