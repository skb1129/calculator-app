package calculator_test

import (
	"testing"

	"github.com/skb1129/calculator-app/internal/calculator"
)

func tokenizeHelper(t *testing.T, expr string) []calculator.Token {
	t.Helper()
	lexer := calculator.NewLexer(expr)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("failed to tokenize %q: %v", expr, err)
	}
	return tokens
}

func tokenTypesToStrings(tokens []calculator.Token) []string {
	res := make([]string, len(tokens))
	for i, tok := range tokens {
		if tok.Type == calculator.TokenNumber {
			res[i] = tok.Value
		} else {
			res[i] = string(tok.Type)
		}
	}
	return res
}

func TestParserEmptyTokens(t *testing.T) {
	parser := calculator.NewParser([]calculator.Token{})
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("expected error parsing empty tokens, got nil")
	}

	calcErr, ok := err.(*calculator.CalculationError)
	if !ok {
		t.Fatalf("expected *CalculationError, got %T", err)
	}
	if calcErr.Code != calculator.ErrCodeEmptyExpr {
		t.Errorf("expected code %s, got %s", calculator.ErrCodeEmptyExpr, calcErr.Code)
	}
}

func TestParserPrecedenceAndAssociativity(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expectedRPN []string
	}{
		{
			name:        "Simple addition",
			expression:  "2 + 3",
			expectedRPN: []string{"2", "3", string(calculator.TokenPlus)},
		},
		{
			name:        "Multiplication before addition",
			expression:  "2 + 3 * 4",
			expectedRPN: []string{"2", "3", "4", string(calculator.TokenMultiply), string(calculator.TokenPlus)},
		},
		{
			name:        "Parentheses override precedence",
			expression:  "(2 + 3) * 4",
			expectedRPN: []string{"2", "3", string(calculator.TokenPlus), "4", string(calculator.TokenMultiply)},
		},
		{
			name:       "Right-associative exponentiation (2 ^ 3 ^ 2)",
			expression: "2 ^ 3 ^ 2",
			// 3 ^ 2 is evaluated first, so RPN is: 2 3 2 ^ ^
			expectedRPN: []string{"2", "3", "2", string(calculator.TokenPower), string(calculator.TokenPower)},
		},
		{
			name:        "Explicit left grouping for exponentiation ((2 ^ 3) ^ 2)",
			expression:  "(2 ^ 3) ^ 2",
			expectedRPN: []string{"2", "3", string(calculator.TokenPower), "2", string(calculator.TokenPower)},
		},
		{
			name:        "Division and subtraction",
			expression:  "20 / 4 - 2",
			expectedRPN: []string{"20", "4", string(calculator.TokenDivide), "2", string(calculator.TokenMinus)},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := tokenizeHelper(t, tc.expression)
			parser := calculator.NewParser(tokens)
			rpn, err := parser.Parse()
			if err != nil {
				t.Fatalf("unexpected error parsing %q: %v", tc.expression, err)
			}

			actualStrings := tokenTypesToStrings(rpn)
			if len(actualStrings) != len(tc.expectedRPN) {
				t.Fatalf("expected RPN length %d, got %d: %v", len(tc.expectedRPN), len(actualStrings), actualStrings)
			}

			for i := range actualStrings {
				if actualStrings[i] != tc.expectedRPN[i] {
					t.Errorf("at index %d: expected %s, got %s (full: %v)", i, tc.expectedRPN[i], actualStrings[i], actualStrings)
				}
			}
		})
	}
}

func TestParserUnaryOperators(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expectedRPN []string
	}{
		{
			name:        "Leading unary minus",
			expression:  "-5 + 3",
			expectedRPN: []string{"5", string(calculator.TokenUnaryMinus), "3", string(calculator.TokenPlus)},
		},
		{
			name:        "Unary minus after multiplication",
			expression:  "4 * -2",
			expectedRPN: []string{"4", "2", string(calculator.TokenUnaryMinus), string(calculator.TokenMultiply)},
		},
		{
			name:        "Unary plus",
			expression:  "+7 * 3",
			expectedRPN: []string{"7", string(calculator.TokenUnaryPlus), "3", string(calculator.TokenMultiply)},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := tokenizeHelper(t, tc.expression)
			parser := calculator.NewParser(tokens)
			rpn, err := parser.Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			actualStrings := tokenTypesToStrings(rpn)
			if len(actualStrings) != len(tc.expectedRPN) {
				t.Fatalf("expected RPN length %d, got %d: %v", len(tc.expectedRPN), len(actualStrings), actualStrings)
			}

			for i := range actualStrings {
				if actualStrings[i] != tc.expectedRPN[i] {
					t.Errorf("at index %d: expected %s, got %s", i, tc.expectedRPN[i], actualStrings[i])
				}
			}
		})
	}
}

func TestParserSquareRootAndPercentage(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expectedRPN []string
	}{
		{
			name:        "Square root with parens",
			expression:  "sqrt(16) + 2",
			expectedRPN: []string{"16", string(calculator.TokenSqrt), "2", string(calculator.TokenPlus)},
		},
		{
			name:        "Square root sub-expression",
			expression:  "sqrt(9 + 16)",
			expectedRPN: []string{"9", "16", string(calculator.TokenPlus), string(calculator.TokenSqrt)},
		},
		{
			name:        "Percentage postfix operator in multiplication",
			expression:  "100 * 15%",
			expectedRPN: []string{"100", "15", string(calculator.TokenPercentage), string(calculator.TokenMultiply)},
		},
		{
			name:        "Percentage after addition",
			expression:  "1000 + 10%",
			expectedRPN: []string{"1000", "10", string(calculator.TokenPercentage), string(calculator.TokenPlus)},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := tokenizeHelper(t, tc.expression)
			parser := calculator.NewParser(tokens)
			rpn, err := parser.Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			actualStrings := tokenTypesToStrings(rpn)
			if len(actualStrings) != len(tc.expectedRPN) {
				t.Fatalf("expected RPN length %d, got %d: %v", len(tc.expectedRPN), len(actualStrings), actualStrings)
			}

			for i := range actualStrings {
				if actualStrings[i] != tc.expectedRPN[i] {
					t.Errorf("at index %d: expected %s, got %s", i, tc.expectedRPN[i], actualStrings[i])
				}
			}
		})
	}
}

func TestParserErrors(t *testing.T) {
	t.Run("Unmatched opening parenthesis", func(t *testing.T) {
		tokens := []calculator.Token{
			{Type: calculator.TokenLParen, Value: "(", Position: 0},
			{Type: calculator.TokenNumber, Value: "2", Number: 2, Position: 1},
			{Type: calculator.TokenPlus, Value: "+", Position: 3},
			{Type: calculator.TokenNumber, Value: "3", Number: 3, Position: 5},
		}
		parser := calculator.NewParser(tokens)
		_, err := parser.Parse()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		calcErr := err.(*calculator.CalculationError)
		if calcErr.Code != calculator.ErrCodeUnmatchedParen {
			t.Errorf("expected code %s, got %s", calculator.ErrCodeUnmatchedParen, calcErr.Code)
		}
	})

	t.Run("Mismatched closing parenthesis", func(t *testing.T) {
		tokens := []calculator.Token{
			{Type: calculator.TokenNumber, Value: "2", Number: 2, Position: 0},
			{Type: calculator.TokenPlus, Value: "+", Position: 2},
			{Type: calculator.TokenNumber, Value: "3", Number: 3, Position: 4},
			{Type: calculator.TokenRParen, Value: ")", Position: 5},
		}
		parser := calculator.NewParser(tokens)
		_, err := parser.Parse()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		calcErr := err.(*calculator.CalculationError)
		if calcErr.Code != calculator.ErrCodeUnmatchedParen {
			t.Errorf("expected code %s, got %s", calculator.ErrCodeUnmatchedParen, calcErr.Code)
		}
	})

	t.Run("Unknown token type", func(t *testing.T) {
		tokens := []calculator.Token{
			{Type: calculator.TokenType("UNKNOWN"), Value: "??", Position: 0},
		}
		parser := calculator.NewParser(tokens)
		_, err := parser.Parse()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		calcErr := err.(*calculator.CalculationError)
		if calcErr.Code != calculator.ErrCodeInvalidSyntax {
			t.Errorf("expected code %s, got %s", calculator.ErrCodeInvalidSyntax, calcErr.Code)
		}
	})
}
