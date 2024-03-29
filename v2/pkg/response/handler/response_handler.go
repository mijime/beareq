package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mijime/beareq/v2/pkg/response/attrs"
	"github.com/mijime/beareq/v2/pkg/response/handler/jq"
	"github.com/mijime/beareq/v2/pkg/response/handler/raw"
)

type ResponseHandler struct {
	JSONQuery attrs.JSONQuery
	Verbose   bool
	Fail      bool
}

func NewResponseHandler() ResponseHandler {
	return ResponseHandler{
		JSONQuery: attrs.NewJSONQuery(),
		Verbose:   false,
		Fail:      false,
	}
}

func (h ResponseHandler) HandleResponse(ctx context.Context, resp *http.Response) error {
	if h.Verbose {
		dumpResponse(resp)
	}

	if h.Fail && resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("failed to request: %s", resp.Status)
	}

	if h.JSONQuery.Exists() && resp.StatusCode < http.StatusBadRequest {
		rh := &jq.ResponseHandler{
			Writer: os.Stdout,
			Query:  h.JSONQuery.Query,
		}

		if err := rh.HandleResponse(ctx, resp); err != nil {
			return fmt.Errorf("failed to handle response: %w", err)
		}

		return nil
	}

	rh := &raw.ResponseHandler{
		Writer: os.Stdout,
	}

	if err := rh.HandleResponse(ctx, resp); err != nil {
		return fmt.Errorf("failed to handle response: %w", err)
	}

	return nil
}

func dumpHeader(h http.Header) {
	for k, vs := range h {
		for _, v := range vs {
			log.Printf("%s: %s", k, v)
		}
	}
}

func dumpResponse(resp *http.Response) {
	log.Printf("%s %s %s", resp.Request.Method, resp.Request.URL, resp.Request.Proto)
	dumpHeader(resp.Request.Header)

	log.Printf("%s %s", resp.Proto, resp.Status)
	dumpHeader(resp.Header)
}
