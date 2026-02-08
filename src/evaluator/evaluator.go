package evaluator

import (
	"fmt"
	"rubo-lang/src/parser"
)

const HOT_THRESHOLD = 1000

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

func Eval(node parser.Node, env *Environment) Object {
	switch node := node.(type) {

	// Statements
	case *parser.Program:
		return evalProgram(node, env)

	case *parser.ExpressionStatement:
		return Eval(node.Expression, env)

	case *parser.BlockStatement:
		return evalBlockStatement(node, env)

	// Expressions
	case *parser.IntegerLiteral:
		return &Integer{Value: node.Value}

	case *parser.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(node.Operator, left, right)

	case *parser.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)

	case *parser.Identifier:
		return evalIdentifier(node, env)

	case *parser.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		fn := &Function{Name: node.Name.Value, Parameters: params, Env: env, Body: body}
		env.SetValue(node.Name.Value, fn)
		return fn

	case *parser.CallExpression:
		function := Eval(node.Function, env)
		args := evalExpressions(node.Arguments, env)
		return applyFunction(function, args)
	}

	return nil
}

func evalProgram(program *parser.Program, env *Environment) Object {
	var result Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalBlockStatement(block *parser.BlockStatement, env *Environment) Object {
	var result Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil && result.Type() == RETURN_VALUE_OBJ {
			return result
		}
	}

	return result
}

func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left != nil && right != nil && left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return NULL
	}
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right == nil || right.Type() != INTEGER_OBJ {
		return NULL
	}

	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		return &Integer{Value: leftVal / rightVal}
	default:
		return NULL
	}
}

func evalIdentifier(node *parser.Identifier, env *Environment) Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return NULL
}

func evalExpressions(exps []parser.Expression, env *Environment) []Object {
	var result []Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn Object, args []Object) Object {
	function, ok := fn.(*Function)
	if !ok {
		return NULL
	}

	if function.Name != "" {
		// Step A: Check if a "Native Version" exists
		if nativeFn, exists := function.Env.NativeFunctions[function.Name]; exists {
			// Convert args to int64 for the native function (assuming integers for now)
			var intArgs []int64
			for _, arg := range args {
				if integer, ok := arg.(*Integer); ok {
					intArgs = append(intArgs, integer.Value)
				}
			}
			// Step B: If it exists, call it.
			result := nativeFn(intArgs...)
			return &Integer{Value: result}
		}

		// Profiler logic
		function.Env.CallCounts[function.Name]++
		if function.Env.CallCounts[function.Name] == HOT_THRESHOLD {
			fmt.Printf("OPTIMIZATION TRIGGERED: Function %s is hot. Handing off to Rust...\n", function.Name)
			if function.Env.OnHotFunction != nil {
				function.Env.OnHotFunction(function)
			}
		}
	}

	// Step C: If it doesn't, use the standard Eval
	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.SetValue(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}
