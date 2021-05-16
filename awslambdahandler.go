package awslambdahandler

import (
	"go/ast"
	"go/types"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

const doc = `awslambdahandler checks whether aws lambda handler signature is valid

Valid AWS Lambda signatures are as follows, 

* func ()
* func () error
* func (TIn) error
* func () (TOut, error)
* func (TIn) (TOut, error)
* func (context.Context) error
* func (context.Context, TIn) error
* func (context.Context) (TOut, error)
* func (context.Context, TIn) (TOut, error)

Where "TIn" and "TOut" are types compatible with the "encoding/json" standard library.
`

var Analyzer = &analysis.Analyzer{
	Name: "awslambdahandler",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	start := analysisutil.TypeOf(pass, "github.com/aws/aws-lambda-go/lambda", "Start")
	startWithCtx := analysisutil.TypeOf(pass, "github.com/aws/aws-lambda-go/lambda", "StartWithContext")

	errType := types.Universe.Lookup("error").Type()
	ctxType := analysisutil.TypeOf(pass, "context", "Context")

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		fn := typeutil.StaticCallee(pass.TypesInfo, call)
		if fn == nil {
			return // not a static call
		}

		// Classify the callee.
		recv := fn.Type().(*types.Signature).Recv()
		if recv != nil {
			return // not a function we are interested in
		}

		argIdx := -1
		if types.Identical(fn.Type(), start) {
			argIdx = 0
		} else if types.Identical(fn.Type(), startWithCtx) {
			argIdx = 1
		}
		if argIdx < 0 {
			return // not a function we are interested in
		}

		// handler must be a function
		t := pass.TypesInfo.Types[call.Args[argIdx]].Type
		s, ok := t.Underlying().(*types.Signature)
		if ok {
			valid := true

			// validate handler params
			// handler may take between 0 and two arguments
			vars := s.Params()
			switch vars.Len() {
			case 0, 1:
			case 2:
				// if there are two arguments, the first argument must satisfy the "context.Context" interface
				if !types.Identical(vars.At(0).Type(), ctxType) {
					valid = false
				}
			default:
				valid = false
			}

			// validate handler returns
			// handler may return between 0 and two arguments
			vars = s.Results()
			switch vars.Len() {
			case 0:
			case 1:
				// if there is one return value it must be an error
				if !types.Identical(vars.At(0).Type(), errType) {
					valid = false
				}
			case 2:
				// if there are two return values, the second argument must be an error
				if !types.Identical(vars.At(1).Type(), errType) {
					valid = false
				}
			default:
				valid = false
			}

			if valid {
				return // valid signatures
			}
		}

		var name string
		if i, ok := call.Args[argIdx].(*ast.Ident); ok {
			name = i.Name
		}
		if name == "" {
			pass.Reportf(call.Lparen, "invalid lambda signature, see https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start")
		} else {
			pass.Reportf(call.Lparen, `lambda handler of "%s" is invalid lambda signature, see https://pkg.go.dev/github.com/aws/aws-lambda-go/lambda#Start`, name)
		}

	})

	return nil, nil
}
