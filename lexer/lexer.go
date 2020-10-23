package lexer

import (
	"interpt/token"
)

// inputで入ってきたmokey言語は、全てこのlexerにかけられる。
// Lexerでは入力の文字列と、読み終えた文字列index,次に読み出す文字等を管理する。
type Lexer struct {
	input        string // 入力の文字列
	position     int    // 読み終わったindex
	readPosition int    // 次に読むindex
	ch           byte   // 対象としている実際の文字
}

// inputが初めてlexerにかけられた時に呼ばれる関数。
func New(input string) *Lexer {
	l := &Lexer{input: input}

	// readCharを1回だけ呼び出すと、先頭文字がsetされたり,
	// 各indexに値がsetされる。
	l.readChar()
	return l
}

// 字句を読みtokenを返すという,leserのメイン処理。
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
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
	// 1文字字句以外の処理
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// この後のreadCharを回避するため。
			return tok

		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			// *******************Todo: ここ何してる？
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

// nexttokenのヘルパー。
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// whitespaceを食い潰す関数
func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// nexttokenのヘルパー。
//
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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 単語の最初の1文字を読んでletterか数字かを判別
func isLetter(ch byte) bool {
	// good:byte列は文字列と比較できる
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// 識別しをひたすら読み進めて、その結果を返す関数。
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}
