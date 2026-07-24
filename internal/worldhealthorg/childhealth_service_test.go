package worldhealthorg_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	. "candybar/internal/worldhealthorg"
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

func TestFetchInfantNutrition(t *testing.T) {
	tests := []struct {
		name       string
		transport  roundTripFunc
		country    string
		wantErr    string
		assertions func(t *testing.T, got *InfantNutrition)
	}{
		{
			name: "returns wrapped error when request fails",
			transport: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network down")
			},
			country: "MEX",
			wantErr: "error on getting Infant Nutrition",
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
			country: "MEX",
			wantErr: "read failed",
		},
		{
			name: "parses valid json payload",
			transport: func(req *http.Request) (*http.Response, error) {
				if !strings.Contains(req.URL.String(), "COUNTRY:MEX") {
					t.Fatalf("expected country filter in URL, got %q", req.URL.String())
				}

				body := `{
					"Copyright": "WHO",
					"Dataset": [{"Label": "WHS_PBR", "Display": "Dataset Display"}],
					"Attribute": [],
					"Dimension": [],
					"Fact": []
				}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			},
			country: "MEX",
			assertions: func(t *testing.T, got *InfantNutrition) {
				if got == nil {
					t.Fatalf("expected non-nil result")
				}
				if got.Copyright != "WHO" {
					t.Fatalf("expected Copyright to be WHO, got %q", got.Copyright)
				}
				if len(got.Dataset) != 1 {
					t.Fatalf("expected 1 dataset, got %d", len(got.Dataset))
				}
				if got.Dataset[0].Label != "WHS_PBR" {
					t.Fatalf("unexpected dataset label: %q", got.Dataset[0].Label)
				}
			},
		},
		{
			name: "invalid json returns zero value without error",
			transport: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{invalid-json")),
					Header:     make(http.Header),
				}, nil
			},
			country: "MEX",
			assertions: func(t *testing.T, got *InfantNutrition) {
				if got == nil {
					t.Fatalf("expected non-nil result")
				}
				if got.Copyright != "" || len(got.Dataset) != 0 {
					t.Fatalf("expected zero-value result when unmarshal fails, got %+v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			withMockTransport(t, tt.transport)

			result, err := FetchInfantNutrition(tt.country)

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
