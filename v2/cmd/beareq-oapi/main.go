package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"

	"github.com/mijime/beareq/v2/pkg/beareq"
	"github.com/mijime/beareq/v2/pkg/client/builder"
	"github.com/mijime/beareq/v2/pkg/openapi"
	"github.com/mijime/beareq/v2/pkg/response/handler"
	"github.com/mijime/beareq/v2/pkg/suggest"
	"github.com/pelletier/go-toml"
)

func main() {
	cb := builder.NewClientBuilder()
	flag.StringVar(&cb.Profile, "profile", cb.Profile, "")
	flag.StringVar(&cb.ProfilesPath, "profiles", cb.ProfilesPath, "")
	flag.StringVar(&cb.TokenDir, "tokens", cb.TokenDir, "")
	flag.BoolVar(&cb.RefreshToken, "refresh-token", cb.RefreshToken, "")

	rh := handler.NewResponseHandler()
	flag.Var(&rh.JSONQuery, "jq", "")
	flag.BoolVar(&rh.Fail, "fail", rh.Fail, "Fail silently (no output at all) on HTTP errors")

	verbose, _ := strconv.ParseBool(os.Getenv("BEAREQ_VERBOSE"))

	flag.BoolVar(&verbose, "verbose", verbose, "")

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
	usr, _ := user.Current()

	for _, oc := range config.OpenAPI {
		for _, specPath := range oc.Specs {
			specPath = strings.Replace(specPath, "~", usr.HomeDir, 1)

			subcmds, err := openapi.GenerateOperationFromPath(oc.BaseURL, specPath)
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
		fmt.Fprintf(os.Stderr, "unsupported subcommand %s: the more similar command is\n", args[0])

		cmdNames := make([]string, 0, len(cmds))

		for cmdName := range cmds {
			cmdNames = append(cmdNames, cmdName)
		}

		suggestCmdNames := suggest.Suggest(cmdNames, args[0], 10)

		for _, cmdName := range suggestCmdNames {
			fmt.Fprintf(os.Stderr, "\t- %s\n", cmdName)
		}

		os.Exit(1)
	}

	cmdFs := cmd.FlagSet(envPrefix)
	cmdFs.Var(&rh.JSONQuery, "jq", "")
	cmdFs.BoolVar(&verbose, "verbose", verbose, "")

	if err := cmdFs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}

	rh.Verbose = verbose
	cmd.Verbose = verbose

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
