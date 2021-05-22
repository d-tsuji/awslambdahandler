package b

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func(customContext, interface{}) error { return nil })                             // OK
	lambda.StartWithContext(customContext{}, func(customContext, interface{}) error { return nil }) // OK
}

type customContext struct {
	context.Context
}
