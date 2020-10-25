package ast

import "interpt/token"

// Node :astにおけるルートノード。
// これ以下のノードinterfaceにはこのNodeを保有させることで依存関係を作っている。
type Node interface {
	TokenLiteral() string
}

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

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementnode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) identifiermentnode()  {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

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
