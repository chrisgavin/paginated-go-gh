package paginated

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"

	"dario.cat/mergo"
	"github.com/tomnomnom/linkheader"
)

type PaginatingRoundTripper struct {
	transport http.RoundTripper
}

func NewRoundTripper(transport http.RoundTripper) *PaginatingRoundTripper {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &PaginatingRoundTripper{transport: transport}
}

func (p *PaginatingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != http.MethodGet {
		return p.transport.RoundTrip(req)
	}

	originalURL := req.URL.String()
	paginatedURL, err := setPerPage(originalURL)
	if err != nil {
		return nil, err
	}

	if paginatedURL != originalURL {
		newReq := req.Clone(req.Context())
		newReq.URL, err = url.Parse(paginatedURL)
		if err != nil {
			return nil, err
		}
		req = newReq
	}

	var httpResponse *http.Response
	var combinedMap map[string]interface{}
	var combinedSlice []interface{}
	currentReq := req
	isFirstRequest := true

	for {
		httpResponse, err = p.transport.RoundTrip(currentReq)
		if err != nil {
			return nil, err
		}

		contentType, _, err := mime.ParseMediaType(httpResponse.Header.Get("Content-Type"))
		if err != nil {
			return httpResponse, nil
		}
		if contentType != "application/json" {
			return httpResponse, nil
		}

		if httpResponse.Header.Get("Link") == "" && isFirstRequest {
			return httpResponse, nil
		}

		var response interface{}
		err = json.NewDecoder(httpResponse.Body).Decode(&response)
		httpResponse.Body.Close()
		if err != nil {
			return nil, err
		}

		switch response := response.(type) {
		case []interface{}:
			combinedSlice = append(combinedSlice, response...)
		case map[string]interface{}:
			err = mergo.Merge(&combinedMap, response, mergo.WithAppendSlice)
		default:
			err = errors.New("unexpected response type")
		}

		if err != nil {
			return nil, err
		}

		isFirstRequest = false

		links := linkheader.Parse(httpResponse.Header.Get("Link"))
		next := links.FilterByRel("next")
		if len(next) == 0 {
			break
		}

		nextURL, err := url.Parse(next[0].URL)
		if err != nil {
			return nil, err
		}

		currentReq = req.Clone(req.Context())
		currentReq.URL = nextURL
	}

	var marshaled []byte
	if len(combinedSlice) > 0 {
		marshaled, err = json.Marshal(combinedSlice)
	} else if combinedMap != nil {
		marshaled, err = json.Marshal(combinedMap)
	} else {
		marshaled, err = json.Marshal(combinedSlice)
	}
	if err != nil {
		return nil, err
	}

	httpResponse.Body = io.NopCloser(bytes.NewReader(marshaled))
	httpResponse.ContentLength = int64(len(marshaled))
	httpResponse.Header.Set("Content-Length", fmt.Sprintf("%d", len(marshaled)))

	return httpResponse, nil
}

func setPerPage(input string) (string, error) {
	parsed, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	parsedQuery := parsed.Query()
	if parsedQuery.Get("per_page") == "" {
		parsedQuery.Set("per_page", "100")
		parsed.RawQuery = parsedQuery.Encode()
		return parsed.String(), nil
	}
	return input, nil
}
