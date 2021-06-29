package attrs

import (
	"fmt"
	"net/http"
	"strings"
)

type HTTPRequestMethod struct {
	Method string
}

func NewHTTPRequestMethod() HTTPRequestMethod {
	return HTTPRequestMethod{Method: ""}
}

func (a *HTTPRequestMethod) Exists() bool {
	return len(a.Method) > 0
}

func (a *HTTPRequestMethod) String() string {
	if a == nil {
		return ""
	}

	return a.Method
}

func (a *HTTPRequestMethod) Set(value string) error {
	switch strings.ToUpper(value) {
	case http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		a.Method = value

		return nil
	default:
		return fmt.Errorf("not support http method: %s", value)
	}
}
