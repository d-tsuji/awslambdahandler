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

const (
	doc = `awslambdahandler checks whether aws lambda handler signature is valid

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

	awsLambdaGoPath = "github.com/aws/aws-lambda-go/lambda"
)

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
		if identical(fn, awsLambdaGoPath, "Start") {
			argIdx = 0
		} else if identical(fn, awsLambdaGoPath, "StartWithContext") {
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
				if !isImplementContext(pass, vars.At(0).Type()) {
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
				if !isErrorType(vars.At(0).Type()) {
					valid = false
				}
			case 2:
				// if there are two return values, the second argument must be an error
				if !isErrorType(vars.At(1).Type()) {
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

func identical(obj types.Object, path, name string) bool {
	return analysisutil.RemoveVendor(obj.Pkg().Path()) == path && obj.Name() == name
}

func isErrorType(t types.Type) bool {
	return types.Identical(t, types.Universe.Lookup("error").Type())
}

func isImplementContext(pass *analysis.Pass, t types.Type) bool {
	// TODO(d-tsuji): If the code does not import the context package, it will be nil.
	// 	Therefore, it is possible that context.Context is implemented,
	// 	but a false positive will result in an error.
	ctxType := analysisutil.TypeOf(pass, "context", "Context")
	ctxIntf, ok := ctxType.Underlying().(*types.Interface)
	if !ok {
		return false
	}
	return types.Implements(t, ctxIntf)
}
