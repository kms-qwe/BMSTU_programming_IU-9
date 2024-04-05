package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type Lexem struct {
	Tag   Tag
	Image string
}

type Tag int

const (
	ERROR Tag = 1 << iota
	NUMBER
	VAR
	PLUS
	MINUS
	MUL
	DIV
	LPAREN
	RPAREN
)

var UNDEFINEDLEXEM = Lexem{ERROR, ""}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: arith <expression>")
		return
	}

	expr := os.Args[1]
	lexems := make(chan Lexem)
	go lexer(expr, lexems)
	result, ok := analyzer(lexems)
	if !ok {
		fmt.Println("error")
	} else {
		fmt.Println(result)
	}
}

func lexer(expr string, lexems chan Lexem) {
	i := 0
	for i < len(expr) {
		switch expr[i] {
		case ' ':
			i++
		case '(':
			lexems <- Lexem{LPAREN, "("}
			i++
		case ')':
			lexems <- Lexem{RPAREN, ")"}
			i++
		case '+':
			lexems <- Lexem{PLUS, "+"}
			i++
		case '-':
			lexems <- Lexem{MINUS, "-"}
			i++
		case '*':
			lexems <- Lexem{MUL, "*"}
			i++
		case '/':
			lexems <- Lexem{DIV, "/"}
			i++
		default:
			if unicode.IsDigit(rune(expr[i])) {
				number := ""
				for i < len(expr) && unicode.IsDigit(rune(expr[i])) {
					number += string(expr[i])
					i++
				}
				lexems <- Lexem{NUMBER, number}
			} else if unicode.IsLetter(rune(expr[i])) {
				varname := ""
				for i < len(expr) && (unicode.IsLetter(rune(expr[i])) || unicode.IsDigit(rune(expr[i]))) {
					varname += string(expr[i])
					i++
				}
				lexems <- Lexem{VAR, varname}
			} else {
				lexems <- Lexem{ERROR, string(expr[i])}
				i++
			}
		}
	}
	close(lexems)
}

func analyzer(lexems chan Lexem) (result int, ok bool) {
	defer func() {
		if x := recover(); x != nil {
			result = 0
			ok = false
		}
	}()

	var parseExpression func() int
	var parseTerm func() int
	var parseFactor func() int

	currentLexem := UNDEFINEDLEXEM

	env := make(map[string]int)

	getVarValue := func(lexem Lexem) (value int) {
		if value, ok := env[lexem.Image]; ok {
			return value
		}
		_, err := fmt.Scanf("%d", &value)
		if err != nil {
			panic("Input error")
		}
		env[lexem.Image] = value
		return value
	}

	getNextLexem := func() Lexem {
		if currentLexem == UNDEFINEDLEXEM {
			lexem, ok := <-lexems
			if !ok {
				return UNDEFINEDLEXEM
			}
			if lexem.Tag&ERROR != 0 {
				panic("Error lexem")
			}
			return lexem
		}
		lexem := currentLexem
		currentLexem = UNDEFINEDLEXEM
		return lexem
	}

	expectLexem := func(expectedTag Tag) {
		lexem := getNextLexem()
		if lexem.Tag != expectedTag {
			panic("Expected " + lexem.Image + ", but got " + strconv.Itoa(int(lexem.Tag)))
		}
	}

	parseExpression = func() int {
		T := parseTerm()
		for {
			lexem := getNextLexem()
			switch lexem.Tag {
			case PLUS:
				T += parseTerm()
			case MINUS:
				T -= parseTerm()
			default:
				currentLexem = lexem
				return T
			}
		}
	}

	parseTerm = func() int {
		F := parseFactor()
		for {
			lexem := getNextLexem()
			switch lexem.Tag {
			case MUL:
				F *= parseFactor()
			case DIV:
				F /= parseFactor()
			default:
				currentLexem = lexem
				return F
			}
		}
	}

	parseFactor = func() int {
		lexem := getNextLexem()
		switch lexem.Tag {
		case NUMBER:
			value, _ := strconv.Atoi(lexem.Image)
			return value
		case VAR:
			val := getVarValue(lexem)
			return val
		case LPAREN:
			result := parseExpression()
			expectLexem(RPAREN)
			return result
		case MINUS:
			return -parseFactor()
		default:
			panic("Unexpected lexem")
		}
	}

	result = parseExpression()
	ok = true
	if _, ok := <-lexems; ok {
		panic("Channel is not closed")
	}
	return result, ok
}
