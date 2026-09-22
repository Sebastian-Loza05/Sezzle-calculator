package calculator

import (
	"errors"
	"math"
)

var (
	ErrInvalidOperation               = errors.New("invalid operation")
	ErrDivisionByZero                 = errors.New("division by zero")
	ErrNonFiniteOperand               = errors.New("operands must be finite numbers")
	ErrNonFiniteResult                = errors.New("result is outside the supported range")
	ErrNegativeSquareRoot             = errors.New("cannot take the square root of a negative number")
	ErrNegativeBaseFractionalExponent = errors.New("a negative base requires an integer exponent")
	ErrZeroToNegativePower            = errors.New("zero cannot be raised to a negative power")
)

func Calculate(a, b float64, operation string) (float64, error) {
	if !isFinite(a) || !isFinite(b) {
		return 0, ErrNonFiniteOperand
	}

	var result float64
	var err error
	switch operation {
	case "add":
		result = Add(a, b)
	case "subtract":
		result = Subtract(a, b)
	case "multiply":
		result = Multiply(a, b)
	case "divide":
		result, err = Divide(a, b)
	case "power":
		result, err = Power(a, b)
	case "sqrt":
		result, err = SquareRoot(a)
	case "percentage":
		result = Percentage(a, b)
	default:
		return 0, ErrInvalidOperation
	}
	if err != nil {
		return 0, err
	}
	if !isFinite(result) {
		return 0, ErrNonFiniteResult
	}
	return result, nil
}

func Add(a, b float64) float64 {
	return a + b
}

func Subtract(a, b float64) float64 {
	return a - b
}

func Multiply(a, b float64) float64 {
	return a * b
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}

	return a / b, nil
}

func Power(a, b float64) (float64, error) {
	if a == 0 && b < 0 {
		return 0, ErrZeroToNegativePower
	}
	if a < 0 && math.Trunc(b) != b {
		return 0, ErrNegativeBaseFractionalExponent
	}
	return math.Pow(a, b), nil
}

func SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSquareRoot
	}
	return math.Sqrt(a), nil
}

func Percentage(a, b float64) float64 {
	return (a / 100) * b
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
