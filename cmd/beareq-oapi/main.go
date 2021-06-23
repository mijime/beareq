package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/agnivade/levenshtein"
	"github.com/mijime/beareq/pkg/beareq"
	cbuilder "github.com/mijime/beareq/pkg/client/builder"
	"github.com/mijime/beareq/pkg/openapi"
	"github.com/mijime/beareq/pkg/response/handler"
	"github.com/pelletier/go-toml"
)

func main() {
	cb := cbuilder.NewClientBuilder()
	flag.StringVar(&cb.Profile, "profile", cb.Profile, "")
	flag.StringVar(&cb.ProfilesPath, "profiles", cb.ProfilesPath, "")
	flag.StringVar(&cb.TokenDir, "tokens", cb.TokenDir, "")

	rh := handler.NewResponseHandler()
	flag.Var(&rh.JSONQuery, "jq", "")
	flag.BoolVar(&rh.Verbose, "verbose", rh.Verbose, "")
	flag.BoolVar(&rh.Fail, "fail", rh.Fail, "Fail silently (no output at all) on HTTP errors")

	flag.Parse()

	config, err := fetchConfigByProfile(profileConfig{
		ProfilesPath: cb.ProfilesPath,
		Profile:      cb.Profile,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to fecth config: %w", err))
	}

	cmds := make(map[string]*openapi.Operation)

	for _, oc := range config.OpenAPI {
		for _, specPath := range oc.Specs {
			fp, err := os.Open(specPath)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to open openapi spec: %w", err))
			}

			subcmds, err := openapi.GenerateOperation(oc.BaseURL, fp)
			if err != nil {
				fp.Close()
				log.Fatal(fmt.Errorf("failed to generate openapi: %w", err))
			}

			fp.Close()

			for cmdName, cmd := range subcmds {
				cmds[cmdName] = cmd
			}
		}
	}

	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, "supported subcommands:\n")

		cmdNames := make([]string, 0, len(cmds))

		for cmdName := range cmds {
			cmdNames = append(cmdNames, cmdName)
		}

		sort.Strings(cmdNames)

		for _, cmdName := range cmdNames {
			fmt.Fprintf(os.Stderr, "\t- %s\n", cmdName)
		}

		os.Exit(1)
	}

	cmd, ok := cmds[args[0]]
	if !ok {
		minD := -1
		suggestCmdName := ""

		for cmdName := range cmds {
			d := levenshtein.ComputeDistance(args[0], cmdName)
			if minD < 0 || d < minD {
				minD = d
				suggestCmdName = cmdName
			}
		}

		log.Fatalf("unsupported command: %s. the most similar command is %s", args[0], suggestCmdName)
	}

	if err := cmd.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	if err := beareq.Run(ctx, cb, cmd, rh, cmd.BaseURL); err != nil {
		log.Fatal(err)
	}
}

type profileConfig struct {
	ProfilesPath string
	Profile      string
}

type openapiConfig struct {
	OpenAPI []struct {
		BaseURL string
		Specs   []string
	}
}

func fetchConfigByProfile(b profileConfig) (openapiConfig, error) {
	confp, err := os.Open(b.ProfilesPath)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("not found profile:", err)

		return openapiConfig{}, nil
	}

	if err != nil {
		return openapiConfig{}, fmt.Errorf("failed to open profile: %w", err)
	}

	defer confp.Close()

	config := make(map[string]openapiConfig)
	if err := toml.NewDecoder(confp).Decode(&config); err != nil {
		return openapiConfig{}, fmt.Errorf("failed to decode openapi config: %w", err)
	}

	c := config[b.Profile]

	return c, nil
}
