package ast

import (
	"bytes"
	"interpt/token"
)

// astのパッケージ。
// 文章を各コンポーネントに分けていく。

// Node :astにおけるルートノード。
// これ以下のノードinterfaceにはこのNodeを保有させることで依存関係を作っている。
type Node interface {
	TokenLiteral() string
	// String()は、各ノードがカバーしている部分のリテラルをstringで返します。
	String() string
}

// letやreturnのようなstatementを格納するstruct
type Statement interface {
	Node
	statementnode()
}

type Expression interface {
	Node
	expressionnode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var buf bytes.Buffer

	for _, s := range p.Statements {
		buf.WriteString(s.String())
	}
	return buf.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) statementnode()       {}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier // let文に続き識別子
	Value Expression  // 右辺の式
}

func (ls *LetStatement) statementnode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(ls.Token.Literal + " ")
	buf.WriteString(ls.Name.Value)
	buf.WriteString(" = ")

	if ls.Value != nil {
		// *******************Todo:
		buf.WriteString(ls.Value.String())
	}

	buf.WriteString(";")

	return buf.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementnode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(rs.Token.Literal + " ")

	if rs.ReturnValue != nil {
		buf.WriteString(rs.ReturnValue.String())
	}

	return buf.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) identifiermentnode()  {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	return i.Value
}

// 自身のtokenを返す
func (p *Program) TokenLiteral() string {
	// *******************Todo: Program構造体ってどの単位で作成するんだ？
	// 「letノード 」「returnノード」みたいに、node事にProgramインスタンスが生成されると予想
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
