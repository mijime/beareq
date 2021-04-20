package attrs

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// HTTPBody is wrapper request body reader.
type HTTPBody struct {
	Reader io.ReadCloser
}

func NewHTTPBody() HTTPBody {
	return HTTPBody{Reader: nil}
}

func (b *HTTPBody) Exists() bool {
	return b.Reader != nil
}

func (b *HTTPBody) String() string {
	return ""
}

func (b *HTTPBody) Set(value string) error {
	if value == "@-" {
		b.Reader = os.Stdin

		return nil
	}

	if strings.HasPrefix(value, "@") {
		rfp, err := os.Open(value[1:])
		if err != nil {
			return fmt.Errorf("failed to open request data: %w", err)
		}

		b.Reader = rfp

		return nil
	}

	b.Reader = io.NopCloser(strings.NewReader(value))

	return nil
}
