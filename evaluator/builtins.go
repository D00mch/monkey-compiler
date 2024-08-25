package evaluator

import "dumch/monkey/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. "+
					"Got %d, want 1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},

	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. "+
					"Got %d, want 1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				if len(arg.Value) == 0 {
					return NULL
				}
				return &object.String{Value: arg.Value[0:1]}

			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[0]
			}

			return NULL
		},
	},

	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. "+
					"Got %d, want 1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				l := len(arg.Value)
				if l == 0 {
					return NULL
				}
				return &object.String{Value: arg.Value[l-1:]}

			case *object.Array:
				l := len(arg.Elements)
				if l == 0 {
					return NULL
				}
				return arg.Elements[l-1]
			}

			return NULL
		},
	},
}
