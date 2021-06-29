package beareq

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientBuilder interface {
	BuildClient(context.Context) (Client, error)
}

type RequestBuilder interface {
	BuildRequest(context.Context, string) (*http.Request, error)
}

type ResponseHandler interface {
	HandleResponse(context.Context, *http.Response) error
}

func Run(ctx context.Context, cb ClientBuilder, rb RequestBuilder, rh ResponseHandler, urls ...string) error {
	client, err := cb.BuildClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to build client: %w", err)
	}

	if closer, ok := client.(io.Closer); ok {
		defer closer.Close()
	}

	for _, u := range urls {
		err := func(u string) error {
			req, err := rb.BuildRequest(ctx, u)
			if err != nil {
				return fmt.Errorf("failed to build request: %w", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to request: %w", err)
			}

			defer resp.Body.Close()

			if err := rh.HandleResponse(ctx, resp); err != nil {
				return fmt.Errorf("failed to handle response: %w", err)
			}

			return nil
		}(u)
		if err != nil {
			return err
		}
	}

	return nil
}
