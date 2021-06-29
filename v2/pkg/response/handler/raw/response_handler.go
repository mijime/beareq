package raw

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type ResponseHandler struct {
	io.Writer
}

func (h *ResponseHandler) HandleResponse(ctx context.Context, resp *http.Response) error {
	if _, err := io.Copy(h.Writer, resp.Body); err != nil {
		return fmt.Errorf("failed to copy response body: %w", err)
	}

	return nil
}
