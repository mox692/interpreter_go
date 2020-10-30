package parser

import (
	"fmt"
	"interpt/ast"
	"interpt/lexer"
	"interpt/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// parserパッケージでは、lexerで生成されたトークン列から、具体的なASTを作成する。
// perserはparse(token=>AST)の手段を提供するParseProgramメソッドを実装している構造体である。
// lはlexerのポインタを保有している。
// curTokenは現在解析を行っているtokenを保有している。
// peekは次のtokenを保有している。
type Perser struct {
	l         *lexer.Lexer
	errors    []string
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// perserを生成する関数。lexerを引数にとる。
func New(l *lexer.Lexer) *Perser {
	p := &Perser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}

// parserの解析tokenを1つ進める。
// perserに紐づいているlexer内部のnextToken()メソッドを呼び出すことによに実現させている。
func (p *Perser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// inputからASTを生成。
func (p *Perser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.perseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Perser) perseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.perseLetStatement()
	case token.RETURN:
		return p.perseReturnStatement()
	default:
		return p.perseExpressionStatement()
	}
}

func (p *Perser) perseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// 次の文字が識別子出なかったらerr
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// todo.....
	// stmt.Value =
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Perser) perseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Perser) perseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.perseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Perser) peekTokenIs(t token.TokenType) bool {
	if p.peekToken.Type == t {
		return true
	} else {
		return false
	}
}

func (p *Perser) curTokenIs(t token.TokenType) bool {
	if p.curToken.Type == t {
		return true
	} else {
		return false
	}
}

func (p *Perser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Perser) Errors() []string {
	return p.errors
}

func (p *Perser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected %s, but got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Perser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Perser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
