package main

import (
	"fmt"
	"os"
	"regexp"
	"unicode"
	"unicode/utf8"
)

type Token struct {
	Tag   string
	Line  int
	Col   int
	Value string
}

type Domain struct {
	Tag      string
	Pattern  *regexp.Regexp
	Priority int
}

type Lexer struct {
	src  string
	pos  int
	line int
	col  int
}

func NewLexer(text string) *Lexer {
	return &Lexer{
		src:  text,
		pos:  0,
		line: 1,
		col:  1,
	}
}

func (lx *Lexer) eof() bool {
	return lx.pos >= len(lx.src)
}

func (lx *Lexer) peekSlice() string {
	if lx.eof() {
		return ""
	}
	return lx.src[lx.pos:]
}

func (lx *Lexer) advanceByString(s string) {
	i := 0
	for i < len(s) {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			lx.col++
			i++
			continue
		}

		if r == '\r' {
			if i+size < len(s) {
				r2, size2 := utf8.DecodeRuneInString(s[i+size:])
				if r2 == '\n' {
					i += size + size2
					lx.line++
					lx.col = 1
					continue
				}
			}
			i += size
			lx.line++
			lx.col = 1
			continue
		}
		if r == '\n' {
			i += size
			lx.line++
			lx.col = 1
			continue
		}

		i += size
		lx.col++
	}
	lx.pos += len(s)
}

func (lx *Lexer) skipWhitespace() {
	for !lx.eof() {
		r, size := utf8.DecodeRuneInString(lx.src[lx.pos:])
		if r == utf8.RuneError && size == 1 {
			return
		}
		if !unicode.IsSpace(r) {
			return
		}
		lx.advanceByString(lx.src[lx.pos : lx.pos+size])
	}
}

func sanitizeValue(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		switch r {
		case '\n':
			out = append(out, []rune(`\n`)...)
		case '\r':
			out = append(out, []rune(`\r`)...)
		case '\t':
			out = append(out, []rune(`\t`)...)
		default:
			out = append(out, r)
		}
	}
	return string(out)
}

func (lx *Lexer) NextToken(domains []Domain) (Token, bool) {
	for {
		lx.skipWhitespace()
		if lx.eof() {
			return Token{}, false
		}

		text := lx.peekSlice()

		bestIdx := -1
		bestLen := -1
		bestPri := 1 << 30

		for i := range domains {
			d := domains[i]
			loc := d.Pattern.FindStringIndex(text)
			if loc == nil {
				continue
			}
			mLen := loc[1]
			if mLen > bestLen || (mLen == bestLen && d.Priority < bestPri) {
				bestLen = mLen
				bestPri = d.Priority
				bestIdx = i
			}
		}

		if bestIdx >= 0 {
			startLine, startCol := lx.line, lx.col
			lex := text[:bestLen]
			lx.advanceByString(lex)
			return Token{
				Tag:   domains[bestIdx].Tag,
				Line:  startLine,
				Col:   startCol,
				Value: lex,
			}, true
		}

		fmt.Printf("syntax error (%d,%d)\n", lx.line, lx.col)

		for {
			if lx.eof() {
				return Token{}, false
			}

			_, size := utf8.DecodeRuneInString(lx.src[lx.pos:])
			if size <= 0 {
				size = 1
			}
			lx.advanceByString(lx.src[lx.pos : lx.pos+size])

			lx.skipWhitespace()
			if lx.eof() {
				return Token{}, false
			}

			text = lx.peekSlice()
			found := false
			for i := range domains {
				loc := domains[i].Pattern.FindStringIndex(text)
				if loc != nil {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <input_utf8_file>\n", os.Args[0])
		os.Exit(2)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
		os.Exit(1)
	}
	text := string(data)

	domains := []Domain{
		{Tag: "COMMENT", Pattern: regexp.MustCompile(`(?s)^\{\-.*?\-\}`), Priority: 0},
		{Tag: "COMMENT", Pattern: regexp.MustCompile(`^--[^\r\n]*`), Priority: 1},

		{Tag: "WHERE", Pattern: regexp.MustCompile(`^where\b`), Priority: 2},
		{Tag: "ARROW", Pattern: regexp.MustCompile(`^->`), Priority: 3},
		{Tag: "FATARROW", Pattern: regexp.MustCompile(`^=>`), Priority: 4},

		{Tag: "OP", Pattern: regexp.MustCompile(`^` + "`" + `[A-Za-z][A-Za-z0-9]*` + "`"), Priority: 5},

		{Tag: "IDENT", Pattern: regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*`), Priority: 6},
		{Tag: "INT", Pattern: regexp.MustCompile(`^[0-9]+`), Priority: 7},

		{Tag: "OP", Pattern: regexp.MustCompile(`^[!#$%&*+./<=>?@\\^|\-~]+`), Priority: 8},
		// Координата в полярной системе вида ddd@ddd
		{Tag: "POL", Pattern: regexp.MustCompile(`^[0-9]+@[0-9]+`), Priority: 9},
	}

	lx := NewLexer(text)
	for {
		tok, ok := lx.NextToken(domains)
		if !ok {
			break
		}
		fmt.Printf("%s (%d, %d): %s\n", tok.Tag, tok.Line, tok.Col, sanitizeValue(tok.Value))
	}
}
