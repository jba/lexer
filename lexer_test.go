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
	in := "The 20 quick-foxes ran >=the (big dogs).<="
	want := []string{
		"The", "20", "quick", "-", "foxes", "ran", ">=", "the", "(", "big", "dogs", ")", ".", "<=",
	}

	l := lexer.New(strings.NewReader(in))
	l.Install(unicode.IsSpace, lexer.SkipWhile(unicode.IsSpace))
	l.Install(unicode.IsLetter, lexer.ReadWhile(isIdent))
	l.Install(unicode.IsDigit, lexer.ReadWhile(unicode.IsDigit))
	for _, r := range "+-().," {
		l.Install(lexer.IsRune(r), lexer.ReadRune(r))
	}
	l.Install(lexer.IsRune('>'), lexer.ReadOneOrTwo('='))
	l.Install(lexer.IsRune('<'), lexer.ReadOneOrTwo('='))

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
