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
				if !strings.Contains(req.URL.String(), "ghoapi.azureedge.net/api/WHOSIS_000006") {
					t.Fatalf("expected GHO OData indicator URL, got %q", req.URL.String())
				}
				if !strings.Contains(req.URL.String(), "SpatialDim") || !strings.Contains(req.URL.String(), "MEX") {
					t.Fatalf("expected country filter in URL, got %q", req.URL.String())
				}

				body := `{
					"@odata.context": "https://ghoapi.azureedge.net/api/$metadata#WHOSIS_000006",
					"value": [{
						"Id": 1888461,
						"IndicatorCode": "WHOSIS_000006",
						"SpatialDimType": "COUNTRY",
						"SpatialDim": "MEX",
						"TimeDimType": "YEAR",
						"TimeDim": 2012,
						"NumericValue": 14.4,
						"Low": 11.5,
						"High": 17.3,
						"Value": "14.4 [11.5-17.3]"
					}]
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
				if len(got.Value) != 1 {
					t.Fatalf("expected 1 observation, got %d", len(got.Value))
				}
				if got.Value[0].SpatialDim != "MEX" {
					t.Fatalf("unexpected SpatialDim: %q", got.Value[0].SpatialDim)
				}
				if got.Value[0].NumericValue != 14.4 {
					t.Fatalf("unexpected NumericValue: %v", got.Value[0].NumericValue)
				}
			},
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
			country: "MEX",
			wantErr: "invalid character 'i' looking for beginning of object key string",
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
