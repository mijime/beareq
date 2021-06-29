package builder

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mijime/beareq/v2/pkg/request/attrs"
)

type RequestBuilder struct {
	Method     attrs.HTTPRequestMethod
	Data       attrs.HTTPBody
	Header     attrs.HTTPHeader
	JSONObject attrs.JSONObject
}

func NewRequestBuilder() RequestBuilder {
	return RequestBuilder{
		Method:     attrs.NewHTTPRequestMethod(),
		Header:     attrs.NewHTTPHeader(),
		Data:       attrs.NewHTTPBody(),
		JSONObject: attrs.NewJSONObject(),
	}
}

func (b RequestBuilder) body() (io.Reader, error) {
	if b.Data.Exists() {
		return b.Data.Reader, nil
	}

	if b.JSONObject.Exists() {
		resp, err := b.JSONObject.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse json object: %w", err)
		}

		return resp, nil
	}

	return nil, nil
}

func (b RequestBuilder) BuildRequest(ctx context.Context, url string) (*http.Request, error) {
	body, err := b.body()
	if err != nil {
		return nil, err
	}

	var method string

	switch {
	case len(b.Method.Method) > 0:
		method = b.Method.Method
	case body != nil:
		method = http.MethodPost
	default:
		method = http.MethodGet
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to new request: %w", err)
	}

	header := b.Header.Values

	if len(header.Get("Content-Type")) == 0 && body != nil {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	req.Header = header

	return req, nil
}
