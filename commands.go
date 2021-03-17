package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/pelletier/go-toml"
)

type commandsConfig struct {
	URL    string
	Method string
	Header map[string][]string
	Data   string
}

func (cc *commandsConfig) Create(u *url.URL, opts option) (*http.Request, error) {
	urlTmpl := template.Must(template.New("").Parse(cc.URL))
	urlBuf := bytes.NewBuffer([]byte{})

	if err := urlTmpl.Execute(urlBuf, u.Query()); err != nil {
		return nil, fmt.Errorf("failed to build command url: %w", err)
	}

	reqURL, err := url.Parse(urlBuf.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse command url: %w", err)
	}

	req := &http.Request{
		URL:    reqURL,
		Method: cc.Method,
		Header: opts.Header.Values,
	}

	for h, vs := range cc.Header {
		if len(req.Header.Get(h)) > 0 {
			continue
		}

		for _, v := range vs {
			req.Header.Add(h, v)
		}
	}

	if opts.Data.Exists() {
		req.Body = io.NopCloser(opts.Data.Reader)

		return req, nil
	}

	if len(cc.Data) == 0 {
		return req, nil
	}

	dataTmpl, err := template.New("").Parse(cc.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command body: %w", err)
	}

	dataBuf := bytes.NewBuffer([]byte{})
	if err := dataTmpl.Execute(dataBuf, u.Query()); err != nil {
		return nil, fmt.Errorf("failed to build command body: %w", err)
	}

	req.Body = io.NopCloser(dataBuf)

	return req, nil
}

type commandsConfigs = map[string]map[string]*commandsConfig

func fetchCommands(opts option) (commandsConfigs, error) {
	config := make(commandsConfigs)

	confp, err := os.Open(opts.CommandsPath)
	if errors.Is(err, os.ErrNotExist) {
		return config, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open commands config: %w", err)
	}

	if err := toml.NewDecoder(confp).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode commands config: %w", err)
	}

	return config, nil
}
