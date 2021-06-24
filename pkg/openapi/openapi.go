package openapi

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
)

type Operation struct {
	*flag.FlagSet
	BaseURL   string
	Path      string
	Method    string
	Variables map[string]map[string]*string
}

func NewOperation(url, path, method string, op *openapi3.Operation) *Operation {
	fs := flag.NewFlagSet(op.OperationID, flag.ExitOnError)

	variables := make(map[string]map[string]*string)

	for _, prm := range op.Parameters {
		argName := strcase.ToKebab(prm.Value.Name)

		if _, ok := variables[prm.Value.In]; !ok {
			variables[prm.Value.In] = make(map[string]*string)
		}

		variables[prm.Value.In][prm.Value.Name] = fs.String(argName, "", prm.Value.Description)
	}

	return &Operation{
		BaseURL:   url,
		Path:      path,
		Method:    method,
		FlagSet:   fs,
		Variables: variables,
	}
}

func (cmd *Operation) BuildRequest(ctx context.Context, baseURI string) (*http.Request, error) {
	p := cmd.Path

	for k, v := range cmd.Variables["path"] {
		if v == nil || len(*v) == 0 {
			continue
		}

		p = strings.ReplaceAll(p, "{"+k+"}", *v)
	}

	req, err := http.NewRequestWithContext(ctx, cmd.Method, baseURI+p, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if m, ok := cmd.Variables["header"]; !ok {
		for k, v := range m {
			if v == nil || len(*v) == 0 {
				continue
			}

			req.Header.Add(k, *v)
		}
	}

	body := make(url.Values)

	if m, ok := cmd.Variables["formData"]; !ok {
		for k, v := range m {
			if v == nil || len(*v) == 0 {
				continue
			}

			body.Add(k, *v)
		}
	}

	if len(body) > 0 {
		req.PostForm = body
	}

	query := make(url.Values)

	if m, ok := cmd.Variables["query"]; !ok {
		for k, v := range m {
			if v == nil || len(*v) == 0 {
				continue
			}

			query.Add(k, *v)
		}
	}

	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}

	return req, nil
}

func GenerateOperation(baseURL string, openapiJSON io.Reader) (map[string]*Operation, error) {
	cmds := make(map[string]*Operation)

	var doc3 openapi3.T
	if err := json.NewDecoder(openapiJSON).Decode(&doc3); err != nil {
		return nil, fmt.Errorf("failed to decode openapi v3: %w", err)
	}

	for _, v := range doc3.Servers {
		if len(baseURL) == 0 {
			baseURL = v.URL
		}
	}

	for path, v := range doc3.Paths {
		for method, op := range v.Operations() {
			cmd := NewOperation(baseURL, path, method, op)
			cmds[cmd.Name()] = cmd
		}
	}

	return cmds, nil
}
