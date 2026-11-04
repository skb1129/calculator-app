package calculator

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Lexer converts an expression string into a slice of tokens.
type Lexer struct {
	input []rune
	pos   int
}

// NewLexer creates a new Lexer instance.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		pos:   0,
	}
}

// Tokenize processes the entire input and returns tokens or a CalculationError.
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token

	if len(strings.TrimSpace(string(l.input))) == 0 {
		return nil, &CalculationError{
			Message:  "Expression cannot be empty",
			Code:     ErrCodeEmptyExpr,
			Position: 0,
		}
	}

	for l.pos < len(l.input) {
		r := l.input[l.pos]

		// Skip whitespace
		if unicode.IsSpace(r) {
			l.pos++
			continue
		}

		startPos := l.pos
		var prevToken *Token
		if len(tokens) > 0 {
			prevToken = &tokens[len(tokens)-1]
		}

		// Numbers (digits or leading dot like .5)
		if unicode.IsDigit(r) || (r == '.' && l.peekNextIsDigit()) {
			numToken, err := l.readNumber()
			if err != nil {
				return nil, err
			}

			// Check for implicit multiplication: e.g. "5(2)" or "5sqrt(4)"
			// But for number, this is handled when the next token is parsed.
			tokens = append(tokens, numToken)
			continue
		}

		// Operators & Parentheses
		switch r {
		case '+':
			if l.isUnaryContext(prevToken) {
				tokens = append(tokens, Token{Type: TokenUnaryPlus, Value: "+", Position: startPos})
			} else {
				tokens = append(tokens, Token{Type: TokenPlus, Value: "+", Position: startPos})
			}
			l.pos++

		case '-':
			if l.isUnaryContext(prevToken) {
				tokens = append(tokens, Token{Type: TokenUnaryMinus, Value: "-", Position: startPos})
			} else {
				tokens = append(tokens, Token{Type: TokenMinus, Value: "-", Position: startPos})
			}
			l.pos++

		case '*', '×':
			tokens = append(tokens, Token{Type: TokenMultiply, Value: "*", Position: startPos})
			l.pos++

		case '/', '÷':
			tokens = append(tokens, Token{Type: TokenDivide, Value: "/", Position: startPos})
			l.pos++

		case '^':
			tokens = append(tokens, Token{Type: TokenPower, Value: "^", Position: startPos})
			l.pos++

		case '%':
			// Percentage is a postfix operator; must follow a number, rparen, or another percentage
			if prevToken == nil || (prevToken.Type != TokenNumber && prevToken.Type != TokenRParen && prevToken.Type != TokenPercentage) {
				return nil, &CalculationError{
					Message:  fmt.Sprintf("Invalid position for percentage operator '%%' at index %d", startPos),
					Code:     ErrCodeInvalidSyntax,
					Position: startPos,
				}
			}
			tokens = append(tokens, Token{Type: TokenPercentage, Value: "%", Position: startPos})
			l.pos++

		case '(':
			// Check for implicit multiplication before '(': e.g. "2(3)" -> "2*(3)" or "(2)(3)" -> "(2)*(3)"
			if prevToken != nil && (prevToken.Type == TokenNumber || prevToken.Type == TokenRParen || prevToken.Type == TokenPercentage) {
				tokens = append(tokens, Token{Type: TokenMultiply, Value: "*", Position: startPos})
			}
			tokens = append(tokens, Token{Type: TokenLParen, Value: "(", Position: startPos})
			l.pos++

		case ')':
			tokens = append(tokens, Token{Type: TokenRParen, Value: ")", Position: startPos})
			l.pos++

		case '√':
			// Unicode square root symbol
			if prevToken != nil && (prevToken.Type == TokenNumber || prevToken.Type == TokenRParen || prevToken.Type == TokenPercentage) {
				tokens = append(tokens, Token{Type: TokenMultiply, Value: "*", Position: startPos})
			}
			tokens = append(tokens, Token{Type: TokenSqrt, Value: "sqrt", Position: startPos})
			l.pos++

		default:
			// Identifier check (e.g. "sqrt")
			if unicode.IsLetter(r) {
				ident := l.readIdentifier()
				lowerIdent := strings.ToLower(ident)
				if lowerIdent == "sqrt" {
					if prevToken != nil && (prevToken.Type == TokenNumber || prevToken.Type == TokenRParen || prevToken.Type == TokenPercentage) {
						tokens = append(tokens, Token{Type: TokenMultiply, Value: "*", Position: startPos})
					}
					tokens = append(tokens, Token{Type: TokenSqrt, Value: "sqrt", Position: startPos})
				} else {
					return nil, &CalculationError{
						Message:  fmt.Sprintf("Unknown identifier '%s' at position %d", ident, startPos),
						Code:     ErrCodeInvalidSyntax,
						Position: startPos,
					}
				}
			} else {
				return nil, &CalculationError{
					Message:  fmt.Sprintf("Unexpected character '%c' at position %d", r, startPos),
					Code:     ErrCodeInvalidSyntax,
					Position: startPos,
				}
			}
		}
	}

	return tokens, nil
}

func (l *Lexer) peekNextIsDigit() bool {
	if l.pos+1 < len(l.input) {
		return unicode.IsDigit(l.input[l.pos+1])
	}
	return false
}

func (l *Lexer) isUnaryContext(prevToken *Token) bool {
	if prevToken == nil {
		return true
	}
	switch prevToken.Type {
	case TokenPlus, TokenMinus, TokenMultiply, TokenDivide, TokenPower,
		TokenLParen, TokenUnaryMinus, TokenUnaryPlus, TokenSqrt:
		return true
	default:
		return false
	}
}

func (l *Lexer) readNumber() (Token, error) {
	startPos := l.pos
	hasDot := false
	hasExp := false

	var sb strings.Builder

	for l.pos < len(l.input) {
		r := l.input[l.pos]

		if unicode.IsDigit(r) {
			sb.WriteRune(r)
			l.pos++
		} else if r == '.' {
			if hasDot {
				return Token{}, &CalculationError{
					Message:  fmt.Sprintf("Invalid multiple decimal points in number at position %d", l.pos),
					Code:     ErrCodeInvalidSyntax,
					Position: l.pos,
				}
			}
			hasDot = true
			sb.WriteRune(r)
			l.pos++
		} else if (r == 'e' || r == 'E') && !hasExp {
			hasExp = true
			sb.WriteRune(r)
			l.pos++

			// Check for optional +/- in scientific notation exponent
			if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
				sb.WriteRune(l.input[l.pos])
				l.pos++
			}
		} else {
			break
		}
	}

	strVal := sb.String()
	num, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		return Token{}, &CalculationError{
			Message:  fmt.Sprintf("Failed to parse number '%s' at position %d", strVal, startPos),
			Code:     ErrCodeInvalidSyntax,
			Position: startPos,
		}
	}

	return Token{
		Type:     TokenNumber,
		Value:    strVal,
		Number:   num,
		Position: startPos,
	}, nil
}

func (l *Lexer) readIdentifier() string {
	var sb strings.Builder
	for l.pos < len(l.input) && unicode.IsLetter(l.input[l.pos]) {
		sb.WriteRune(l.input[l.pos])
		l.pos++
	}
	return sb.String()
}
