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

	var envPrefix string

	flag.StringVar(&envPrefix, "env-prefix", "", "")

	flag.Parse()

	if len(envPrefix) == 0 {
		envPrefix = "BEAREQ_OAPI_" + cb.Profile
	}

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
			subcmds, err := openapi.GenerateOperation(oc.BaseURL, specPath)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to generate openapi: %w", err))
			}

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
		cmdNames := make([]struct {
			n string
			d int
		}, 0, len(cmds))

		for cmdName := range cmds {
			var d int

			if len(cmdName) > len(args[0]) {
				d = levenshtein.ComputeDistance(args[0], cmdName[:len(args[0])])*10 +
					levenshtein.ComputeDistance(args[0], cmdName[len(args[0]):])
			} else {
				d = levenshtein.ComputeDistance(args[0], cmdName) * 10
			}

			cmdNames = append(cmdNames, struct {
				n string
				d int
			}{n: cmdName, d: d})
		}

		sort.Slice(cmdNames, func(i, j int) bool {
			return cmdNames[i].d < cmdNames[j].d
		})

		viewSize := len(cmdNames)
		if viewSize > 5 {
			viewSize = 5
		}

		fmt.Fprintf(os.Stderr, "unsupported subcommand %s: the more similar command is\n", args[0])

		for _, cmdName := range cmdNames[:viewSize] {
			fmt.Fprintf(os.Stderr, "\t- %s\n", cmdName.n)
		}

		os.Exit(1)
	}

	if err := cmd.Parse(envPrefix, args[1:]); err != nil {
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

	c, ok := config[b.Profile]
	if !ok {
		return openapiConfig{}, fmt.Errorf("not found profile in profiles.toml: %s", b.Profile)
	}

	return c, nil
}
