package builder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/mijime/beareq/pkg/beareq"
	"github.com/pelletier/go-toml"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/oauth2"
)

type ClientBuilder struct {
	Profile      string
	ProfilesPath string
	TokenDir     string
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

func (b *ClientBuilder) fetchConfigByProfile() (*oauth2.Config, error) {
	confp, err := os.Open(b.ProfilesPath)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("not found oauth config:", err)

		return &oauth2.Config{}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open oauth config: %w", err)
	}

	config := make(map[string]oauth2.Config)
	if err := toml.NewDecoder(confp).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode oauth config: %w", err)
	}

	c := config[b.Profile]

	return &c, nil
}

type profileClient struct {
	client *http.Client
	closer func() error
}

func (c *profileClient) Client() *http.Client {
	return c.client
}

func (c *profileClient) Close() error {
	return c.closer()
}

func (b *ClientBuilder) BuildClient(ctx context.Context) (beareq.Client, error) {
	config, err := b.fetchConfigByProfile()
	if err != nil {
		return nil, err
	}

	tok, err := b.fetchToken(config)
	if err != nil {
		return nil, err
	}

	tokSrc := config.TokenSource(ctx, tok)
	cli := oauth2.NewClient(ctx, tokSrc)

	return &profileClient{
		client: cli,
		closer: func() error {
			return b.saveToken(tokSrc)
		},
	}, nil
}

func (b ClientBuilder) saveToken(tokSrc oauth2.TokenSource) error {
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

func (b ClientBuilder) fetchToken(config *oauth2.Config) (*oauth2.Token, error) {
	tokenPath := path.Join(b.TokenDir, b.Profile+".json")

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

var errInvalidAuthCodeURL = errors.New("invalid auth code url")

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
		close(code)

		return nil, err
	}

	ctx := context.Background()

	tok, err := config.Exchange(ctx, <-code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange: %w", err)
	}

	return tok, nil
}
