package jq

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/itchyny/gojq"
)

type ResponseHandler struct {
	io.Writer
	*gojq.Query
}

func (h *ResponseHandler) HandleResponse(ctx context.Context, resp *http.Response) error {
	var v interface{}

	dec := json.NewDecoder(resp.Body)

	if err := dec.Decode(&v); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}

	enc := json.NewEncoder(h.Writer)
	enc.SetIndent("", "  ")

	iter := h.Query.RunWithContext(ctx, v)

	for {
		nv, ok := iter.Next()

		if !ok {
			return nil
		}

		if ns, ok := nv.(string); ok {
			if _, err := fmt.Fprintln(h.Writer, ns); err != nil {
				return fmt.Errorf("failed to copy string: %w", err)
			}

			continue
		}

		if err := enc.Encode(nv); err != nil {
			return fmt.Errorf("failed to encode json: %w", err)
		}
	}
}
