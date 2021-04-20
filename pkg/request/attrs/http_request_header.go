package attrs

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// HTTPHeader is wrapper http.Header for flags.
type HTTPHeader struct {
	Values http.Header
}

func NewHTTPHeader() HTTPHeader {
	return HTTPHeader{Values: make(http.Header)}
}

func (h *HTTPHeader) String() string {
	return ""
}

// Set is parse to argument for http header.
func (h *HTTPHeader) Set(values string) error {
	if strings.HasPrefix(values, "@") {
		rfp, err := os.Open(values[1:])
		if err != nil {
			return fmt.Errorf("failed to open request header: %w", err)
		}

		sc := bufio.NewScanner(rfp)
		for sc.Scan() {
			if err := h.Set(sc.Text()); err != nil {
				return err
			}
		}

		return nil
	}

	vs := strings.SplitN(values, ":", 2)
	if len(vs[0]) == 0 {
		return nil
	}

	h.Values.Add(strings.TrimSpace(vs[0]), strings.TrimSpace(vs[1]))

	return nil
}
