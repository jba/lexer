// Copyright 2021 Jonathan Amsterdam.

package lexer

import (
	"fmt"
	"io"
)

type (
	LexFunc  func(io.RuneScanner) (string, error)
	RunePred func(rune) bool
)

type Config struct {
	predFuncs []predFunc
}

func (c *Config) Install(p RunePred, f LexFunc) {
	c.predFuncs = append(c.predFuncs, predFunc{p, f})
}

type predFunc struct {
	p RunePred
	f LexFunc
}

type Lexer struct {
	rs io.RuneScanner
	c  *Config
}

func New(rs io.RuneScanner, c *Config) *Lexer {
	return &Lexer{
		rs: rs,
		c:  c,
	}
}

func (l *Lexer) Next() (string, error) {
loop:
	for {
		r, _, err := l.rs.ReadRune()
		if err != nil {
			return "", err
		}
		l.rs.UnreadRune()
		for _, pf := range l.c.predFuncs {
			if pf.p(r) {
				s, err := pf.f(l.rs)
				if err != nil {
					return "", err
				}
				if s != "" {
					return s, nil
				}
				continue loop
			}
		}
		return "", fmt.Errorf("illegal rune: %c (%[1]d)", r)
	}
}

func IsRune(r rune) RunePred {
	return func(r2 rune) bool { return r == r2 }
}

func ReadRune(r rune) LexFunc {
	return func(rs io.RuneScanner) (string, error) {
		r, _, err := rs.ReadRune()
		if err != nil {
			return "", err
		}
		return string(r), nil
	}
}

func SkipWhile(p RunePred) LexFunc {
	return func(rs io.RuneScanner) (string, error) {
		if err := skipWhile(rs, p); err != nil {
			return "", err
		}
		return "", nil
	}
}

func skipWhile(rs io.RuneScanner, p RunePred) error {
	for {
		r, _, err := rs.ReadRune()
		if err != nil {
			return err
		}
		if !p(r) {
			rs.UnreadRune()
			return nil
		}
	}
}

func ReadWhile(p RunePred) LexFunc {
	return func(rs io.RuneScanner) (string, error) {
		return readWhile(rs, p)
	}
}

func readWhile(rs io.RuneScanner, p RunePred) (string, error) {
	var s []rune

	for {
		r, _, err := rs.ReadRune()
		if err != nil && err != io.EOF {
			return "", err
		}
		if err == io.EOF || !p(r) {
			rs.UnreadRune()
			return string(s), nil
		}
		s = append(s, r)
	}
}

// func Not(p RunePred) RunePred {
// 	return func(r rune) bool { return !p(r) }
// }

// func In(s string) func(rune) bool {
// 	m := map[rune]bool{}
// 	for _, r := range s {
// 		m[r] = true
// 	}
// 	return func(r rune) bool {
// 		return m[r]
// 	}
// }

// Read the next rune. If the one after that is r2, include that in the token too.
func ReadOneOrTwo(r rune) LexFunc {
	return func(rs io.RuneScanner) (string, error) {
		r1, _, err := rs.ReadRune()
		if err != nil {
			return "", err
		}
		r2, _, err := rs.ReadRune()
		if r != r2 {
			rs.UnreadRune()
			return string(r), nil
		}
		if err != nil {
			return "", err
		}
		return string([]rune{r1, r2}), nil
	}
}
