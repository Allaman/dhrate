package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kong"
)

const (
	TOKEN_URL     = "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull"
	RATELIMIT_URL = "https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest"
)

// will be overwritten in release pipeline
var Version = "dev"

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("dhrate"),
		kong.Description("Get your current Dockerhub rate"))
	err := ctx.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

// rate prints Dockerhub rate limits to stdout
func rate() error {
	client := createClient()
	dockerHubClientToken, err := client.getClientToken()
	if err != nil {
		return err
	}
	limit, err := client.getRateLimit(dockerHubClientToken)
	if err != nil {
		return err
	}
	fmt.Println(limit)
	return nil
}
