package calculator

import (
	"fmt"
)

// Parser converts an infix slice of tokens into Reverse Polish Notation (RPN) using Shunting-Yard.
type Parser struct {
	tokens []Token
}

// NewParser creates a new Parser instance.
func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

// Parse converts infix tokens into an RPN token queue.
func (p *Parser) Parse() ([]Token, error) {
	if len(p.tokens) == 0 {
		return nil, &CalculationError{
			Message:  "No tokens to parse",
			Code:     ErrCodeEmptyExpr,
			Position: 0,
		}
	}

	var output []Token
	var opStack []Token

	for i := 0; i < len(p.tokens); i++ {
		token := p.tokens[i]

		switch token.Type {
		case TokenNumber:
			output = append(output, token)

		case TokenPercentage:
			// Postfix unary operator: applies directly to preceding operand
			output = append(output, token)

		case TokenUnaryMinus, TokenUnaryPlus:
			// Prefix unary operators: push directly without popping binary operators to the left
			opStack = append(opStack, token)

		case TokenSqrt:
			// Square root function / prefix operator
			opStack = append(opStack, token)

		case TokenPlus, TokenMinus, TokenMultiply, TokenDivide, TokenPower:
			prec := precedence(token.Type)
			rightAssoc := isRightAssociative(token.Type)

			for len(opStack) > 0 {
				top := opStack[len(opStack)-1]
				if top.Type == TokenLParen {
					break
				}
				topPrec := precedence(top.Type)
				if (topPrec > prec) || (topPrec == prec && !rightAssoc) {
					output = append(output, top)
					opStack = opStack[:len(opStack)-1]
				} else {
					break
				}
			}
			opStack = append(opStack, token)

		case TokenLParen:
			opStack = append(opStack, token)

		case TokenRParen:
			matched := false
			for len(opStack) > 0 {
				top := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]
				if top.Type == TokenLParen {
					matched = true
					break
				}
				output = append(output, top)
			}
			if !matched {
				return nil, &CalculationError{
					Message:  fmt.Sprintf("Mismatched closing parenthesis ')' at position %d", token.Position),
					Code:     ErrCodeUnmatchedParen,
					Position: token.Position,
				}
			}
			// If top of stack is now a function like sqrt, pop it to output
			if len(opStack) > 0 && opStack[len(opStack)-1].Type == TokenSqrt {
				output = append(output, opStack[len(opStack)-1])
				opStack = opStack[:len(opStack)-1]
			}

		default:
			return nil, &CalculationError{
				Message:  fmt.Sprintf("Unknown token type '%s' at position %d", token.Type, token.Position),
				Code:     ErrCodeInvalidSyntax,
				Position: token.Position,
			}
		}
	}

	// Drain remaining operators from stack
	for len(opStack) > 0 {
		top := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]
		if top.Type == TokenLParen {
			return nil, &CalculationError{
				Message:  fmt.Sprintf("Unmatched opening parenthesis '(' at position %d", top.Position),
				Code:     ErrCodeUnmatchedParen,
				Position: top.Position,
			}
		}
		output = append(output, top)
	}

	return output, nil
}

func precedence(t TokenType) int {
	switch t {
	case TokenPlus, TokenMinus:
		return 1
	case TokenMultiply, TokenDivide:
		return 2
	case TokenUnaryMinus, TokenUnaryPlus:
		return 3
	case TokenPower:
		return 4
	case TokenSqrt:
		return 5
	default:
		return 0
	}
}

func isRightAssociative(t TokenType) bool {
	switch t {
	case TokenPower, TokenUnaryMinus, TokenUnaryPlus, TokenSqrt:
		return true
	default:
		return false
	}
}
