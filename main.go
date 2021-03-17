package main

import (
	"context"
	"flag"
	"log"

	"golang.org/x/oauth2"
)

func run(urls []string, opts option) error {
	config, err := fetchProfile(opts)
	if err != nil {
		return err
	}

	tok, err := fetchToken(config, opts)
	if err != nil {
		return err
	}

	ctx := context.Background()
	tokSrc := config.TokenSource(ctx, tok)
	cli := oauth2.NewClient(ctx, tokSrc)

	for _, rawurl := range urls {
		if err := doRequest(cli, rawurl, opts); err != nil {
			return err
		}
	}

	return saveToken(tokSrc, opts)
}

func main() {
	opts, err := newOption()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&opts.Request, "request", opts.Request, "Specify request command to use")
	flag.StringVar(&opts.Profile, "profile", opts.Profile, "")
	flag.StringVar(&opts.ProfilesPath, "config", opts.ProfilesPath, "")
	flag.StringVar(&opts.TokenDir, "tokens", opts.TokenDir, "")
	flag.Var(&opts.Data, "data", "HTTP POST data")
	flag.Var(&opts.Header, "header", "Pass custom header(s) to server")
	flag.BoolVar(&opts.Verbose, "verbose", opts.Verbose, "")
	flag.BoolVar(&opts.Fail, "fail", opts.Fail, "Fail silently (no output at all) on HTTP errors")
	flag.Parse()

	urls := flag.Args()

	if err := run(urls, opts); err != nil {
		log.Fatal(err)
	}
}
