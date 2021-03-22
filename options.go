package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func osGetEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}

	return val
}

type option struct {
	Request      string
	Profile      string
	ProfilesPath string
	TokenDir     string
	Data         httpRequestBody
	Header       httpRequestHeader
	Verbose      bool
	Fail         bool
}

func newOption() (option, error) {
	profilesPath := os.Getenv("HOME") + "/.config/beareq/profiles.toml"
	tokenDir := os.Getenv("HOME") + "/.config/beareq/tokens"
	verbose, err := strconv.ParseBool(osGetEnv("BEAREQ_VERBOSE", "False"))
	if err != nil {
		return option{}, fmt.Errorf("failed to parse verbose flag: %w", err)
	}

	return option{
		Request:      "",
		Header:       httpRequestHeader{Values: make(http.Header)},
		Data:         httpRequestBody{Reader: nil},
		Profile:      osGetEnv("BEAREQ_PROFILE", "default"),
		ProfilesPath: osGetEnv("BEAREQ_PROFILES_PATH", profilesPath),
		TokenDir:     osGetEnv("BEAREQ_TOKENS_DIR", tokenDir),
		Verbose:      verbose,
		Fail:         false,
	}, nil
}

type httpRequestBody struct {
	io.Reader
}

func (b *httpRequestBody) Exists() bool {
	return b.Reader != nil
}

func (b *httpRequestBody) String() string {
	return ""
}

func (b *httpRequestBody) Set(value string) error {
	if value == "-" {
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

	b.Reader = strings.NewReader(value)

	return nil
}

type httpRequestHeader struct {
	Values http.Header
}

func (h *httpRequestHeader) String() string {
	return ""
}

func (h *httpRequestHeader) Set(values string) error {
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

	h.Values.Add(strings.Trim(vs[0], " "), strings.Trim(vs[1], " "))

	return nil
}
