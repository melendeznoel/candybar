package drug_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	. "candybar/internal/drug"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type errReadCloser struct{}

func (e errReadCloser) Read(p []byte) (int, error) {
	return 0, errors.New("read failed")
}

func (e errReadCloser) Close() error {
	return nil
}

func withMockTransport(t *testing.T, fn roundTripFunc) {
	t.Helper()

	original := http.DefaultTransport
	http.DefaultTransport = fn
	t.Cleanup(func() {
		http.DefaultTransport = original
	})
}

func TestFetchRecallsByProductDescription(t *testing.T) {
	tests := []struct {
		name       string
		transport  roundTripFunc
		product    string
		wantErr    string
		assertions func(t *testing.T, got *RecallResponse)
	}{
		{
			name: "returns error when request fails",
			transport: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network down")
			},
			product: "ibuprofen",
			wantErr: "Get \"https://api.fda.gov/drug/enforcement.json?limit=100&search=product_description%3Aibuprofen+AND+status%3AOngoing&skip=0\": network down",
		},
		{
			name: "returns read error when body cannot be read",
			transport: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       errReadCloser{},
					Header:     make(http.Header),
				}, nil
			},
			product: "ibuprofen",
			wantErr: "read failed",
		},
		{
			name: "404 returns empty response with no error",
			transport: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"error":{"code":"NOT_FOUND"}}`)),
					Header:     make(http.Header),
				}, nil
			},
			product: "nonexistent drug",
			assertions: func(t *testing.T, got *RecallResponse) {
				if got == nil {
					t.Fatalf("expected non-nil result")
				}
				if len(got.Results) != 0 {
					t.Fatalf("expected 0 results, got %d", len(got.Results))
				}
			},
		},
		{
			name: "non-200/404 status returns error with body",
			transport: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("boom")),
					Header:     make(http.Header),
				}, nil
			},
			product: "ibuprofen",
			wantErr: "openFDA returned status 500: boom",
		},
		{
			name: "invalid json returns unmarshal error",
			transport: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{invalid-json")),
					Header:     make(http.Header),
				}, nil
			},
			product: "ibuprofen",
			wantErr: "invalid character 'i' looking for beginning of object key string",
		},
		{
			name: "parses valid json payload and builds query",
			transport: func(req *http.Request) (*http.Response, error) {
				if !strings.HasPrefix(req.URL.String(), "https://api.fda.gov/drug/enforcement.json") {
					t.Fatalf("expected openFDA drug enforcement URL, got %q", req.URL.String())
				}
				if !strings.Contains(req.URL.String(), "product_description%3Aibuprofen") {
					t.Fatalf("expected product_description filter in URL, got %q", req.URL.String())
				}
				if !strings.Contains(req.URL.String(), "status%3AOngoing") {
					t.Fatalf("expected status:Ongoing filter in URL, got %q", req.URL.String())
				}

				body := `{
					"meta": {
						"disclaimer": "test",
						"results": {"skip": 0, "limit": 100, "total": 1}
					},
					"results": [{
						"recall_number": "D-1234-2026",
						"reason_for_recall": "Mislabeled",
						"status": "Ongoing",
						"city": "Springfield",
						"state": "IL",
						"country": "United States",
						"classification": "Class II",
						"recalling_firm": "Acme Pharma",
						"product_description": "Ibuprofen 200mg Tablets"
					}]
				}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			},
			product: "ibuprofen",
			assertions: func(t *testing.T, got *RecallResponse) {
				if got == nil {
					t.Fatalf("expected non-nil result")
				}
				if len(got.Results) != 1 {
					t.Fatalf("expected 1 result, got %d", len(got.Results))
				}
				if got.Results[0].RecallNumber != "D-1234-2026" {
					t.Fatalf("unexpected RecallNumber: %q", got.Results[0].RecallNumber)
				}
				if got.Results[0].ProductDescription != "Ibuprofen 200mg Tablets" {
					t.Fatalf("unexpected ProductDescription: %q", got.Results[0].ProductDescription)
				}
				if got.Meta.Results.Total != 1 {
					t.Fatalf("unexpected Meta.Results.Total: %d", got.Meta.Results.Total)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			withMockTransport(t, tt.transport)

			result, err := FetchRecallsByProductDescription(tt.product)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				if result != nil {
					t.Fatalf("expected nil result, got %+v", result)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if tt.assertions != nil {
				tt.assertions(t, result)
			}
		})
	}
}

func TestHasRecalls(t *testing.T) {
	tests := []struct {
		name      string
		transport roundTripFunc
		product   string
		want      bool
		wantErr   string
	}{
		{
			name: "returns error when fetch fails",
			transport: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network down")
			},
			product: "ibuprofen",
			wantErr: "Get \"https://api.fda.gov/drug/enforcement.json?limit=100&search=product_description%3Aibuprofen+AND+status%3AOngoing&skip=0\": network down",
		},
		{
			name: "returns false when no results",
			transport: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{}`)),
					Header:     make(http.Header),
				}, nil
			},
			product: "nonexistent drug",
			want:    false,
		},
		{
			name: "returns true when results exist",
			transport: func(req *http.Request) (*http.Response, error) {
				body := `{"results": [{"recall_number": "D-1234-2026"}]}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			},
			product: "ibuprofen",
			want:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			withMockTransport(t, tt.transport)

			got, err := HasRecalls(tt.product)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				if got != false {
					t.Fatalf("expected false result, got %v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
