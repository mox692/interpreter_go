package parser

import (
	"fmt"
	"interpt/ast"
	"interpt/lexer"
	"interpt/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

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

	// tokenと中置き演算子(前置演算子)を解析する関数とを結び付けているmap
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// perserを生成する関数。lexerを引数にとる。
func New(l *lexer.Lexer) *Perser {
	p := &Perser{
		l:      l,
		errors: []string{},
	}
	// 前置演算子であるtokenと解析関数を結び付けている
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.persePrefixExpression)
	p.registerPrefix(token.MINUS, p.persePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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

// ParseProgram はinputからAST(Program構造体)を生成します。
// ParseProgramの中では、curTokenを進める毎にparseStatementを呼んで処理を移譲します。
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

// parserのcurTokenに応じて、構文解析関数を呼び出す関数です。
// tokenを進める作業はそれぞれのparse関数に移譲します
// curtokenがletかreturn以外だった場合は、全て式文というASTだとみなされます。
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

// let文を引き受けてLetstatementの構造体を返す。
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
	// *******************Todo: returnのvalueをどっかでsetしたいね。
	return stmt
}

// return文とlet文以外のtokenは一度全てこの関数にかけられます。
// exporessionstatement構造体を返します。
func (p *Perser) perseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	// この段階での「式文」を下位のperseExpressionに渡して、astノードを受け取る
	stmt.Expression = p.perseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Perser) perseExpression(precedence int) ast.Expression {
	// 現在読んでいるtokenに紐づけられたprefix構文解析関数がもしあれば、prefixに代入(式のはじめは必ずprefixタイプのtokenがくる)
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseError(p.curToken.Type)
		return nil
	}
	// 構文解析関数を実行、返ってきたASTノードをleftExpに格納。
	// prefixが[-]とか[!]だった際は、裏でtokenも進められる。
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		// もし中置き関数が見つからない(つまり、もう次のtokenがない場合は)
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		// infix構文解析関数を呼び出し、その結果をleftExpに代入していく(leftExpを使いまわしてるのに注意)
		// 左側のASTをどんどん取り込んでいくイメージ
		// このinfix関数の中で、leftExpがどんどん肥大化している様子が見えるはず？？
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Perser) persePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.perseExpression(PREFIX)

	// RIghtが代入されたASTを返す
	return expression
}

// parseInfixExpression中置き演算子の構文解析関数です。
// 最終的にinfixのASTを返します。
// prefix関数と違って引数を持つことに注意。
// 内部でparseEexpressionを呼んでいて、再起になってい明日。
func (p *Perser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	// leftとrightを結合して、新たな1つのASTを生成している部分
	expression.Right = p.perseExpression(precedence)

	// LeftとRightがセットされた中置きASTが返る
	return expression
}

func (p *Perser) noPrefixParseError(t token.TokenType) {
	msg := fmt.Sprintf("no parse Fn is found for token `%s`\n", t)
	p.errors = append(p.errors, msg)
}

func (p *Perser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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

func (p *Perser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Perser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Perser) Errors() []string {
	return p.errors
}

func (p *Perser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected %s, but got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// 第１引数のtokenに第２引数の構文解析関数を紐付けている。
func (p *Perser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Perser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// curTokenからintValueを抜き出してastを返す関数
func (p *Perser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}
