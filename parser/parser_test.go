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
	let x = 5
	let y = 10
	let foo = 4583
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatal("parse err !!")
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

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	return false
}
