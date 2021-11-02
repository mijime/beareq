package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
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
	PathItem *openapi3.PathItem

	BaseURL string
	Path    string
	Method  string

	args map[string]map[string]flag.Value
}

func NewOperation(url, path, method string, op *openapi3.Operation, pi *openapi3.PathItem) *Operation {
	return &Operation{
		Operation: op,
		PathItem:  pi,
		BaseURL:   url,
		Path:      path,
		Method:    method,

		args: make(map[string]map[string]flag.Value),
	}
}

func (op *Operation) Name() string {
	return op.OperationID
}

func (op *Operation) FlagSet(envPrefix string) *flag.FlagSet {
	fs := flag.NewFlagSet(op.Name(), flag.ExitOnError)

	for _, prm := range append(op.PathItem.Parameters, op.Parameters...) {
		if _, ok := op.args[prm.Value.In]; !ok {
			op.args[prm.Value.In] = make(map[string]flag.Value)
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
		op.args[prm.Value.In][prm.Value.Name] = &argString{}

		if err := op.args[prm.Value.In][prm.Value.Name].Set(defaultVal); err != nil {
			log.Printf("[WARN] failed to set default value: %s = %v", argName, err)
		}

		fs.Var(op.args[prm.Value.In][prm.Value.Name], argName, prm.Value.Description)
	}

	if op.RequestBody != nil {
		for mimeName, schemas := range op.RequestBody.Value.Content {
			op.args[mimeName] = make(map[string]flag.Value)
			ao := &argObject{V: make(map[string]flag.Value)}
			op.args[mimeName]["body"] = ao
			op.buildBodyArgs(fs, ao, envPrefix, "body-", schemas.Schema.Value)
		}
	}

	return fs
}

func (op *Operation) buildBodyArgs(fs *flag.FlagSet, ao *argObject, envPrefix, namePrefix string, schema *openapi3.Schema) {
	for name, prm := range schema.Properties {
		var defaultVal string

		for _, v := range []string{
			os.Getenv(strings.ToUpper(strcase.ToSnake(envPrefix + "_" + namePrefix + name))),
			os.Getenv(strings.ToUpper(strcase.ToSnake(envPrefix + "_" + op.Name() + "_" + namePrefix + name))),
		} {
			if len(v) > 0 {
				defaultVal = v
			}
		}

		argName := strcase.ToKebab(namePrefix + name)

		switch prm.Value.Type {
		case "integer":
			ao.V[namePrefix+name] = &argInteger{}

			if err := ao.V[namePrefix+name].Set(defaultVal); err != nil {
				log.Printf("[WARN] failed to set default value: %s = %v", argName, err)
			}

			fs.Var(ao.V[namePrefix+name], argName, prm.Value.Description)
		case "number":
			ao.V[namePrefix+name] = &argFloat{}

			if err := ao.V[namePrefix+name].Set(defaultVal); err != nil {
				log.Printf("[WARN] failed to set default value: %s = %v", argName, err)
			}

			fs.Var(ao.V[namePrefix+name], argName, prm.Value.Description)
		case "boolean":
			ao.V[namePrefix+name] = &argBoolean{}

			if err := ao.V[namePrefix+name].Set(defaultVal); err != nil {
				log.Printf("[WARN] failed to set default value: %s = %v", argName, err)
			}

			fs.Var(ao.V[namePrefix+name], argName, prm.Value.Description)
		case "string":
			ao.V[namePrefix+name] = &argString{}

			if err := ao.V[namePrefix+name].Set(defaultVal); err != nil {
				log.Printf("[WARN] failed to set default value: %s = %v", argName, err)
			}

			fs.Var(ao.V[namePrefix+name], argName, prm.Value.Description)
		case "array":
			ao.V[namePrefix+name] = &argArray{V: make([]string, 0)}

			if err := ao.V[namePrefix+name].Set(defaultVal); err != nil {
				log.Printf("[WARN] failed to set default value: %s = %v", argName, err)
			}

			fs.Var(ao.V[namePrefix+name], argName, prm.Value.Description)
		case "object":
			cao := &argObject{V: make(map[string]flag.Value)}
			ao.V[namePrefix+name] = cao
			op.buildBodyArgs(fs, cao, envPrefix, namePrefix+name+"-", prm.Value)
		default:
			log.Printf("[WARN] unspported type: %s = %v", strcase.ToKebab(namePrefix+name), prm.Value.Type)
		}
	}
}

func buildReuqestBodyFromSchema(ao *argObject, namePrefix string, schema *openapi3.Schema) map[string]interface{} {
	body := make(map[string]interface{})

	for name, prm := range schema.Properties {
		v, ok := ao.V[namePrefix+name]
		if !ok {
			continue
		}

		switch tv := v.(type) {
		case *argFloat:
			if !tv.IsSet {
				continue
			}

			body[name] = tv.V
		case *argInteger:
			if !tv.IsSet {
				continue
			}

			body[name] = tv.V
		case *argBoolean:
			if !tv.IsSet {
				continue
			}

			body[name] = tv.V
		case *argString:
			if !tv.IsSet {
				continue
			}

			body[name] = tv.V
		case *argArray:
			if !tv.IsSet {
				continue
			}

			body[name] = tv.V
		case *argObject:
			b := buildReuqestBodyFromSchema(tv, namePrefix+name+"-", prm.Value)

			if len(b) > 0 {
				body[name] = b
			}
		default:
			log.Printf("[WARN] unspported type: %s = %v", strcase.ToKebab(namePrefix+name), prm.Value.Type)
		}
	}

	return body
}

func (op *Operation) BuildRequest(ctx context.Context, baseURI string) (*http.Request, error) {
	path := op.Path
	query := make(url.Values)
	header := make(http.Header)
	formData := make(url.Values)

	for _, prm := range append(op.PathItem.Parameters, op.Parameters...) {
		v, ok := op.args[prm.Value.In][prm.Value.Name]
		if !ok {
			continue
		}

		switch tv := v.(type) {
		case *argString:
			if !tv.IsSet {
				continue
			}

			switch prm.Value.In {
			case "path":
				path = strings.ReplaceAll(path, "{"+prm.Value.Name+"}", tv.V)
			case "query":
				query.Add(prm.Value.Name, tv.V)
			case "header":
				header.Add(prm.Value.Name, tv.V)
			case "formData":
				formData.Add(prm.Value.Name, tv.V)
			default:
				log.Printf("[WARN] unspported type: %s = %v", prm.Value.Name, prm.Value.In)
			}
		default:
			log.Printf("[WARN] unspported type: %s = %v", prm.Value.Name, tv)
		}
	}

	var requestBody io.Reader

	if op.RequestBody != nil {
		var requestBodyMap map[string]interface{}

		for mimeName, schemas := range op.RequestBody.Value.Content {
			bargs, ok := op.args[mimeName]["body"]
			if !ok {
				continue
			}

			bvo, ok := bargs.(*argObject)
			if !ok {
				continue
			}

			requestBodyMap = buildReuqestBodyFromSchema(bvo, "body-", schemas.Schema.Value)
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
	for _, srv := range doc3.Servers {
		if len(baseURL) == 0 {
			baseURL = srv.URL
		}
	}

	cmds := make(map[string]*Operation)

	for path, pi := range doc3.Paths {
		for method, op := range pi.Operations() {
			cmd := NewOperation(baseURL, path, method, op, pi)
			cmds[cmd.Name()] = cmd
		}
	}

	return cmds, nil
}

type argInteger struct {
	V     int
	IsSet bool
}

func (a *argInteger) String() string {
	return strconv.Itoa(a.V)
}

func (a *argInteger) Set(v string) error {
	tv, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("failed to set integer: %w", err)
	}

	a.V = tv

	return nil
}

type argBoolean struct {
	V     bool
	IsSet bool
}

func (a *argBoolean) String() string {
	return strconv.FormatBool(a.V)
}

func (a *argBoolean) Set(v string) error {
	if len(v) == 0 {
		return nil
	}

	tv, err := strconv.ParseBool(v)
	if err != nil {
		return fmt.Errorf("failed to set argBoolean: %w", err)
	}

	a.IsSet = true
	a.V = tv

	return nil
}

type argFloat struct {
	V     float64
	IsSet bool
}

func (a *argFloat) String() string {
	return ""
}

func (a *argFloat) Set(v string) error {
	if len(v) == 0 {
		return nil
	}

	tv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("failed to set argFloat: %w", err)
	}

	a.IsSet = true
	a.V = tv

	return nil
}

type argString struct {
	V     string
	IsSet bool
}

func (a *argString) String() string {
	return a.V
}

func (a *argString) Set(v string) error {
	if len(v) == 0 {
		return nil
	}

	a.IsSet = true
	a.V = v

	return nil
}

type argArray struct {
	V     []string
	IsSet bool
}

func (a *argArray) String() string {
	return strings.Join(a.V, ", ")
}

func (a *argArray) Set(v string) error {
	if len(v) == 0 {
		return nil
	}

	a.IsSet = true
	a.V = append(a.V, v)

	return nil
}

type argObject struct {
	V map[string]flag.Value
}

func (a *argObject) String() string {
	return ""
}

func (a *argObject) Set(v string) error {
	return errors.New("failed to set")
}
