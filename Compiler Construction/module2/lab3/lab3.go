package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokIdent TokenType = iota
	TokString
	TokAxiom
	TokIs
	TokOr
	TokEnd
	TokEpsilon
	TokEOF
)

func (t TokenType) String() string {
	switch t {
	case TokIdent:
		return "IDENT"
	case TokString:
		return "STRING"
	case TokAxiom:
		return "`axiom"
	case TokIs:
		return "`is"
	case TokOr:
		return "`or"
	case TokEnd:
		return "`end"
	case TokEpsilon:
		return "`epsilon"
	case TokEOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Col    int
}

type Lexer struct {
	input []rune
	pos   int
	line  int
	col   int
}

func NewLexer(text string) *Lexer {
	return &Lexer{
		input: []rune(text),
		pos:   0,
		line:  1,
		col:   1,
	}
}

func (l *Lexer) eof() bool {
	return l.pos >= len(l.input)
}

func (l *Lexer) peek() rune {
	if l.eof() {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) next() rune {
	if l.eof() {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) skipWhitespace() {
	for !l.eof() {
		ch := l.peek()

		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			l.next()
			continue
		}

		if ch == '$' {
			l.skipComment()
			continue
		}

		break
	}
}

func (l *Lexer) skipComment() {
	for !l.eof() {
		ch := l.next()
		if ch == '\n' {
			break
		}
	}
}

func (l *Lexer) readIdentifier() string {
	var b strings.Builder
	for !l.eof() {
		ch := l.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			b.WriteRune(l.next())
		} else {
			break
		}
	}
	return b.String()
}

func (l *Lexer) readString() (string, error) {
	startLine, startCol := l.line, l.col
	var b strings.Builder
	if l.next() != '"' {
		return "", fmt.Errorf("internal lexer error")
	}
	b.WriteRune('"')

	for !l.eof() {
		ch := l.next()
		b.WriteRune(ch)

		if ch == '\\' {
			if l.eof() {
				return "", fmt.Errorf("lexical error at %d:%d: unfinished escape sequence", startLine, startCol)
			}
			esc := l.next()
			b.WriteRune(esc)
			continue
		}

		if ch == '"' {
			return b.String(), nil
		}
	}

	return "", fmt.Errorf("lexical error at %d:%d: unterminated string", startLine, startCol)
}

func (l *Lexer) NextToken() (Token, error) {
	l.skipWhitespace()

	if l.eof() {
		return Token{Type: TokEOF, Lexeme: "EOF", Line: l.line, Col: l.col}, nil
	}

	startLine, startCol := l.line, l.col
	ch := l.peek()

	if ch == '`' {
		l.next()
		word := l.readIdentifier()
		lexeme := "`" + word
		switch lexeme {
		case "`axiom":
			return Token{Type: TokAxiom, Lexeme: lexeme, Line: startLine, Col: startCol}, nil
		case "`is":
			return Token{Type: TokIs, Lexeme: lexeme, Line: startLine, Col: startCol}, nil
		case "`or":
			return Token{Type: TokOr, Lexeme: lexeme, Line: startLine, Col: startCol}, nil
		case "`end":
			return Token{Type: TokEnd, Lexeme: lexeme, Line: startLine, Col: startCol}, nil
		case "`epsilon":
			return Token{Type: TokEpsilon, Lexeme: lexeme, Line: startLine, Col: startCol}, nil
		default:
			return Token{}, fmt.Errorf("lexical error at %d:%d: unknown keyword %q", startLine, startCol, lexeme)
		}
	}

	if ch == '"' {
		s, err := l.readString()
		if err != nil {
			return Token{}, err
		}
		return Token{Type: TokString, Lexeme: s, Line: startLine, Col: startCol}, nil
	}

	if unicode.IsLetter(ch) || ch == '_' {
		id := l.readIdentifier()
		return Token{Type: TokIdent, Lexeme: id, Line: startLine, Col: startCol}, nil
	}

	return Token{}, fmt.Errorf("lexical error at %d:%d: unexpected character %q", startLine, startCol, ch)
}

func LexAll(text string) ([]Token, error) {
	lexer := NewLexer(text)
	var tokens []Token
	for {
		tok, err := lexer.NextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Type == TokEOF {
			break
		}
	}
	return tokens, nil
}

type NonTerminal int

const (
	NGrammar NonTerminal = iota
	NRules
	NRule
	NAxiomOpt
	NAlts
	NAltsTail
	NAlt
	NSymbols
	NSymbolsTail
	NSymbol
)

func (n NonTerminal) String() string {
	switch n {
	case NGrammar:
		return "Grammar"
	case NRules:
		return "Rules"
	case NRule:
		return "Rule"
	case NAxiomOpt:
		return "AxiomOpt"
	case NAlts:
		return "Alts"
	case NAltsTail:
		return "AltsTail"
	case NAlt:
		return "Alt"
	case NSymbols:
		return "Symbols"
	case NSymbolsTail:
		return "SymbolsTail"
	case NSymbol:
		return "Symbol"
	default:
		return "<?>"
	}
}

type SymKind int

const (
	SymTerminal SymKind = iota
	SymNonTerminal
)

type SymbolRef struct {
	Kind SymKind
	Tok  TokenType
	NT   NonTerminal
}

func T(t TokenType) SymbolRef {
	return SymbolRef{Kind: SymTerminal, Tok: t}
}

func N(nt NonTerminal) SymbolRef {
	return SymbolRef{Kind: SymNonTerminal, NT: nt}
}

type ParseNode struct {
	Label    string
	Token    *Token
	Children []*ParseNode
}

func NewNode(label string) *ParseNode {
	return &ParseNode{Label: label}
}

type ParseTable map[NonTerminal]map[TokenType][]SymbolRef

func buildParseTable() ParseTable {
	return ParseTable{
		NGrammar: {
			TokAxiom: {N(NRules)},
			TokIdent: {N(NRules)},
			TokEOF:   {N(NRules)},
		},
		NRules: {
			TokAxiom: {N(NRule), N(NRules)},
			TokIdent: {N(NRule), N(NRules)},
			TokEOF:   {},
		},
		NRule: {
			TokAxiom: {N(NAxiomOpt), T(TokIdent), T(TokIs), N(NAlts), T(TokEnd)},
			TokIdent: {N(NAxiomOpt), T(TokIdent), T(TokIs), N(NAlts), T(TokEnd)},
		},
		NAxiomOpt: {
			TokAxiom: {T(TokAxiom)},
			TokIdent: {},
		},
		NAlts: {
			TokEpsilon: {N(NAlt), N(NAltsTail)},
			TokIdent:   {N(NAlt), N(NAltsTail)},
			TokString:  {N(NAlt), N(NAltsTail)},
		},
		NAltsTail: {
			TokOr:  {T(TokOr), N(NAlt), N(NAltsTail)},
			TokEnd: {},
		},
		NAlt: {
			TokEpsilon: {T(TokEpsilon)},
			TokIdent:   {N(NSymbols)},
			TokString:  {N(NSymbols)},
		},
		NSymbols: {
			TokIdent:  {N(NSymbol), N(NSymbolsTail)},
			TokString: {N(NSymbol), N(NSymbolsTail)},
		},
		NSymbolsTail: {
			TokIdent:  {N(NSymbol), N(NSymbolsTail)},
			TokString: {N(NSymbol), N(NSymbolsTail)},
			TokOr:     {},
			TokEnd:    {},
		},
		NSymbol: {
			TokIdent:  {T(TokIdent)},
			TokString: {T(TokString)},
		},
	}
}

func tokenDisplay(tok Token) string {
	switch tok.Type {
	case TokIdent, TokString:
		return fmt.Sprintf("%s(%s)\\n(%d:%d)", tok.Type.String(), tok.Lexeme, tok.Line, tok.Col)
	default:
		return fmt.Sprintf("%s\\n(%d:%d)", tok.Type.String(), tok.Line, tok.Col)
	}
}

func expectedTokens(table ParseTable, nt NonTerminal) string {
	row := table[nt]
	if row == nil {
		return ""
	}
	var parts []string
	for tok := range row {
		parts = append(parts, tok.String())
	}
	return strings.Join(parts, ", ")
}

type stackItem struct {
	symbol SymbolRef
	node   *ParseNode
}

func Parse(tokens []Token) (*ParseNode, error) {
	table := buildParseTable()
	root := NewNode(NGrammar.String())
	stack := []stackItem{
		{symbol: N(NGrammar), node: root},
	}

	pos := 0
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		lookahead := tokens[pos]

		if top.symbol.Kind == SymTerminal {
			if top.symbol.Tok != lookahead.Type {
				return nil, fmt.Errorf(
					"syntax error at %d:%d: expected %s, got %s (%q)",
					lookahead.Line, lookahead.Col,
					top.symbol.Tok.String(), lookahead.Type.String(), lookahead.Lexeme,
				)
			}

			tok := lookahead
			top.node.Token = &tok
			top.node.Label = tokenDisplay(tok)

			pos++
			continue
		}

		row := table[top.symbol.NT]
		prod, ok := row[lookahead.Type]
		if !ok {
			return nil, fmt.Errorf(
				"syntax error at %d:%d: unexpected %s (%q) while parsing %s; expected one of: %s",
				lookahead.Line, lookahead.Col,
				lookahead.Type.String(), lookahead.Lexeme,
				top.symbol.NT.String(), expectedTokens(table, top.symbol.NT),
			)
		}

		if len(prod) == 0 {
			top.node.Children = append(top.node.Children, NewNode("ε"))
			continue
		}

		children := make([]*ParseNode, 0, len(prod))
		for _, sym := range prod {
			var child *ParseNode
			if sym.Kind == SymTerminal {
				child = NewNode("")
			} else {
				child = NewNode(sym.NT.String())
			}
			children = append(children, child)
		}

		top.node.Children = append(top.node.Children, children...)

		for i := len(prod) - 1; i >= 0; i-- {
			stack = append(stack, stackItem{
				symbol: prod[i],
				node:   children[i],
			})
		}
	}

	if pos >= len(tokens) {
		return nil, fmt.Errorf("internal parser error")
	}
	if tokens[pos].Type != TokEOF {
		tok := tokens[pos]
		return nil, fmt.Errorf("syntax error at %d:%d: extra input starting from %q", tok.Line, tok.Col, tok.Lexeme)
	}

	return root, nil
}

func escapeLabel(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

func nodeLabel(n *ParseNode) string {
	if n.Token != nil {
		return tokenDisplay(*n.Token)
	}
	return n.Label
}

func writeDOT(w io.Writer, root *ParseNode) error {
	fmt.Fprintln(w, "digraph ParseTree {")
	fmt.Fprintln(w, "  rankdir=TB;")
	fmt.Fprintln(w, "  node [shape=box];")

	idCounter := 0
	ids := make(map[*ParseNode]string)

	var assign func(*ParseNode)
	assign = func(n *ParseNode) {
		id := fmt.Sprintf("n%d", idCounter)
		idCounter++
		ids[n] = id
		for _, c := range n.Children {
			assign(c)
		}
	}
	assign(root)

	var emit func(*ParseNode)
	emit = func(n *ParseNode) {
		fmt.Fprintf(w, "  %s [label=\"%s\"];\n", ids[n], escapeLabel(nodeLabel(n)))

		for _, c := range n.Children {
			fmt.Fprintf(w, "  %s -> %s;\n", ids[n], ids[c])
		}

		if len(n.Children) >= 2 {
			fmt.Fprint(w, "  { rank=same; ")
			for i, c := range n.Children {
				if i > 0 {
					fmt.Fprint(w, " -> ")
				}
				fmt.Fprint(w, ids[c])
			}
			fmt.Fprintln(w, " [style=invis] }")
		}

		for _, c := range n.Children {
			emit(c)
		}
	}
	emit(root)

	fmt.Fprintln(w, "}")
	return nil
}

func readInput() (string, error) {
	if len(os.Args) > 1 {
		data, err := os.ReadFile(os.Args[1])
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	var b strings.Builder
	reader := bufio.NewReader(os.Stdin)
	for {
		part, err := reader.ReadString('\n')
		b.WriteString(part)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func main() {
	input, err := readInput()
	if err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
		os.Exit(1)
	}

	tokens, err := LexAll(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tree, err := Parse(tokens)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := writeDOT(os.Stdout, tree); err != nil {
		fmt.Fprintln(os.Stderr, "write error:", err)
		os.Exit(1)
	}
}
