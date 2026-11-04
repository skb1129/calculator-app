package calculator

import (
	"fmt"
	"math"
	"strconv"
)

// CalculationResult holds the computed value and metadata.
type CalculationResult struct {
	Expression string   `json:"expression"`
	Result     float64  `json:"result"`
	Formatted  string   `json:"formatted"`
	Tokens     []string `json:"tokens,omitempty"`
}

// Evaluate evaluates an RPN token queue and returns the scalar result.
func Evaluate(rpn []Token) (float64, string, error) {
	if len(rpn) == 0 {
		return 0, "", &CalculationError{
			Message:  "Empty expression queue",
			Code:     ErrCodeEmptyExpr,
			Position: 0,
		}
	}

	var stack []float64

	for _, token := range rpn {
		switch token.Type {
		case TokenNumber:
			stack = append(stack, token.Number)

		case TokenPercentage:
			if len(stack) < 1 {
				return 0, "", missingOperandError("%", token.Position)
			}
			stack[len(stack)-1] /= 100.0

		case TokenUnaryMinus:
			if len(stack) < 1 {
				return 0, "", missingOperandError("-", token.Position)
			}
			stack[len(stack)-1] = -stack[len(stack)-1]

		case TokenUnaryPlus:
			if len(stack) < 1 {
				return 0, "", missingOperandError("+", token.Position)
			}
			// Unary plus is identity: +a = a

		case TokenSqrt:
			if len(stack) < 1 {
				return 0, "", missingOperandError("sqrt", token.Position)
			}
			a := stack[len(stack)-1]
			if a < 0 {
				return 0, "", &CalculationError{
					Message:  fmt.Sprintf("Cannot calculate square root of negative number (%v) at position %d", a, token.Position),
					Code:     ErrCodeMathDomain,
					Position: token.Position,
				}
			}
			stack[len(stack)-1] = math.Sqrt(a)

		case TokenPlus, TokenMinus, TokenMultiply, TokenDivide, TokenPower:
			if len(stack) < 2 {
				return 0, "", missingOperandError(token.Value, token.Position)
			}

			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			res, err := applyBinaryOp(token.Type, a, b, token.Position)
			if err != nil {
				return 0, "", err
			}
			stack = append(stack, res)

		default:
			return 0, "", &CalculationError{
				Message:  fmt.Sprintf("Unsupported execution token '%s' at position %d", token.Type, token.Position),
				Code:     ErrCodeInvalidSyntax,
				Position: token.Position,
			}
		}
	}

	if len(stack) != 1 {
		return 0, "", &CalculationError{
			Message:  "Invalid expression syntax: redundant operands",
			Code:     ErrCodeInvalidSyntax,
			Position: 0,
		}
	}

	cleaned, formatted := FormatResult(stack[0])
	return cleaned, formatted, nil
}

func missingOperandError(op string, pos int) *CalculationError {
	return &CalculationError{
		Message:  fmt.Sprintf("Missing operand for operator '%s' at position %d", op, pos),
		Code:     ErrCodeInvalidSyntax,
		Position: pos,
	}
}

func applyBinaryOp(op TokenType, a, b float64, pos int) (float64, error) {
	switch op {
	case TokenPlus:
		return a + b, nil

	case TokenMinus:
		return a - b, nil

	case TokenMultiply:
		return a * b, nil

	case TokenDivide:
		if math.Abs(b) < 1e-15 {
			return 0, &CalculationError{
				Message:  fmt.Sprintf("Division by zero at position %d", pos),
				Code:     ErrCodeDivisionByZero,
				Position: pos,
			}
		}
		return a / b, nil

	case TokenPower:
		if a == 0 && b < 0 {
			return 0, &CalculationError{
				Message:  fmt.Sprintf("Zero cannot be raised to negative power at position %d", pos),
				Code:     ErrCodeDivisionByZero,
				Position: pos,
			}
		}
		if a < 0 && b != math.Trunc(b) {
			return 0, &CalculationError{
				Message:  fmt.Sprintf("Negative base (%v) cannot be raised to fractional power (%v) at position %d", a, b, pos),
				Code:     ErrCodeMathDomain,
				Position: pos,
			}
		}

		res := math.Pow(a, b)
		if math.IsNaN(res) {
			return 0, &CalculationError{
				Message:  fmt.Sprintf("Mathematical operation produced undefined result (NaN) at position %d", pos),
				Code:     ErrCodeMathDomain,
				Position: pos,
			}
		}
		if math.IsInf(res, 0) {
			return 0, &CalculationError{
				Message:  fmt.Sprintf("Calculation result overflowed at position %d", pos),
				Code:     ErrCodeMathOverflow,
				Position: pos,
			}
		}
		return res, nil

	default:
		return 0, &CalculationError{
			Message:  fmt.Sprintf("Unsupported operator '%s' at position %d", op, pos),
			Code:     ErrCodeInvalidSyntax,
			Position: pos,
		}
	}
}

// FormatResult rounds float values to prevent 64-bit IEEE representation artifacts (e.g. 0.1 + 0.2).
func FormatResult(val float64) (float64, string) {
	rounded := math.Round(val*1e12) / 1e12
	if math.Abs(rounded) == 0 {
		rounded = 0
	}
	formatted := strconv.FormatFloat(rounded, 'g', 12, 64)
	return rounded, formatted
}

// Calculate parses and evaluates a raw expression string.
func Calculate(expression string) (*CalculationResult, error) {
	tokens, err := NewLexer(expression).Tokenize()
	if err != nil {
		return nil, err
	}

	rpn, err := NewParser(tokens).Parse()
	if err != nil {
		return nil, err
	}

	result, formatted, err := Evaluate(rpn)
	if err != nil {
		return nil, err
	}

	tokenStrings := make([]string, len(tokens))
	for i, t := range tokens {
		tokenStrings[i] = t.Value
	}

	return &CalculationResult{
		Expression: expression,
		Result:     result,
		Formatted:  formatted,
		Tokens:     tokenStrings,
	}, nil
}
