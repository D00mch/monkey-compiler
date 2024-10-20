package vm

import (
	"dumch/monkey/ast"
	"dumch/monkey/compiler"
	"dumch/monkey/lexer"
	"dumch/monkey/object"
	"dumch/monkey/parser"
	"fmt"
	"testing"
)

type vmTestCase struct {
	input    string
	expected any
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"(50 / 2) * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-50 + 100 + -50", 0},
		{"(5+10*2 +15/3)*2 + -10", 50},
	}

	runVmTests(t, tests)
}

func TestBooleanExpresisons(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5; })", true},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}
	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 } else { 20 }", 10},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
	}
	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{input: "let one=1; one", expected: 1},
		{input: "let one=1; let two=2; two+one", expected: 3},
		{input: "let one=1; let two=one+one; two+one", expected: 3},
	}
	runVmTests(t, tests)
}

func TestStringExpression(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}
	runVmTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1,2,3]", []int{1, 2, 3}},
		{"[1+2,3*4,5+6]", []int{3, 12, 11}},
	}
	runVmTests(t, tests)
}

func TestHashListernals(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    "{}",
			expected: map[object.HashKey]int64{},
		},
		{
			input: "{1+1: 2*2, 3+3: 4*4}",
			expected: map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}
	runVmTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", Null},
		{"[1, 2, 3][99]", Null},
		{"[1][-1]", Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", Null},
		{"{}[0]", Null},
	}
	runVmTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
		let fivePlusTen = fn() { 5 + 10; };
		fivePlusTen();`,
			expected: 15,
		},
		{
			input: `
		let one = fn() { 1; };
		let two = fn() { 2; };
		one() + two()`,
			expected: 3,
		},
		{
			input: `
		let a = fn() { 1 };
		let b = fn() { a() + 1 };
		let c = fn() { b() + 1 };
		c();`,
			expected: 3,
		},
	}
	runVmTests(t, tests)
}

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
		let earlyExit = fn() { return 99; 100; };
		earlyExit();`,
			expected: 99,
		},
		{
			input: `
		let earlyExit = fn() { return 99; return 100; };
		earlyExit(); `,
			expected: 99},
	}
	runVmTests(t, tests)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `let noReturn = fn() { };
				noReturn();`,
			expected: Null,
		},
		{
			input: `let noReturn = fn() {};
			        let noReturnTwo = fn() { noReturn(); };
				noReturn();
				noReturnTwo();`,
			expected: Null,
		},
	}
	runVmTests(t, tests)
}

func TestFirstClassFunction(t *testing.T) {
	test := []vmTestCase{
		{
			input: `let returnsOne = fn() { 1; }
				let returnsOneReturner = fn() { returnsOne; }
				returnsOneReturner()();`,
			expected: 1,
		},
	}
	runVmTests(t, test)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()
		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(
	t *testing.T,
	expected any,
	actual object.Object,
) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []int:
		err := testArrayObject(expected, actual)
		if err != nil {
			t.Errorf("testArrayObject failed: %s", err)
		}
	case map[object.HashKey]int64:
		err := testHashObject(expected, actual)
		if err != nil {
			t.Errorf("testHashObject failed: %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
		}
	}
}

func testArrayObject(expected []int, actual object.Object) error {
	array, ok := actual.(*object.Array)
	if !ok {
		return fmt.Errorf("object not Array: %T (%+v)", actual, actual)
	}

	if len(array.Elements) != len(expected) {
		return fmt.Errorf("wrong num of elements. want=%d, got=%d",
			len(expected), len(array.Elements))
	}

	for i, expectedElem := range expected {
		err := testIntegerObject(int64(expectedElem), array.Elements[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func testHashObject(expected map[object.HashKey]int64, actual object.Object) error {
	hashMap, ok := actual.(*object.Hash)
	if !ok {
		return fmt.Errorf("object is not HashMap. got=%T (%+v)", actual, actual)
	}
	if len(hashMap.Pairs) != len(expected) {
		return fmt.Errorf("hash has wrong number of Pairs. want=%d, got=%d",
			len(expected), len(hashMap.Pairs))
	}
	for expKey, expVal := range expected {
		pair, ok := hashMap.Pairs[expKey]
		if !ok {
			return fmt.Errorf("no pair for given key in Pairs")
		}
		err := testIntegerObject(expVal, pair.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
	}

	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)",
			actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q",
			actual, expected)
	}
	return nil
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
