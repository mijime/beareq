package builder

import (
	"context"
	"crypto/tls"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/mijime/beareq/v2/pkg/beareq"
	"github.com/pelletier/go-toml"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/oauth2"
)

type ClientBuilder struct {
	Profile      string
	ProfilesPath string
	TokenDir     string
	RefreshToken bool
}

func osGetEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}

	return val
}

func NewClientBuilder() *ClientBuilder {
	profilesPath := os.Getenv("HOME") + "/.config/beareq/profiles.toml"
	tokenDir := os.Getenv("HOME") + "/.config/beareq/tokens"

	return &ClientBuilder{
		Profile:      osGetEnv("BEAREQ_PROFILE", "default"),
		ProfilesPath: osGetEnv("BEAREQ_PROFILES_PATH", profilesPath),
		TokenDir:     osGetEnv("BEAREQ_TOKENS_DIR", tokenDir),
	}
}

type profileConfig struct {
	OAuth  *oauth2.Config
	Header map[string][]string
}

func (b *ClientBuilder) fetchConfigByProfile() (profileConfig, error) {
	confp, err := os.Open(b.ProfilesPath)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("not found oauth config:", err)

		return profileConfig{}, nil
	}

	if err != nil {
		return profileConfig{}, fmt.Errorf("failed to open oauth config: %w", err)
	}

	defer confp.Close()

	config := make(map[string]profileConfig)
	if err := toml.NewDecoder(confp).Decode(&config); err != nil {
		return profileConfig{}, fmt.Errorf("failed to decode oauth config: %w", err)
	}

	c, ok := config[b.Profile]
	if !ok {
		return profileConfig{}, fmt.Errorf("not found profile in profiles.toml: %s", b.Profile)
	}

	return c, nil
}

type headerOverride map[string][]string

func (r headerOverride) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, vs := range r {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("failed to default round trip: %w", err)
	}

	return resp, nil
}

func (b *ClientBuilder) BuildClient(ctx context.Context) (beareq.Client, error) {
	config, err := b.fetchConfigByProfile()
	if err != nil {
		return nil, err
	}

	if config.OAuth != nil {
		tok, err := b.fetchToken(ctx, config.OAuth)
		if err != nil {
			return nil, err
		}

		tokSrc := config.OAuth.TokenSource(ctx, tok)
		cli := oauth2.NewClient(ctx, tokSrc)

		return &profileClient{
			client:  cli,
			builder: b,
			source:  tokSrc,
		}, nil
	}

	if config.Header != nil {
		return &http.Client{Transport: headerOverride(config.Header)}, nil
	}

	return http.DefaultClient, nil
}

func (b *ClientBuilder) saveToken(tokSrc oauth2.TokenSource) error {
	if err := os.MkdirAll(b.TokenDir, 0o700); err != nil {
		return fmt.Errorf("failed to create token dir: %w", err)
	}

	tokenPath := path.Join(b.TokenDir, b.Profile+".json")

	tokwp, err := os.Create(tokenPath)
	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	newTok, err := tokSrc.Token()
	if err != nil {
		return fmt.Errorf("failed to create new token: %w", err)
	}

	if err := json.NewEncoder(tokwp).Encode(newTok); err != nil {
		return fmt.Errorf("failed to encode token: %w", err)
	}

	return nil
}

func (b *ClientBuilder) fetchToken(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	tokenPath := path.Join(b.TokenDir, b.Profile+".json")

	tokfp, err := os.Open(tokenPath)
	if b.RefreshToken || errors.Is(err, os.ErrNotExist) {
		return generateToken(ctx, config)
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

type profileClient struct {
	client *http.Client

	builder *ClientBuilder
	source  oauth2.TokenSource
}

func (c *profileClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	return resp, nil
}

func (c *profileClient) Close() error {
	return c.builder.saveToken(c.source)
}

var errInvalidAuthCodeURL = errors.New("invalid auth code url")

var (
	//go:embed certs/tmp.key
	keyPEM []byte
	//go:embed certs/tmp.crt
	certPEM []byte
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

	if redirectURL.Scheme == "" {
		return errInvalidAuthCodeURL
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

	var listener net.Listener

	switch redirectURL.Scheme {
	case "https":
		cert, err := tls.X509KeyPair(certPEM, keyPEM)
		if err != nil {
			return fmt.Errorf("failed to load x509 key pair: %w", err)
		}

		cfg := &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
		}

		listener, err = tls.Listen("tcp", addr, cfg)
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}

	case "http":
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}
	}

	go func() {
		if err := http.Serve(listener, nil); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func generateToken(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	code := make(chan string, 1)

	if err := fetchCode(code, config); err != nil {
		close(code)

		return nil, err
	}

	tok, err := config.Exchange(ctx, <-code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange: %w", err)
	}

	return tok, nil
}
