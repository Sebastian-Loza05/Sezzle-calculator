package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculateRoute(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		body        string
		wantStatus  int
		wantCode    string
		wantResult  float64
	}{
		{name: "success", method: http.MethodPost, contentType: "application/json", body: `{"operation":"add","a":2,"b":3}`, wantStatus: http.StatusOK, wantResult: 5},
		{name: "zero is present", method: http.MethodPost, contentType: "application/json", body: `{"operation":"subtract","a":0,"b":3}`, wantStatus: http.StatusOK, wantResult: -3},
		{name: "power with fractional exponent", method: http.MethodPost, contentType: "application/json", body: `{"operation":"power","a":9,"b":0.5}`, wantStatus: http.StatusOK, wantResult: 3},
		{name: "square root without b", method: http.MethodPost, contentType: "application/json", body: `{"operation":"sqrt","a":9}`, wantStatus: http.StatusOK, wantResult: 3},
		{name: "percentage is a percent of b", method: http.MethodPost, contentType: "application/json", body: `{"operation":"percentage","a":20,"b":50}`, wantStatus: http.StatusOK, wantResult: 10},
		{name: "division by zero", method: http.MethodPost, contentType: "application/json", body: `{"operation":"divide","a":3,"b":0}`, wantStatus: http.StatusBadRequest, wantCode: "division_by_zero"},
		{name: "negative square root", method: http.MethodPost, contentType: "application/json", body: `{"operation":"sqrt","a":-9}`, wantStatus: http.StatusBadRequest, wantCode: "negative_square_root"},
		{name: "square root rejects b", method: http.MethodPost, contentType: "application/json", body: `{"operation":"sqrt","a":9,"b":2}`, wantStatus: http.StatusBadRequest, wantCode: "unexpected_operand"},
		{name: "negative base fractional exponent", method: http.MethodPost, contentType: "application/json", body: `{"operation":"power","a":-9,"b":0.5}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_exponent"},
		{name: "zero to negative power", method: http.MethodPost, contentType: "application/json", body: `{"operation":"power","a":0,"b":-1}`, wantStatus: http.StatusBadRequest, wantCode: "zero_to_negative_power"},
		{name: "missing operand", method: http.MethodPost, contentType: "application/json", body: `{"operation":"add","a":2}`, wantStatus: http.StatusBadRequest, wantCode: "missing_operand"},
		{name: "missing square root operand", method: http.MethodPost, contentType: "application/json", body: `{"operation":"sqrt"}`, wantStatus: http.StatusBadRequest, wantCode: "missing_operand"},
		{name: "missing operation", method: http.MethodPost, contentType: "application/json", body: `{"a":2,"b":3}`, wantStatus: http.StatusBadRequest, wantCode: "missing_operation"},
		{name: "invalid operation", method: http.MethodPost, contentType: "application/json", body: `{"operation":"square","a":2,"b":3}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_operation"},
		{name: "unknown field", method: http.MethodPost, contentType: "application/json", body: `{"operation":"add","a":2,"b":3,"extra":1}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_request"},
		{name: "trailing JSON", method: http.MethodPost, contentType: "application/json", body: `{"operation":"add","a":2,"b":3} {}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_request"},
		{name: "wrong operand type", method: http.MethodPost, contentType: "application/json", body: `{"operation":"add","a":"2","b":3}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_request"},
		{name: "wrong content type", method: http.MethodPost, contentType: "text/plain", body: `{"operation":"add","a":2,"b":3}`, wantStatus: http.StatusUnsupportedMediaType, wantCode: "unsupported_media_type"},
		{name: "wrong method", method: http.MethodGet, wantStatus: http.StatusMethodNotAllowed, wantCode: "method_not_allowed"},
		{name: "oversized body", method: http.MethodPost, contentType: "application/json", body: strings.Repeat(" ", 8<<10) + `{}`, wantStatus: http.StatusRequestEntityTooLarge, wantCode: "request_too_large"},
		{name: "overflowing result", method: http.MethodPost, contentType: "application/json", body: `{"operation":"multiply","a":1e308,"b":2}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_result"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/calculate", strings.NewReader(tt.body))
			if tt.contentType != "" {
				request.Header.Set("Content-Type", tt.contentType)
			}
			recorder := httptest.NewRecorder()
			NewRouter().ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
			if got := recorder.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", got)
			}
			if tt.method == http.MethodGet && recorder.Header().Get("Allow") != "POST" {
				t.Fatalf("Allow = %q, want POST", recorder.Header().Get("Allow"))
			}
			if tt.wantCode == "" {
				var response Response
				if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if response.Result != tt.wantResult {
					t.Fatalf("result = %v, want %v", response.Result, tt.wantResult)
				}
				return
			}
			var response ErrorResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("decode error response: %v", err)
			}
			if response.Error.Code != tt.wantCode {
				t.Fatalf("error code = %q, want %q", response.Error.Code, tt.wantCode)
			}
		})
	}
}
