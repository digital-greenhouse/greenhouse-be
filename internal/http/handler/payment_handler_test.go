package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
	"digital-greenhouse/greenhouse-be/internal/http/dto"
	"digital-greenhouse/greenhouse-be/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

type mockPaymentServiceForHandler struct {
	ProcessPaymentProofFunc func(ctx context.Context, bookingID uint, amount float64, method domain.PaymentMethod, proofData string, mimeType string) (*domain.Payment, error)
	VerifyPaymentFunc       func(ctx context.Context, paymentID uint, verifierID uint, status domain.PaymentStatus, reason string) error
	GetPaymentProofFunc     func(ctx context.Context, paymentID uint, requesterID uint) (*domain.Payment, error)
}

func (m *mockPaymentServiceForHandler) ProcessPaymentProof(ctx context.Context, bookingID uint, amount float64, method domain.PaymentMethod, proofData string, mimeType string) (*domain.Payment, error) {
	if m.ProcessPaymentProofFunc != nil {
		return m.ProcessPaymentProofFunc(ctx, bookingID, amount, method, proofData, mimeType)
	}
	return nil, nil
}

func (m *mockPaymentServiceForHandler) VerifyPayment(ctx context.Context, paymentID uint, verifierID uint, status domain.PaymentStatus, reason string) error {
	if m.VerifyPaymentFunc != nil {
		return m.VerifyPaymentFunc(ctx, paymentID, verifierID, status, reason)
	}
	return nil
}

func (m *mockPaymentServiceForHandler) GetPaymentProof(ctx context.Context, paymentID uint, requesterID uint) (*domain.Payment, error) {
	if m.GetPaymentProofFunc != nil {
		return m.GetPaymentProofFunc(ctx, paymentID, requesterID)
	}
	return nil, nil
}

func TestPaymentHandler_UploadProof(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        string
		mockSetup      func(m *mockPaymentServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:    "Happy Path",
			reqBody: `{"booking_id":1,"amount":150.0,"payment_method":"TRANSFERENCIA","proof_data":"data:image/png;base64,iVBORw0KGgo=","proof_mime_type":"image/png"}`,
			mockSetup: func(m *mockPaymentServiceForHandler) {
				m.ProcessPaymentProofFunc = func(ctx context.Context, bookingID uint, amount float64, method domain.PaymentMethod, proofData string, mimeType string) (*domain.Payment, error) {
					if bookingID != 1 || amount != 150.0 || method != domain.PaymentMethodTransfer || proofData != "data:image/png;base64,iVBORw0KGgo=" || mimeType != "image/png" {
						return nil, errors.New("invalid arguments to service")
					}
					return &domain.Payment{
						ID:            10,
						BookingID:     bookingID,
						Amount:        amount,
						PaymentMethod: method,
						Status:        domain.PaymentPending,
						PaymentDate:   time.Unix(1600000000, 0),
					}, nil
				}
			},
			expectedStatus: http.StatusCreated,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.PaymentResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if resp.ID != 10 || resp.Status != domain.PaymentPending {
					t.Errorf("unexpected response: %+v", resp)
				}
			},
		},
		{
			name:           "Invalid JSON",
			reqBody:        `{"booking_id":`,
			mockSetup:      func(m *mockPaymentServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("payload inválido")) {
					t.Errorf("expected payload invalid, got: %s", body)
				}
			},
		},
		{
			name:    "Service Error",
			reqBody: `{"booking_id":1,"amount":150.0,"payment_method":"TRANSFERENCIA"}`,
			mockSetup: func(m *mockPaymentServiceForHandler) {
				m.ProcessPaymentProofFunc = func(ctx context.Context, bookingID uint, amount float64, method domain.PaymentMethod, proofData string, mimeType string) (*domain.Payment, error) {
					return nil, errors.New("booking not found")
				}
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("booking not found")) {
					t.Errorf("expected booking not found error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockPaymentServiceForHandler{}
			tt.mockSetup(m)
			h := NewPaymentHandler(m)

			req := httptest.NewRequest(http.MethodPost, "/payments/proof", bytes.NewBufferString(tt.reqBody))
			rec := httptest.NewRecorder()

			h.UploadProof(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestPaymentHandler_VerifyPayment(t *testing.T) {
	tests := []struct {
		name           string
		paymentIDParam string
		userID         uint
		reqBody        string
		mockSetup      func(m *mockPaymentServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:           "Happy Path",
			paymentIDParam: "10",
			userID:         99,
			reqBody:        `{"status":"VERIFIED","rejection_reason":""}`,
			mockSetup: func(m *mockPaymentServiceForHandler) {
				m.VerifyPaymentFunc = func(ctx context.Context, paymentID uint, verifierID uint, status domain.PaymentStatus, reason string) error {
					if paymentID != 10 || verifierID != 99 || status != domain.PaymentVerified || reason != "" {
						return errors.New("invalid arguments to service")
					}
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("pago procesado exitosamente")) {
					t.Errorf("expected success message, got: %s", body)
				}
			},
		},
		{
			name:           "Invalid Payment ID",
			paymentIDParam: "abc",
			userID:         99,
			reqBody:        `{}`,
			mockSetup:      func(m *mockPaymentServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de pago inválido")) {
					t.Errorf("expected invalid payment id error, got: %s", body)
				}
			},
		},
		{
			name:           "Invalid JSON",
			paymentIDParam: "10",
			userID:         99,
			reqBody:        `{"status":`,
			mockSetup:      func(m *mockPaymentServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("payload inválido")) {
					t.Errorf("expected payload invalid error, got: %s", body)
				}
			},
		},
		{
			name:           "Unauthorized",
			paymentIDParam: "10",
			userID:         0,
			reqBody:        `{"status":"VERIFIED"}`,
			mockSetup:      func(m *mockPaymentServiceForHandler) {},
			expectedStatus: http.StatusUnauthorized,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("se requiere autenticación")) {
					t.Errorf("expected unauthorized, got: %s", body)
				}
			},
		},
		{
			name:           "Service Error",
			paymentIDParam: "10",
			userID:         99,
			reqBody:        `{"status":"REJECTED","rejection_reason":"Insufficent funds"}`,
			mockSetup: func(m *mockPaymentServiceForHandler) {
				m.VerifyPaymentFunc = func(ctx context.Context, paymentID uint, verifierID uint, status domain.PaymentStatus, reason string) error {
					return errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("db error")) {
					t.Errorf("expected db error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockPaymentServiceForHandler{}
			tt.mockSetup(m)
			h := NewPaymentHandler(m)

			r := chi.NewRouter()
			r.Post("/payments/{id}/verify", h.VerifyPayment)

			req := httptest.NewRequest(http.MethodPost, "/payments/"+tt.paymentIDParam+"/verify", bytes.NewBufferString(tt.reqBody))
			if tt.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestPaymentHandler_DownloadProof(t *testing.T) {
	tests := []struct {
		name           string
		paymentIDParam string
		userID         uint
		mockSetup      func(m *mockPaymentServiceForHandler)
		expectedStatus int
		verifyHeaders  func(t *testing.T, w *httptest.ResponseRecorder)
		verifyResponse func(t *testing.T, body []byte)
	}{
		{
			name:           "Happy Path - JPEG with Base64 prefix and whitespace",
			paymentIDParam: "10",
			userID:         42,
			mockSetup: func(m *mockPaymentServiceForHandler) {
				m.GetPaymentProofFunc = func(ctx context.Context, paymentID uint, requesterID uint) (*domain.Payment, error) {
					if paymentID != 10 || requesterID != 42 {
						return nil, errors.New("invalid arguments to service")
					}
					// base64 for "hello world" is "aGVsbG8gd29ybGQ="
					// we put prefix "data:image/jpeg;base64," and some whitespace/newlines
					return &domain.Payment{
						ID:            10,
						ProofMimeType: "image/jpeg",
						ProofData:     "data:image/jpeg;base64, aGVsbG8  gd29ybGQ=\n",
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyHeaders: func(t *testing.T, w *httptest.ResponseRecorder) {
				if w.Header().Get("Content-Type") != "image/jpeg" {
					t.Errorf("unexpected content type: %s", w.Header().Get("Content-Type"))
				}
				if w.Header().Get("Content-Disposition") != "attachment; filename=comprobante_10.jpg" {
					t.Errorf("unexpected content disposition: %s", w.Header().Get("Content-Disposition"))
				}
			},
			verifyResponse: func(t *testing.T, body []byte) {
				if string(body) != "hello world" {
					t.Errorf("expected 'hello world', got: %s", string(body))
				}
			},
		},
		{
			name:           "Happy Path - PNG with raw base64 (no padding)",
			paymentIDParam: "10",
			userID:         42,
			mockSetup: func(m *mockPaymentServiceForHandler) {
				m.GetPaymentProofFunc = func(ctx context.Context, paymentID uint, requesterID uint) (*domain.Payment, error) {
					// base64 raw standard (no padding) for "hello world" is "aGVsbG8gd29ybGQ"
					return &domain.Payment{
						ID:            10,
						ProofMimeType: "image/png",
						ProofData:     "aGVsbG8gd29ybGQ",
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyHeaders: func(t *testing.T, w *httptest.ResponseRecorder) {
				if w.Header().Get("Content-Type") != "image/png" {
					t.Errorf("unexpected content type: %s", w.Header().Get("Content-Type"))
				}
				if w.Header().Get("Content-Disposition") != "attachment; filename=comprobante_10.png" {
					t.Errorf("unexpected content disposition: %s", w.Header().Get("Content-Disposition"))
				}
			},
			verifyResponse: func(t *testing.T, body []byte) {
				if string(body) != "hello world" {
					t.Errorf("expected 'hello world', got: %s", string(body))
				}
			},
		},
		{
			name:           "Happy Path - PDF",
			paymentIDParam: "10",
			userID:         42,
			mockSetup: func(m *mockPaymentServiceForHandler) {
				m.GetPaymentProofFunc = func(ctx context.Context, paymentID uint, requesterID uint) (*domain.Payment, error) {
					return &domain.Payment{
						ID:            10,
						ProofMimeType: "application/pdf",
						ProofData:     "aGVsbG8gd29ybGQ=",
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyHeaders: func(t *testing.T, w *httptest.ResponseRecorder) {
				if w.Header().Get("Content-Disposition") != "attachment; filename=comprobante_10.pdf" {
					t.Errorf("unexpected content disposition: %s", w.Header().Get("Content-Disposition"))
				}
			},
		},
		{
			name:           "Happy Path - WebP",
			paymentIDParam: "10",
			userID:         42,
			mockSetup: func(m *mockPaymentServiceForHandler) {
				m.GetPaymentProofFunc = func(ctx context.Context, paymentID uint, requesterID uint) (*domain.Payment, error) {
					return &domain.Payment{
						ID:            10,
						ProofMimeType: "image/webp",
						ProofData:     "aGVsbG8gd29ybGQ=",
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyHeaders: func(t *testing.T, w *httptest.ResponseRecorder) {
				if w.Header().Get("Content-Disposition") != "attachment; filename=comprobante_10.webp" {
					t.Errorf("unexpected content disposition: %s", w.Header().Get("Content-Disposition"))
				}
			},
		},
		{
			name:           "Invalid Payment ID",
			paymentIDParam: "abc",
			userID:         42,
			mockSetup:      func(m *mockPaymentServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body []byte) {
				if !bytes.Contains(body, []byte("ID de pago inválido")) {
					t.Errorf("expected invalid payment id, got: %s", string(body))
				}
			},
		},
		{
			name:           "Forbidden / Service Error",
			paymentIDParam: "10",
			userID:         42,
			mockSetup: func(m *mockPaymentServiceForHandler) {
				m.GetPaymentProofFunc = func(ctx context.Context, paymentID uint, requesterID uint) (*domain.Payment, error) {
					return nil, errors.New("not allowed")
				}
			},
			expectedStatus: http.StatusForbidden,
			verifyResponse: func(t *testing.T, body []byte) {
				if !bytes.Contains(body, []byte("not allowed")) {
					t.Errorf("expected forbidden/not allowed message, got: %s", string(body))
				}
			},
		},
		{
			name:           "Base64 Decoding Failure",
			paymentIDParam: "10",
			userID:         42,
			mockSetup: func(m *mockPaymentServiceForHandler) {
				m.GetPaymentProofFunc = func(ctx context.Context, paymentID uint, requesterID uint) (*domain.Payment, error) {
					return &domain.Payment{
						ID:            10,
						ProofMimeType: "image/jpeg",
						ProofData:     "!!!!invalid-base64!!!!", // invalid base64 characters that fail both standard & raw
					}, nil
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body []byte) {
				if !bytes.Contains(body, []byte("error decodificando comprobante")) {
					t.Errorf("expected decoding error message, got: %s", string(body))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockPaymentServiceForHandler{}
			tt.mockSetup(m)
			h := NewPaymentHandler(m)

			r := chi.NewRouter()
			r.Get("/payments/{id}/proof", h.DownloadProof)

			req := httptest.NewRequest(http.MethodGet, "/payments/"+tt.paymentIDParam+"/proof", nil)
			if tt.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyHeaders != nil {
				tt.verifyHeaders(t, rec)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.Bytes())
			}
		})
	}
}
