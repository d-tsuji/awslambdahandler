package a

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
	// valid handler signature
	lambda.Start(func() {})                                                                       // OK
	lambda.Start(func() error { return nil })                                                     // OK
	lambda.Start(func(MyEvent) error { return nil })                                              // OK
	lambda.Start(func() (MyResponse, error) { return MyResponse{}, nil })                         // OK
	lambda.Start(func(context.Context) error { return nil })                                      // OK
	lambda.Start(func(context.Context, MyEvent) error { return nil })                             // OK
	lambda.Start(func(context.Context) (MyResponse, error) { return MyResponse{}, nil })          // OK
	lambda.Start(func(context.Context, MyEvent) (MyResponse, error) { return MyResponse{}, nil }) // OK

	// invalid handler signature
	lambda.Start("gopher")                                    // want `invalid lambda signature, see https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start`
	lambda.Start(func() string { return "gopher" })           // want `invalid lambda signature, see https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start`
	lambda.Start(func(MyEvent, MyEvent) error { return nil }) // want `invalid lambda signature, see https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start`
	lambda.Start(MyHandle)                                    // want `lambda handler of "MyHandle" is invalid lambda signature, see https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start`

	// valid handler signature
	ctx := context.Background()
	lambda.StartWithContext(ctx, func() {})                                                                       // OK
	lambda.StartWithContext(ctx, func() error { return nil })                                                     // OK
	lambda.StartWithContext(ctx, func(MyEvent) error { return nil })                                              // OK
	lambda.StartWithContext(ctx, func() (MyResponse, error) { return MyResponse{}, nil })                         // OK
	lambda.StartWithContext(ctx, func(context.Context) error { return nil })                                      // OK
	lambda.StartWithContext(ctx, func(context.Context, MyEvent) error { return nil })                             // OK
	lambda.StartWithContext(ctx, func(context.Context) (MyResponse, error) { return MyResponse{}, nil })          // OK
	lambda.StartWithContext(ctx, func(context.Context, MyEvent) (MyResponse, error) { return MyResponse{}, nil }) // OK

	// invalid handler signature
	lambda.StartWithContext(ctx, "gopher")                                    // want `invalid lambda signature, see https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start`
	lambda.StartWithContext(ctx, func() string { return "gopher" })           // want `invalid lambda signature, see https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start`
	lambda.StartWithContext(ctx, func(MyEvent, MyEvent) error { return nil }) // want `invalid lambda signature, see https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start`

	// not lambda func
	Start("gopher")                                                    // OK
	Start(func() string { return "gopher" })                           // OK
	Start(func(MyEvent, MyEvent) error { return nil })                 // OK
	Start(MyHandle)                                                    // OK
	StartWithContext(ctx, "gopher")                                    // OK
	StartWithContext(ctx, func() string { return "gopher" })           // OK
	StartWithContext(ctx, func(MyEvent, MyEvent) error { return nil }) // OK
}

func MyHandle() func(MyEvent, MyEvent) error        { return nil }
func Start(interface{})                             {}
func StartWithContext(context.Context, interface{}) {}
