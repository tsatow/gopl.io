package ex08

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)

type lexer struct {
	scan scanner.Scanner
	token rune
}

func (lex *lexer) next() {
	lex.token = lex.scan.Scan()
}

func (lex *lexer) text() string {
	return lex.scan.TokenText()
}

func (lex *lexer) consume(want rune) {
	if lex.token != want {
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

type Token interface {}

type Symbol struct {
	Name string
}

type String struct {
	Value string
}

type Int struct {
	Value int
}

type StartList struct {}

type EndList struct {}

func read(lex *lexer) (Token, error) {
	switch lex.token {
	case scanner.Ident:
		// nilもSymbol
		symbol := Symbol{lex.text()}
		return symbol, nil
	case scanner.String:
		str := String{lex.text()}
		return str, nil
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text())
		return Int{i}, nil
	case '(':
		return StartList{}, nil
	case ')':
		return EndList{}, nil
	case scanner.EOF:
		return nil, io.EOF
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

type Decoder struct {
	lex *lexer
	err error
}

func NewDecoder(r io.Reader) *Decoder {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(r)
	return &Decoder{lex: lex}
}

func (decoder *Decoder) Token(out interface{}) (Token, error) {
	decoder.lex.next()
	return read(decoder.lex)
}