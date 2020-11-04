package parser

import (
	"interpt/ast"
	"interpt/lexer"
	"testing"
)

// inputからテストを作り込んでいる。
// そのため、statementsの数が期待通りであるか、みたいなテストも入っている。
// 後半で、1つ1つのstatementに対するテストは別のTestStatementというテスト関数に投げている。
func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foo = 4583;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("parse err")
	}
	// *******************todo: let, 識別し, 式　で長さが3か、入力のlet文が3つで3か
	if len(program.Statements) != 3 {
		t.Fatalf("")
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

// TestLetStatementsで生成された、1つ1つのletstatementに対するテスト。
func testLetStatement(t *testing.T, s ast.Statement, ident string) bool {

	// 一番初めのtokenのリテラルがletじゃなかったらerr
	if s.TokenLiteral() != "let" {
		t.Errorf("tokenLiteral is not '%s', but  '%s'", "let", s.TokenLiteral())
	}

	// sがLetStatement型出なかったらerr
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("ast.Statement is not %sstatement, but %Tstatement", "let", s)
		return false
	}

	// inputから読み込んだ識別子と、テストテーブルとして用意した識別子が一致しなかったらerror
	if letStmt.Name.Value != ident {
		t.Errorf("ident want %s, but got %s", ident, letStmt.Name.Value)
		return false
	}

	// これをテストする意義があまりわからない。。
	if letStmt.Name.TokenLiteral() != ident {
		t.Errorf("")
		return false
	}
	return true
}

//  識別子を識別子であると判別するテスト
func TestIdentifierExpression(t *testing.T) {
	input := "foooobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enoutg statement,want 1, got %d", len(program.Statements))
	}

	// statement interfaceの中で、ExpressionStatement型であるか？(ひとまず式文として判定させる)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program has not enoutg statement,want 1, got %d", len(program.Statements))
	}

	// expressioninterfaceの中で、さらにidentifierか？
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("program has not enoutg statement,want 1, got %d", len(program.Statements))
	}

	if ident.Value != "foooobar" {
		t.Fatalf("program has not enoutg statement,want 1, got %d", len(program.Statements))
	}

}

func checkParserErrors(t *testing.T, p *Perser) {
	errors := p.errors
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors.\n", len(errors))
	for _, v := range errors {
		t.Errorf("parser error: %s\n", v)
	}

	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 9fdsa;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("parse err")
	}
	// *******************todo: let, 識別し, 式　で長さが3か、入力のlet文が3つで3か
	if len(program.Statements) != 3 {
		t.Fatalf("")
	}

	for _, stmt := range program.Statements {
		_, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("not return statement..")
		}
	}
}
