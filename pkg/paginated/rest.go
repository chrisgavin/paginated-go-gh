package paginated

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"

	"github.com/chrisgavin/paginated-go-gh/internal/merge"
	"github.com/cli/go-gh/pkg/api"
	"github.com/tomnomnom/linkheader"
)

type PaginatingRESTClient struct {
	client api.RESTClient
}

func WrapClient(client api.RESTClient) *PaginatingRESTClient {
	return &PaginatingRESTClient{client: client}
}

func (paginating *PaginatingRESTClient) Do(method string, path string, body io.Reader, response interface{}) error {
	return paginating.DoWithContext(context.Background(), method, path, body, response)
}

func (paginating *PaginatingRESTClient) DoWithContext(ctx context.Context, method string, path string, body io.Reader, response interface{}) error {
	httpResponse, err := paginating.RequestWithContext(ctx, method, path, body)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode == http.StatusNoContent {
		return nil
	}

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return err
	}

	return nil
}

func (paginating *PaginatingRESTClient) Delete(path string, response interface{}) error {
	return paginating.Do(http.MethodDelete, path, nil, response)
}

func (paginating *PaginatingRESTClient) Get(path string, response interface{}) error {
	return paginating.Do(http.MethodGet, path, nil, response)
}

func (paginating *PaginatingRESTClient) Patch(path string, body io.Reader, response interface{}) error {
	return paginating.Do(http.MethodPatch, path, body, response)
}

func (paginating *PaginatingRESTClient) Post(path string, body io.Reader, response interface{}) error {
	return paginating.Do(http.MethodPost, path, body, response)
}

func (paginating *PaginatingRESTClient) Put(path string, body io.Reader, response interface{}) error {
	return paginating.Do(http.MethodPut, path, body, response)
}

func (paginating *PaginatingRESTClient) Request(method string, path string, body io.Reader) (*http.Response, error) {
	return paginating.RequestWithContext(context.Background(), method, path, body)
}

func (paginating *PaginatingRESTClient) RequestWithContext(ctx context.Context, method string, path string, body io.Reader) (*http.Response, error) {
	if method != http.MethodGet {
		return paginating.client.RequestWithContext(ctx, method, path, body)
	}

	path, err := setPerPage(path)
	if err != nil {
		return nil, err
	}

	var httpResponse *http.Response
	var combinedResponse interface{}
	for path != "" {
		httpResponse, err = paginating.client.RequestWithContext(ctx, method, path, body)
		if err != nil {
			return nil, err
		}

		contentType, _, err := mime.ParseMediaType(httpResponse.Header.Get("Content-Type"))
		if err != nil {
			return nil, err
		}
		if contentType != "application/json" {
			return httpResponse, nil
		}

		if httpResponse.Header.Get("Link") == "" {
			return httpResponse, nil
		}

		var response interface{}
		err = json.NewDecoder(httpResponse.Body).Decode(&response)
		if err != nil {
			return nil, err
		}
		merge.MergeResponses(&combinedResponse, &response)

		links := linkheader.Parse(httpResponse.Header.Get("Link"))
		next := links.FilterByRel("next")
		if len(next) == 0 {
			path = ""
		} else {
			path = next[0].URL
		}
	}

	marshaled, err := json.Marshal(combinedResponse)
	if err != nil {
		return nil, err
	}
	httpResponse.Body = io.NopCloser(io.Reader(bytes.NewReader(marshaled)))

	return httpResponse, nil
}

func setPerPage(input string) (string, error) {
	parsed, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	parsedQuery, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		return "", err
	}
	parsedQuery.Set("per_page", "100")
	parsed.RawQuery = parsedQuery.Encode()
	return parsed.String(), nil
}
