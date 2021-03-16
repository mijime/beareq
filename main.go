package main

import (
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
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/oauth2"
)

func fetchCode(code chan<- string, config *oauth2.Config) error {
	id := uuid.New()
	state := fmt.Sprint(id)
	authCodeURL := config.AuthCodeURL(state)

	log.Println("authCodeURL", authCodeURL)

	if config.RedirectURL == "urn:ietf:wg:oauth:2.0:oob" {
		var c string

		fmt.Fprint(os.Stderr, "Input Code: ")
		fmt.Scan(&c)
		code <- c

		return nil
	}

	redirectURL, err := url.Parse(config.RedirectURL)
	if err != nil {
		return fmt.Errorf("failed to parse redirect url: %w", err)
	}

	if err := webbrowser.Open(authCodeURL); err != nil {
		return fmt.Errorf("failed to open url: %w", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if state != q.Get("state") {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		code <- q.Get("code")
		w.WriteHeader(http.StatusAccepted)
	})

	addr := redirectURL.Hostname() + ":" + redirectURL.Port()

	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func generateToken(config *oauth2.Config) (*oauth2.Token, error) {
	code := make(chan string, 1)

	if err := fetchCode(code, config); err != nil {
		return nil, err
	}

	ctx := context.Background()

	tok, err := config.Exchange(ctx, <-code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange: %w", err)
	}

	return tok, nil
}

func saveToken(tok *oauth2.Token, opts option) error {
	if err := os.MkdirAll(opts.TokenDir, 0700); err != nil {
		return fmt.Errorf("failed to create token dir: %w", err)
	}

	tokenPath := path.Join(opts.TokenDir, opts.Profile+".json")

	tokwp, err := os.Create(tokenPath)
	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	if err := json.NewEncoder(tokwp).Encode(tok); err != nil {
		return fmt.Errorf("failed to encode token: %w", err)
	}

	return nil
}

func fetchToken(config *oauth2.Config, opts option) (*oauth2.Token, error) {
	tokenPath := path.Join(opts.TokenDir, opts.Profile+".json")

	tokfp, err := os.Open(tokenPath)
	if errors.Is(err, os.ErrNotExist) {
		return generateToken(config)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open token: %w", err)
	}

	tok := &oauth2.Token{}
	if err := json.NewDecoder(tokfp).Decode(tok); err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	return tok, nil
}

func fetchConfig(opts option) (*oauth2.Config, error) {
	confp, err := os.Open(opts.ProfilesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open oauth config: %w", err)
	}

	config := make(map[string]oauth2.Config)
	if err := toml.NewDecoder(confp).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode oauth config: %w", err)
	}

	c := config[opts.Profile]

	return &c, nil
}

func doRequest(cli *http.Client, rawurl string, opts option) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}

	req := &http.Request{
		URL:    u,
		Method: opts.Request,
		Header: opts.Header.Values,
	}

	if opts.Data.Exists() {
		req.Body = io.NopCloser(opts.Data)
	}

	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}

	defer resp.Body.Close()

	if opts.Verbose {
		log.Printf("%s %s %s", resp.Request.Method, resp.Request.URL, resp.Request.Proto)
		for k, vs := range resp.Request.Header {
			log.Printf("%s: %s", k, vs)
		}

		log.Printf("%s %s", resp.Proto, resp.Status)
		for k, vs := range resp.Header {
			log.Printf("%s: %s", k, vs)
		}
	}

	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		return fmt.Errorf("failed to copy body: %w", err)
	}

	return nil
}

type option struct {
	Request      string
	Profile      string
	ProfilesPath string
	TokenDir     string
	Data         HTTPRequestBody
	Header       HTTPHeader
	Verbose      bool
}

type HTTPRequestBody struct {
	io.Reader
}

func (b *HTTPRequestBody) Exists() bool {
	return b.Reader != nil
}

func (b *HTTPRequestBody) String() string {
	return ""
}

func (b *HTTPRequestBody) Set(value string) error {
	if value == "-" {
		b.Reader = os.Stdin

		return nil
	}

	if strings.HasPrefix(value, "@") {
		var err error

		b.Reader, err = os.Open(value[1:])
		if err != nil {
			return fmt.Errorf("failed to open request data: %w", err)
		}

		return nil
	}

	b.Reader = strings.NewReader(value)

	return nil
}

type HTTPHeader struct {
	Values map[string][]string
}

func (h *HTTPHeader) String() string {
	return ""
}

func (h *HTTPHeader) Set(values string) error {
	vs := strings.SplitN(values, ":", 2)
	if len(vs[0]) == 0 {
		return nil
	}

	k, v := strings.ToLower(strings.Trim(vs[0], " ")), strings.Trim(vs[1], " ")
	if h.Values[k] == nil {
		h.Values[k] = make([]string, 0)
	}

	h.Values[k] = append(h.Values[k], v)

	return nil
}

func main() {
	profilesPath, _ := homedir.Expand("~/.config/go-oauth-curl/profiles.toml")
	tokenDir, _ := homedir.Expand("~/.config/go-oauth-curl/tokens")

	opts := option{
		Request:      "",
		Header:       HTTPHeader{Values: make(map[string][]string)},
		Data:         HTTPRequestBody{Reader: nil},
		Profile:      "default",
		ProfilesPath: profilesPath,
		TokenDir:     tokenDir,
		Verbose:      false,
	}

	flag.StringVar(&opts.Request, "request", opts.Request, "")
	flag.StringVar(&opts.Profile, "profile", opts.Profile, "")
	flag.StringVar(&opts.ProfilesPath, "config", opts.ProfilesPath, "")
	flag.StringVar(&opts.TokenDir, "tokens", opts.TokenDir, "")
	flag.Var(&opts.Data, "data", "")
	flag.Var(&opts.Header, "header", "")
	flag.BoolVar(&opts.Verbose, "verbose", opts.Verbose, "")
	flag.Parse()

	if len(opts.Request) == 0 {
		if opts.Data.Exists() {
			opts.Request = http.MethodPost
		} else {
			opts.Request = http.MethodGet
		}
	}

	if len(opts.Header.Values["content-type"]) == 0 && opts.Data.Exists() {
		if err := opts.Header.Set("Content-type:application/x-www-form-urlencoded"); err != nil {
			log.Fatal(err)
		}
	}

	urls := flag.Args()

	config, err := fetchConfig(opts)
	if err != nil {
		log.Fatal(err)
	}

	tok, err := fetchToken(config, opts)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	tokSrc := config.TokenSource(ctx, tok)
	cli := oauth2.NewClient(ctx, tokSrc)

	for _, rawurl := range urls {
		if err := doRequest(cli, rawurl, opts); err != nil {
			log.Fatal(err)
		}
	}

	newTok, err := tokSrc.Token()
	if err != nil {
		log.Fatal(err)
	}

	if err := saveToken(newTok, opts); err != nil {
		log.Fatal(err)
	}
}
