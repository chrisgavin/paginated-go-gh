package paginated

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockRoundTripper struct {
	responses []mockResponse
	callCount int
	capturedURLs []string
}

type mockResponse struct {
	statusCode int
	headers    map[string]string
	body       string
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.callCount >= len(m.responses) {
		return nil, fmt.Errorf("unexpected request: %s", req.URL.String())
	}

	m.capturedURLs = append(m.capturedURLs, req.URL.String())
	resp := m.responses[m.callCount]
	m.callCount++

	headers := make(http.Header)
	for k, v := range resp.headers {
		headers.Set(k, v)
	}

	return &http.Response{
		StatusCode: resp.statusCode,
		Header:     headers,
		Body:       io.NopCloser(strings.NewReader(resp.body)),
		Request:    req,
	}, nil
}

func TestNewRoundTripper(t *testing.T) {
	t.Run("with nil transport", func(t *testing.T) {
		rt := NewRoundTripper(nil)
		if rt.transport != http.DefaultTransport {
			t.Error("expected default transport when nil provided")
		}
	})

	t.Run("with custom transport", func(t *testing.T) {
		custom := &mockRoundTripper{}
		rt := NewRoundTripper(custom)
		if rt.transport != custom {
			t.Error("expected custom transport to be used")
		}
	})
}

func TestNonGETRequestsPassthrough(t *testing.T) {
	methods := []string{"POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			mockTransport := &mockRoundTripper{
				responses: []mockResponse{{
					statusCode: 200,
					headers:    map[string]string{"Content-Type": "application/json"},
					body:       `{"message": "success"}`,
				}},
			}

			rt := NewRoundTripper(mockTransport)

			req, err := http.NewRequest(method, "https://api.github.com/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := rt.RoundTrip(req)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != 200 {
				t.Errorf("expected status 200, got %d", resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			resp.Body.Close()

			expected := `{"message": "success"}`
			if string(body) != expected {
				t.Errorf("expected body %s, got %s", expected, string(body))
			}

			if mockTransport.callCount != 1 {
				t.Errorf("expected 1 call to transport, got %d", mockTransport.callCount)
			}
		})
	}
}

func TestSinglePageResponses(t *testing.T) {
	t.Run("GET request without Link header", func(t *testing.T) {
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{{
				statusCode: 200,
				headers:    map[string]string{"Content-Type": "application/json"},
				body:       `[{"id": 1}, {"id": 2}]`,
			}},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/repos/owner/repo/issues", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := rt.RoundTrip(req)
		if err != nil {
			t.Fatal(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		expected := `[{"id": 1}, {"id": 2}]`
		if string(body) != expected {
			t.Errorf("expected body %s, got %s", expected, string(body))
		}

		if mockTransport.callCount != 1 {
			t.Errorf("expected 1 call to transport, got %d", mockTransport.callCount)
		}
	})

	t.Run("GET request adds per_page parameter", func(t *testing.T) {
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{{
				statusCode: 200,
				headers:    map[string]string{"Content-Type": "application/json"},
				body:       `[{"id": 1}]`,
			}},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/repos/owner/repo/issues", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := rt.RoundTrip(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		if len(mockTransport.capturedURLs) == 0 {
			t.Fatal("no URLs captured")
		}
		if !strings.Contains(mockTransport.capturedURLs[0], "per_page=100") {
			t.Errorf("expected per_page=100 to be added to URL, got: %s", mockTransport.capturedURLs[0])
		}
	})

	t.Run("GET request preserves existing per_page parameter", func(t *testing.T) {
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{{
				statusCode: 200,
				headers:    map[string]string{"Content-Type": "application/json"},
				body:       `[{"id": 1}]`,
			}},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/repos/owner/repo/issues?per_page=50", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := rt.RoundTrip(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		if strings.Contains(req.URL.String(), "per_page=100") {
			t.Error("expected existing per_page=50 to be preserved, but found per_page=100")
		}
	})
}

func TestMultiPageArrayResponses(t *testing.T) {
	t.Run("two page array response", func(t *testing.T) {
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{
				{
					statusCode: 200,
					headers: map[string]string{
						"Content-Type": "application/json",
						"Link":         `<https://api.github.com/repos/owner/repo/issues?page=2&per_page=100>; rel="next"`,
					},
					body: `[{"id": 1}, {"id": 2}]`,
				},
				{
					statusCode: 200,
					headers: map[string]string{
						"Content-Type": "application/json",
					},
					body: `[{"id": 3}, {"id": 4}]`,
				},
			},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/repos/owner/repo/issues", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := rt.RoundTrip(req)
		if err != nil {
			t.Fatal(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		var result []map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Fatal(err)
		}

		expected := 4
		if len(result) != expected {
			t.Errorf("expected %d items, got %d", expected, len(result))
		}

		expectedIDs := []float64{1, 2, 3, 4}
		for i, item := range result {
			if item["id"] != expectedIDs[i] {
				t.Errorf("expected item %d to have id %v, got %v", i, expectedIDs[i], item["id"])
			}
		}

		if mockTransport.callCount != 2 {
			t.Errorf("expected 2 calls to transport, got %d", mockTransport.callCount)
		}
	})

	t.Run("three page array response", func(t *testing.T) {
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{
				{
					statusCode: 200,
					headers: map[string]string{
						"Content-Type": "application/json",
						"Link":         `<https://api.github.com/repos/owner/repo/issues?page=2&per_page=100>; rel="next"`,
					},
					body: `[{"id": 1}]`,
				},
				{
					statusCode: 200,
					headers: map[string]string{
						"Content-Type": "application/json",
						"Link":         `<https://api.github.com/repos/owner/repo/issues?page=3&per_page=100>; rel="next"`,
					},
					body: `[{"id": 2}]`,
				},
				{
					statusCode: 200,
					headers: map[string]string{
						"Content-Type": "application/json",
					},
					body: `[{"id": 3}]`,
				},
			},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/repos/owner/repo/issues", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := rt.RoundTrip(req)
		if err != nil {
			t.Fatal(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		var result []map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Fatal(err)
		}

		expected := 3
		if len(result) != expected {
			t.Errorf("expected %d items, got %d", expected, len(result))
		}

		if mockTransport.callCount != 3 {
			t.Errorf("expected 3 calls to transport, got %d", mockTransport.callCount)
		}
	})
}

func TestMultiPageObjectResponses(t *testing.T) {
	t.Run("two page object response with merging", func(t *testing.T) {
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{
				{
					statusCode: 200,
					headers: map[string]string{
						"Content-Type": "application/json",
						"Link":         `<https://api.github.com/search/issues?page=2&per_page=100>; rel="next"`,
					},
					body: `{"total_count": 100, "items": [{"id": 1}, {"id": 2}]}`,
				},
				{
					statusCode: 200,
					headers: map[string]string{
						"Content-Type": "application/json",
					},
					body: `{"total_count": 100, "items": [{"id": 3}, {"id": 4}]}`,
				},
			},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/search/issues?q=is:open", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := rt.RoundTrip(req)
		if err != nil {
			t.Fatal(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Fatal(err)
		}

		if result["total_count"] != float64(100) {
			t.Errorf("expected total_count 100, got %v", result["total_count"])
		}

		items, ok := result["items"].([]interface{})
		if !ok {
			t.Fatal("expected items to be array")
		}

		expected := 4
		if len(items) != expected {
			t.Errorf("expected %d items, got %d", expected, len(items))
		}

		if mockTransport.callCount != 2 {
			t.Errorf("expected 2 calls to transport, got %d", mockTransport.callCount)
		}
	})
}

func TestNonJSONResponses(t *testing.T) {
	t.Run("text response is passed through", func(t *testing.T) {
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{{
				statusCode: 200,
				headers:    map[string]string{"Content-Type": "text/plain"},
				body:       "This is plain text",
			}},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/repos/owner/repo/readme", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := rt.RoundTrip(req)
		if err != nil {
			t.Fatal(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		expected := "This is plain text"
		if string(body) != expected {
			t.Errorf("expected body %s, got %s", expected, string(body))
		}

		if mockTransport.callCount != 1 {
			t.Errorf("expected 1 call to transport, got %d", mockTransport.callCount)
		}
	})

	t.Run("binary response is passed through", func(t *testing.T) {
		binaryData := []byte{0x89, 0x50, 0x4E, 0x47}
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{{
				statusCode: 200,
				headers:    map[string]string{"Content-Type": "image/png"},
				body:       string(binaryData),
			}},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/repos/owner/repo/contents/image.png", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := rt.RoundTrip(req)
		if err != nil {
			t.Fatal(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		if !bytes.Equal(body, binaryData) {
			t.Errorf("expected binary data to be preserved")
		}

		if mockTransport.callCount != 1 {
			t.Errorf("expected 1 call to transport, got %d", mockTransport.callCount)
		}
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("transport error is returned", func(t *testing.T) {
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/repos/owner/repo/issues", nil)
		if err != nil {
			t.Fatal(err)
		}

		_, err = rt.RoundTrip(req)
		if err == nil {
			t.Error("expected error from transport, got nil")
		}
	})

	t.Run("invalid JSON in paginated response", func(t *testing.T) {
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{{
				statusCode: 200,
				headers: map[string]string{
					"Content-Type": "application/json",
					"Link":         `<https://api.github.com/repos/owner/repo/issues?page=2&per_page=100>; rel="next"`,
				},
				body: `invalid json`,
			}},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/repos/owner/repo/issues", nil)
		if err != nil {
			t.Fatal(err)
		}

		_, err = rt.RoundTrip(req)
		if err == nil {
			t.Error("expected JSON parsing error, got nil")
		}
	})

	t.Run("unexpected response type", func(t *testing.T) {
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{{
				statusCode: 200,
				headers: map[string]string{
					"Content-Type": "application/json",
					"Link":         `<https://api.github.com/repos/owner/repo/issues?page=2&per_page=100>; rel="next"`,
				},
				body: `"just a string"`,
			}},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/repos/owner/repo/issues", nil)
		if err != nil {
			t.Fatal(err)
		}

		_, err = rt.RoundTrip(req)
		if err == nil {
			t.Error("expected unexpected response type error, got nil")
		}
		if !strings.Contains(err.Error(), "unexpected response type") {
			t.Errorf("expected 'unexpected response type' error, got: %v", err)
		}
	})

	t.Run("invalid next URL", func(t *testing.T) {
		mockTransport := &mockRoundTripper{
			responses: []mockResponse{{
				statusCode: 200,
				headers: map[string]string{
					"Content-Type": "application/json",
					"Link":         `<invalid-url-with-control-characters%>; rel="next"`,
				},
				body: `[{"id": 1}]`,
			}},
		}

		rt := NewRoundTripper(mockTransport)

		req, err := http.NewRequest("GET", "https://api.github.com/repos/owner/repo/issues", nil)
		if err != nil {
			t.Fatal(err)
		}

		_, err = rt.RoundTrip(req)
		if err == nil {
			t.Error("expected URL parsing error, got nil")
		}
	})
}