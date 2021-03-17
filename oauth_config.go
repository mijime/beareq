package main

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
	"golang.org/x/oauth2"
)

func fetchConfigByProfile(opts option) (*oauth2.Config, error) {
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
