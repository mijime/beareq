package main

import (
	"context"
	"flag"
	"log"

	"github.com/mijime/beareq/v2/pkg/beareq"
	cbuilder "github.com/mijime/beareq/v2/pkg/client/builder"
	rbuilder "github.com/mijime/beareq/v2/pkg/request/builder"
	"github.com/mijime/beareq/v2/pkg/response/handler"
)

func main() {
	cb := cbuilder.NewClientBuilder()
	flag.StringVar(&cb.Profile, "profile", cb.Profile, "")
	flag.StringVar(&cb.ProfilesPath, "profiles", cb.ProfilesPath, "")
	flag.StringVar(&cb.TokenDir, "tokens", cb.TokenDir, "")
	flag.BoolVar(&cb.RefreshToken, "refresh-token", cb.RefreshToken, "")

	rb := rbuilder.NewRequestBuilder()
	flag.Var(&rb.Method, "request", "Specify request command to use")
	flag.Var(&rb.Method, "X", "Specify request command to use")
	flag.Var(&rb.Header, "header", "Pass custom header(s) to server")
	flag.Var(&rb.Header, "H", "Pass custom header(s) to server")
	flag.Var(&rb.Data, "data", "HTTP POST data")
	flag.Var(&rb.Data, "d", "HTTP POST data")
	flag.Var(&rb.JSONObject, "jo", "")

	rh := handler.NewResponseHandler()
	flag.Var(&rh.JSONQuery, "jq", "")
	flag.BoolVar(&rh.Verbose, "verbose", rh.Verbose, "")
	flag.BoolVar(&rh.Fail, "fail", rh.Fail, "Fail silently (no output at all) on HTTP errors")

	flag.Parse()

	ctx := context.Background()

	err := beareq.Run(ctx, cb, rb, rh, flag.Args()...)
	if err != nil {
		log.Fatal(err)
	}
}
