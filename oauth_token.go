package main

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
	"github.com/toqueteos/webbrowser"
	"golang.org/x/oauth2"
)

var ErrInvalidAuthCodeURL = errors.New("invalid auth code url")

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
		return ErrInvalidAuthCodeURL
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

func saveToken(tokSrc oauth2.TokenSource, opts option) error {
	if err := os.MkdirAll(opts.TokenDir, 0o700); err != nil {
		return fmt.Errorf("failed to create token dir: %w", err)
	}

	tokenPath := path.Join(opts.TokenDir, opts.Profile+".json")

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
