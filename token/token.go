package token

import (
	"bytes"
	"fmt"
	"unicode"
)

type Type string

const (
	ILLEGAL Type = "ILLEGAL"
	EOF     Type = "EOF"

	LeftBrace    Type = "{"
	RightBrace   Type = "}"
	LeftBracket  Type = "["
	RightBracket Type = "]"
	Comma        Type = ","
	Colon        Type = ":"
	Quote        Type = "\""

	String  Type = "STRING"
	Number  Type = "NUMBER"
	Boolean Type = "BOOLEAN"
	Null    Type = "NULL"
)

func (t Type) String() string {
	return string(t)
}

type Token struct {
	Type Type
	Val  string
}

func Tokenizer(input []byte) []Token {
	curr := 0
	var tokens []Token
	stack := NewStack()
	var prevTokenType = ILLEGAL
	for curr < len(input) {
		char := input[curr]

		currTokenType := determineTokenType(char, input, curr)

		if unicode.IsSpace(rune(char)) {
			curr++
			continue
		}

		if !isValidSequences(prevTokenType, currTokenType) {
			errorToken := Token{
				Type: ILLEGAL,
				Val:  fmt.Sprintf("invalid token sequence"),
			}
			tokens = append(tokens, errorToken)
			return tokens
		}

		switch currTokenType {

		case LeftBrace:
			stack.Push(LeftBrace)
			tokens = append(tokens, Token{Type: LeftBrace,Val: string(char)})

		case RightBrace:
			if stack.Peek() == LeftBrace {
				stack.Pop()
			}
			tokens = append(tokens, Token{Type:RightBrace, Val: string(char)})

		case LeftBracket:
			stack.Push(LeftBracket)
			tokens = append(tokens, Token{Type: LeftBracket, Val: string(char)})

		case RightBracket:
			if stack.Peek() == LeftBracket {
				stack.Pop()
			}
			tokens = append(tokens, Token{Type: RightBracket, Val: string(char)})
		case Comma:
			tokens = append(tokens, Token{Type:Comma, Val: string(char)})
		case Colon:
			tokens = append(tokens, Token{Type:Colon, Val: string(char)})
		case Quote:
			curr++
			start := curr

			for curr < len(input) && input[curr] != '"' {
				curr++
			}

			if curr < len(input) {
				value := input[start:curr]
				tokens = append(tokens, Token{Type:String, Val: string(value)})
			} else {
				tokens = append(tokens, Token{Type:ILLEGAL, Val: "Unclosed string literal"})
				return tokens
			}
		default:
			if unicode.IsDigit(rune(char)) {
				start := curr

				for curr < len(input) && isDigit(input[curr]) {
					curr++
				}

				if curr != len(input) && !isTerminatingCharacter((input[curr])) {
					tokens = append(tokens, Token{Type:ILLEGAL, Val: "Invalid number format"})
					return tokens
				} else {
					prevTokenType = Number
					value := input[start:curr]
					tokens = append(tokens, Token{Type: Number, Val: string(value)})
				}
				continue
			} else if char == 't' || char == 'f' {
				if isBoolean(input, curr) {
					var length int 
					if bytes.Equal(input[curr:curr+4], []byte("true")) {
						length = 4
					} else {
						length = 5
					}

					prevTokenType = Boolean
					value := input[curr:curr+length]
					tokens = append(tokens, Token{Type: Boolean, Val: string(value)})

					curr += length
					continue
				} else {
					tokens = append(tokens, Token{Type: ILLEGAL, Val: "Invalid boolean literal"})
					return tokens
				}
			} else if char == 'n' {
				if isNull(input, curr) {
					if bytes.Equal(input[curr:curr+4], []byte("null")) {
						value := input[curr:curr+4]

						tokens = append(tokens, Token{Type:Null, Val: string(value)})
						curr += 4
						prevTokenType = Null
						continue
					}
				} else {
					curr++
					prevTokenType = Null
					continue
				}
			}
		}
		prevTokenType = currTokenType
		curr++
		
	}
	if len(stack.TokenTypes) > 0 {
		tokens = append(tokens, Token{Type: ILLEGAL, Val: "Unclosed token"})
		return tokens
	}

	tokens = append(tokens, Token{Type: EOF, Val: ""})
	return tokens
}

func isTerminatingCharacter(c byte) bool {
	return c == ',' || c == '}' || c == ']' || string(c) == EOF.String()
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isBoolean(input []byte, index int) bool {
	trueLiteral := []byte("true")
	falseLiteral := []byte("false")

	if len(input)-index >= len(trueLiteral) {
		if !bytes.Equal(input[index:index+len(trueLiteral)],trueLiteral) {
			return false
		}

		afterTrueLiteral := input[index+len(trueLiteral)]
		if afterTrueLiteral == ',' || afterTrueLiteral == '}' || afterTrueLiteral == ']' {
			return true
		}
	}
	if len(input)-index >= len(falseLiteral) {
		if !bytes.Equal(input[index:index+len(falseLiteral)],falseLiteral) {
			return false
		}

		afterTrueLiteral := input[index+len(trueLiteral)]
		if afterTrueLiteral == ',' || afterTrueLiteral == '}' || afterTrueLiteral == ']' {
			return true
		}
	}
	return false
}

func isNull(input []byte, index int) bool {
	nullLiteral := []byte("null")
	if len(input)-index >= len(nullLiteral) {
		if !bytes.Equal(input[index:index+len(nullLiteral)], nullLiteral) {
			return false
		}

		afterTrueLiteral := input[index+len(nullLiteral)]
		if afterTrueLiteral == ',' || afterTrueLiteral == '}' || afterTrueLiteral == ']' {
			return true
		}
	}
	return false
}

func determineTokenType(char byte, input []byte, currIndex int) Type {
	switch char {
	case '{':
		return LeftBrace
	case '}':
		return RightBrace
	case '[':
		return LeftBracket
	case ']':
		return RightBracket
	case ',':
		return Comma
	case ':':
		return Colon
	case '"':
		return Quote 
	default:
		if unicode.IsDigit(rune(char)) {
			return Number
		} else if char == 't' || char == 'f' {
			if isBoolean(input, currIndex) {
				return Boolean
			}
		} else if char == 'n' {
			if isNull(input, currIndex) {
				return Null
			}
		}
	}
	return ILLEGAL
}

type Stack struct {
	TokenTypes []Type
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(t Type) {
	s.TokenTypes = append(s.TokenTypes, t)
}

func (s *Stack) Peek() Type {
	return s.TokenTypes[len(s.TokenTypes)-1]
}

func (s *Stack) Pop() {
	s.TokenTypes = s.TokenTypes[:len(s.TokenTypes)-1]
}

func isValidSequences(prevToken, currToken Type) bool {
	validSequences := map[Type][]Type{
		ILLEGAL: {String,Number,Boolean,Null,LeftBrace,LeftBracket,Quote},
		LeftBrace: {String, Number,Boolean,Null,LeftBrace,LeftBracket,RightBrace,RightBracket,Quote},
		RightBrace: {Comma,EOF},
		LeftBracket: {String,Number,Boolean,Null,LeftBrace,LeftBracket,RightBracket,Quote},
		RightBracket: {Comma,EOF},
		Comma: {String,Number,Boolean,Null,LeftBrace,LeftBracket,Quote},
		Colon:        {String, Number, Boolean, Null, LeftBrace, LeftBracket, Quote},
		String: {Comma, RightBrace, RightBracket, Colon},
		Number: {Comma, RightBrace, RightBracket},
		Boolean: {Comma, RightBrace, RightBracket},
		Null: {Comma, RightBrace, RightBracket},
		Quote: {Comma, RightBrace, RightBracket, Colon},
	}

	if nextTokens, ok := validSequences[prevToken]; ok {
		return ContainsInArrays(nextTokens, currToken)
	}
	return false
}

func ContainsInArrays(arr []Type, val Type) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}