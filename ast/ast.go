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

// letやstatement以外の式を格納するノード
type Expression interface {
	Node
	expressionnode()
}

// プログラム全体を格納しているノードを表す。
// 「式」もstatementexpression構造体を定義したことによって
// このprogram構造体に格納する事が可能。
type Program struct {
	Statements []Statement
}

// そのprogram自身が保有している自身のstatementのリテラルを返す。
func (p *Program) String() string {
	var buf bytes.Buffer

	for _, s := range p.Statements {
		buf.WriteString(s.String())
	}
	return buf.String()
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

// ExpressionStatement は式文。式を表す部分全体をカバーするような構造体。
// program > expressionstatement > identifier ...
// みたいなイメージ。expressoionのラッパー??
// ExpressionStatement自体も、内部のASTをいくつか含む
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) statementnode()       {}
func (es *ExpressionStatement) expressionnode()      {}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// let文のast。
// statement interfaceの型を満たす。
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

// let文のast。
// statement interfaceの型を満たす。
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

// 識別子を表すast。
// expression interfaceを満たす。
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionnode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionnode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// PrefixExpression は前置演算子を表す構造体です。
// Operatorに演算子そのものを、Rightにオペランドを格納します。
// RightがExpressionを持つことができるので、再帰的にASTを構築することができる。
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionnode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionnode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}
