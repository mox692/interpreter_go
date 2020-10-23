package lexer

import (
	"testing"

	"../token"
)

func TestNextToken(t *testing.T) {
	input := `let age =  1;`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "age"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
	}

	lexer := New(input)

	for _, tt := range tests {

		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("got: %s, want: %s", tok.Type, tt.expectedType)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("got: %s, want: %s", tok.Literal, tt.expectedLiteral)
		}
	}
}
