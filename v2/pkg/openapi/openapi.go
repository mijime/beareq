package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

	if op.RequestBody != nil {
		for mimeName, schemas := range op.RequestBody.Value.Content {
			op.args[mimeName] = make(map[string]*string)

			for name, prm := range schemas.Schema.Value.Properties {
				var defaultVal string

				for _, v := range []string{
					os.Getenv(strings.ToUpper(strcase.ToSnake(envPrefix + "_BODY_" + name))),
					os.Getenv(strings.ToUpper(strcase.ToSnake(envPrefix + "_" + op.Name() + "_BODY_" + name))),
				} {
					if len(v) > 0 {
						defaultVal = v
					}
				}

				switch prm.Value.Type {
				case "integer":
					op.args[mimeName][name] = fs.String(strcase.ToKebab("body-"+name), defaultVal, prm.Value.Description)
				case "boolean":
					op.args[mimeName][name] = fs.String(strcase.ToKebab("body-"+name), defaultVal, prm.Value.Description)
				case "string":
					op.args[mimeName][name] = fs.String(strcase.ToKebab("body-"+name), defaultVal, prm.Value.Description)
				default:
				}
			}
		}
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
		requestBodyMap := make(map[string]interface{})

		for mimeName, schemas := range op.RequestBody.Value.Content {
			for name, prm := range schemas.Schema.Value.Properties {
				rawv, ok := op.args[mimeName][name]
				if !ok || rawv == nil || len(*rawv) == 0 {
					continue
				}

				switch prm.Value.Type {
				case "integer":
					v, _ := strconv.Atoi(*rawv)
					requestBodyMap[name] = v
				case "boolean":
					v, _ := strconv.ParseBool(*rawv)
					requestBodyMap[name] = v
				case "string":
					requestBodyMap[name] = *rawv
				default:
				}
			}

			if len(requestBodyMap) == 0 {
				continue
			}

			buf, err := json.Marshal(requestBodyMap)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}

			header.Add("Content-type", mimeName)

			requestBody = bytes.NewBuffer(buf)

			break
		}

		if requestBody == nil {
			requestBody = os.Stdin
		}
	}

	if requestBody == nil && len(formData) > 0 {
		header.Add("Content-Type", "application/x-www-form-urlencoded")

		requestBody = strings.NewReader(formData.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, op.Method, baseURI+path, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if len(header) > 0 {
		req.Header = header
	}

	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}

	return req, nil
}

func GenerateOperationFromPath(baseURL string, path string) (map[string]*Operation, error) {
	uri, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse spec path: %w", err)
	}

	loader3 := openapi3.NewLoader()

	doc3, err := loader3.LoadFromURI(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to load openapi v3: %w", err)
	}

	return generateOperation(baseURL, doc3)
}

func GenerateOperationFromData(baseURL string, data []byte) (map[string]*Operation, error) {
	loader3 := openapi3.NewLoader()

	doc3, err := loader3.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to load openapi v3: %w", err)
	}

	return generateOperation(baseURL, doc3)
}

func generateOperation(baseURL string, doc3 *openapi3.T) (map[string]*Operation, error) {
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
