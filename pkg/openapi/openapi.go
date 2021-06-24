package openapi

import (
	"context"
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
	*openapi3.Operation

	BaseURL string
	Path    string
	Method  string

	args map[string]map[string]*string
}

func NewOperation(url, path, method string, op *openapi3.Operation) *Operation {
	return &Operation{
		Operation: op,
		BaseURL:   url,
		Path:      path,
		Method:    method,

		args: make(map[string]map[string]*string),
	}
}

func (op *Operation) Name() string {
	return op.OperationID
}

func (op *Operation) Parse(envPrefix string, args []string) error {
	fs := flag.NewFlagSet(op.Name(), flag.ExitOnError)

	for _, prm := range op.Parameters {
		if _, ok := op.args[prm.Value.In]; !ok {
			op.args[prm.Value.In] = make(map[string]*string)
		}

		var defaultVal string

		for _, v := range []string{
			os.Getenv(strings.ToUpper(strcase.ToSnake(envPrefix + "_" + prm.Value.Name))),
			os.Getenv(strings.ToUpper(strcase.ToSnake(envPrefix + "_" + op.Name() + "_" + prm.Value.Name))),
		} {
			if len(v) > 0 {
				defaultVal = v
			}
		}

		argName := strcase.ToKebab(prm.Value.Name)
		op.args[prm.Value.In][prm.Value.Name] = fs.String(argName, defaultVal, prm.Value.Description)
	}

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse args: %w", err)
	}

	return nil
}

func (op *Operation) BuildRequest(ctx context.Context, baseURI string) (*http.Request, error) {
	path := op.Path
	query := make(url.Values)
	header := make(http.Header)
	formData := make(url.Values)

	for _, prm := range op.Parameters {
		v, ok := op.args[prm.Value.In][prm.Value.Name]
		if !ok || v == nil || len(*v) == 0 {
			continue
		}

		switch prm.Value.In {
		case "path":
			path = strings.ReplaceAll(path, "{"+prm.Value.Name+"}", *v)
		case "query":
			query.Add(prm.Value.Name, *v)
		case "header":
			header.Add(prm.Value.Name, *v)
		case "formData":
			formData.Add(prm.Value.Name, *v)
		default:
		}
	}

	var requestBody io.Reader
	if op.RequestBody != nil {
		requestBody = os.Stdin
	}

	req, err := http.NewRequestWithContext(ctx, op.Method, baseURI+path, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if len(header) > 0 {
		req.Header = header
	}

	if len(formData) > 0 {
		req.PostForm = formData
	}

	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}

	return req, nil
}

func GenerateOperation(baseURL string, specPath string) (map[string]*Operation, error) {
	uri, err := url.Parse(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse spec path: %w", err)
	}

	loader := openapi3.NewLoader()

	doc3, err := loader.LoadFromURI(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to load openapi v3: %w", err)
	}

	for _, v := range doc3.Servers {
		if len(baseURL) == 0 {
			baseURL = v.URL
		}
	}

	cmds := make(map[string]*Operation)

	for path, v := range doc3.Paths {
		for method, op := range v.Operations() {
			cmd := NewOperation(baseURL, path, method, op)
			cmds[cmd.Name()] = cmd
		}
	}

	return cmds, nil
}
