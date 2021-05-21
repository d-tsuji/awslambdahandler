# awslambdahandler

[![Test Status](https://github.com/d-tsuji/awslambdahandler/workflows/test/badge.svg?branch=master)][actions]
[![Go Report Card](https://goreportcard.com/badge/github.com/d-tsuji/awslambdahandler)][go report card]
[![Apache-2.0 license](https://img.shields.io/badge/license-Apache2.0-blue.svg)][license]

[actions]: https://github.com/d-tsuji/awslambdahandler/actions?workflow=test
[go report card]: https://goreportcard.com/report/github.com/d-tsuji/awslambdahandler
[license]: https://github.com/d-tsuji/awslambdahandler/blob/master/LICENSE

`awslambdahandler` checks whether aws lambda handler signature is valid.

Valid AWS Lambda signatures are as follows

```
* func ()
* func () error
* func (TIn) error
* func () (TOut, error)
* func (TIn) (TOut, error)
* func (context.Context) error
* func (context.Context, TIn) error
* func (context.Context) (TOut, error)
* func (context.Context, TIn) (TOut, error)
```

## Detection sample

```go
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"What is your name?"`
	Age  int    `json:"How old are you?"`
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

func main() {
	lambda.Start(MyHandle)
}

func MyHandle(context.Context, MyEvent) MyResponse {
	return MyResponse{}
}
```

- Output

```
$ go vet -vettool=`which awslambdahandler` main.go
./main.go:19:14: lambda handler of "MyHandle" is invalid lambda signature, see https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start
```

For more detection samples, please see [examples](https://github.com/d-tsuji/awslambdahandler/blob/master/testdata/src/a/a.go).

## Usage

### `awslambdahandler` with go vet

`go vet` is a Go standard tool for analyzing source code.

1. Install `awslambdahandler`.
```sh
$ go install github.com/d-tsuji/awslambdahandler/cmd/awslambdahandler@latest
```

2. `awslambdahandler` execute
```sh
$ go vet -vettool=`which awslambdahandler` main.go
./main.go:33:14: lambda handler of "MyHandle" is invalid lambda signature, see https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start
```
