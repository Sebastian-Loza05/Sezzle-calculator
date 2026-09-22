package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		a, b      float64
		operation string
		want      float64
		wantErr   error
	}{
		{name: "add", a: 2, b: 3, operation: "add", want: 5},
		{name: "subtract", a: 2, b: 3, operation: "subtract", want: -1},
		{name: "multiply", a: 2, b: 3, operation: "multiply", want: 6},
		{name: "divide", a: 10, b: 2, operation: "divide", want: 5},
		{name: "integer power", a: -2, b: 3, operation: "power", want: -8},
		{name: "fractional power", a: 9, b: 0.5, operation: "power", want: 3},
		{name: "zero to zero power", operation: "power", want: 1},
		{name: "square root", a: 9, operation: "sqrt", want: 3},
		{name: "percentage", a: 20, b: 50, operation: "percentage", want: 10},
		{name: "percentage over 100", a: 150, b: 80, operation: "percentage", want: 120},
		{name: "invalid operation", operation: "square", wantErr: ErrInvalidOperation},
		{name: "division by zero", a: 4, operation: "divide", wantErr: ErrDivisionByZero},
		{name: "negative square root", a: -4, operation: "sqrt", wantErr: ErrNegativeSquareRoot},
		{name: "negative base fractional power", a: -9, b: 0.5, operation: "power", wantErr: ErrNegativeBaseFractionalExponent},
		{name: "zero to negative power", a: 0, b: -1, operation: "power", wantErr: ErrZeroToNegativePower},
		{name: "infinite operand", a: math.Inf(1), b: 2, operation: "add", wantErr: ErrNonFiniteOperand},
		{name: "NaN operand", a: math.NaN(), b: 2, operation: "add", wantErr: ErrNonFiniteOperand},
		{name: "overflowing result", a: math.MaxFloat64, b: 2, operation: "multiply", wantErr: ErrNonFiniteResult},
		{name: "overflowing power", a: 1e308, b: 2, operation: "power", wantErr: ErrNonFiniteResult},
		{name: "overflowing percentage", a: 200, b: math.MaxFloat64, operation: "percentage", wantErr: ErrNonFiniteResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.a, tt.b, tt.operation)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Calculate() error = %v, want %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
