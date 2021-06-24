package openapi

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
)

type Operation struct {
	BaseURL string
	Path    string
	Method  string

	*openapi3.Operation
	Variables map[string]map[string]*string
}

func NewOperation(url, path, method string, op *openapi3.Operation) *Operation {
	return &Operation{
		BaseURL:   url,
		Path:      path,
		Method:    method,
		Operation: op,

		Variables: make(map[string]map[string]*string),
	}
}

func (op *Operation) Name() string {
	return op.OperationID
}

func (op *Operation) Parse(envPrefix string, args []string) error {
	fs := flag.NewFlagSet(op.Name(), flag.ExitOnError)

	for _, prm := range op.Parameters {
		argName := strcase.ToKebab(prm.Value.Name)

		if _, ok := op.Variables[prm.Value.In]; !ok {
			op.Variables[prm.Value.In] = make(map[string]*string)
		}

		var defaultVal string

		for _, v := range []string{
			os.Getenv(strings.ToUpper(strcase.ToSnake(envPrefix + "_" + argName))),
			os.Getenv(strings.ToUpper(strcase.ToSnake(envPrefix + "_" + op.Name() + "_" + argName))),
		} {
			if len(v) > 0 {
				defaultVal = v
			}
		}

		op.Variables[prm.Value.In][prm.Value.Name] = fs.String(argName, defaultVal, prm.Value.Description)
	}

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse args: %w", err)
	}

	return nil
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
