package lexer

import (
	"../token"
)

type Lexer struct {
	input        string // 入力の文字列
	position     int    // 読み終わったindex
	readPosition int    // 次に読むindex
	ch           byte   // 対象としている実際の文字
}

// lexerを初期化して生成
func New(input string) *Lexer {
	l := &Lexer{input: input}

	// readCharを1回だけ呼び出すと、先頭文字がsetされたり,
	// 各indexに値がsetされる。
	l.readChar()
	return l
}

// lexerのchに対して、tokenを返す。
// tokenを返した後はlexerの更新も行う。
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		tok = newToken(token.BANG, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// l.ch解析すべき文字をセットする。
// l.inputを超えるまで基本的にinputを読み続け、indexを保持し続ける
func (l *Lexer) readChar() {
	// 最後の文字に達した、もしくはoverしている時はchを0にセットする。
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		// まだ読んでいないreadPosition番号の文字を読む部分。
		l.ch = l.input[l.readPosition]
	}
	// 読みが完了した部分を書き出す
	l.position = l.readPosition
	// 次に読み出すindexの更新
	l.readPosition += 1
}
