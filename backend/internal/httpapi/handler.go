package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"github.com/Sebastian-Loza05/Sezzle-calculator/backend/internal/calculator"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type must be application/json")
		return
	}

	var req Request
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeDecodeError(w, err)
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			writeError(w, http.StatusBadRequest, "invalid_request", "request body must contain one JSON object")
		} else {
			writeDecodeError(w, err)
		}
		return
	}
	if req.Operation == "" {
		writeError(w, http.StatusBadRequest, "missing_operation", "operation is required")
		return
	}
	if req.A == nil {
		writeError(w, http.StatusBadRequest, "missing_operand", "a is required")
		return
	}
	var b float64
	if req.Operation == "sqrt" {
		if req.B != nil {
			writeError(w, http.StatusBadRequest, "unexpected_operand", "b is not used for square root")
			return
		}
	} else {
		if req.B == nil {
			writeError(w, http.StatusBadRequest, "missing_operand", "b is required")
			return
		}
		b = *req.B
	}

	result, err := calculator.Calculate(*req.A, b, req.Operation)
	if err != nil {
		switch {
		case errors.Is(err, calculator.ErrInvalidOperation):
			writeError(w, http.StatusBadRequest, "invalid_operation", err.Error())
		case errors.Is(err, calculator.ErrDivisionByZero):
			writeError(w, http.StatusBadRequest, "division_by_zero", err.Error())
		case errors.Is(err, calculator.ErrNonFiniteOperand):
			writeError(w, http.StatusBadRequest, "invalid_operand", err.Error())
		case errors.Is(err, calculator.ErrNonFiniteResult):
			writeError(w, http.StatusBadRequest, "invalid_result", err.Error())
		case errors.Is(err, calculator.ErrNegativeSquareRoot):
			writeError(w, http.StatusBadRequest, "negative_square_root", err.Error())
		case errors.Is(err, calculator.ErrNegativeBaseFractionalExponent):
			writeError(w, http.StatusBadRequest, "invalid_exponent", err.Error())
		case errors.Is(err, calculator.ErrZeroToNegativePower):
			writeError(w, http.StatusBadRequest, "zero_to_negative_power", err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, Response{Result: result})
}
