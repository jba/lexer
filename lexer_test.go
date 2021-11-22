// Copyright 2021 Jonathan Amsterdam.

package lexer_test

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/jba/lexer"
)

func Test(t *testing.T) {
	in := "The 20 quick-foxes ran >=the (big dogs).<= yes"
	want := []string{
		"The", "20", "quick", "-", "foxes", "ran", ">=", "the", "(", "big", "dogs", ")", ".", "<=", "yes",
	}

	var c lexer.Config
	c.Install(unicode.IsSpace, lexer.SkipWhile(unicode.IsSpace))
	c.Install(unicode.IsLetter, lexer.ReadWhile(isIdent))
	c.Install(unicode.IsDigit, lexer.ReadWhile(unicode.IsDigit))
	for _, r := range "+-().," {
		c.Install(lexer.IsRune(r), lexer.ReadRune(r))
	}
	c.Install(lexer.IsRune('>'), lexer.ReadOneOrTwo('='))
	c.Install(lexer.IsRune('<'), lexer.ReadOneOrTwo('='))

	l := lexer.New(strings.NewReader(in), &c)

	var got []string
	for {
		tok, err := l.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, tok)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got\n%s\nwant\n%s", got, want)
	}
}

func isIdent(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
