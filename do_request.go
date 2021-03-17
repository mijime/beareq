package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func dumpHeader(h http.Header) {
	for k, vs := range h {
		for _, v := range vs {
			log.Printf("%s: %s", k, v)
		}
	}
}

func dumpResponse(resp *http.Response) {
	log.Printf("%s %s %s", resp.Request.Method, resp.Request.URL, resp.Request.Proto)
	dumpHeader(resp.Request.Header)

	log.Printf("%s %s", resp.Proto, resp.Status)
	dumpHeader(resp.Header)
}

func selectMethod(opts option) string {
	if len(opts.Request) > 0 {
		return opts.Request
	}

	if opts.Data.Exists() {
		return http.MethodPost
	}

	return http.MethodGet
}

func createRequest(u *url.URL, opts option) *http.Request {
	header := opts.Header.Values
	if len(header.Get("Content-Type")) == 0 && opts.Data.Exists() {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	req := &http.Request{
		URL:    u,
		Method: selectMethod(opts),
		Header: header,
	}

	if opts.Data.Exists() {
		req.Body = io.NopCloser(opts.Data)
	}

	return req
}

func matchCommand(cc commandsConfigs, u *url.URL) *commandsConfig {
	if cc[u.Scheme] == nil {
		return nil
	}

	if cc[u.Scheme][u.Hostname()] == nil {
		return nil
	}

	return cc[u.Scheme][u.Hostname()]
}

func doRequest(cli *http.Client, rawurl string, opts option) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}

	cmds, err := fetchCommands(opts)
	if err != nil {
		return err
	}

	var req *http.Request

	if cmd := matchCommand(cmds, u); cmd != nil {
		req, err = cmd.Create(u, opts)
		if err != nil {
			return err
		}
	} else {
		req = createRequest(u, opts)
	}

	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}

	defer resp.Body.Close()

	if opts.Verbose {
		dumpResponse(resp)
	}

	if opts.Fail && resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("failed to request: %w", errors.New(resp.Status))
	}

	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		return fmt.Errorf("failed to copy body: %w", err)
	}

	return nil
}
