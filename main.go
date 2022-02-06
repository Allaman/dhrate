package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/tidwall/gjson"
)

const (
	TOKEN_URL     = "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull"
	RATELIMIT_URL = "https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest"
)

var Version = "dev"

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("dhrate"),
		kong.Description("Get your current Dockerhub rate"))
	err := ctx.Run()
	if err != nil {
		panic(err)
	}
}

// getRateLimit returns the value of the headers 'ratelimit-limit' and 'ratelimit-remaining' as string
func getRateLimit(c *http.Client, token string) (string, error) {
	var bearer = "Bearer " + token
	req, err := http.NewRequest("HEAD", RATELIMIT_URL, nil)
	if err != nil {
		return "", errors.New("error while creating request")
	}
	req.Header.Add("Authorization", bearer)
	resp, err := c.Do(req)
	if err != nil {
		return "", errors.New("error while requesting rate limit")
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return fmt.Sprintf("{\"rate limit\": \"%s\",\n\"rate limit remaining\": \"%s\"}", resp.Header.Get("ratelimit-limit"), resp.Header.Get("ratelimit-remaining")), nil
	} else {
		fmt.Printf("Status is: %s\n", resp.Status)
		return "", errors.New("error while getting token")
	}
}

// getToken returns a token as string for an unauthenticated user
func getToken(c *http.Client) (string, error) {
	resp, err := c.Get(TOKEN_URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", errors.New("error while parsing body")
		}
		return getKey(string(b), "token"), nil
	} else {
		fmt.Printf("Status is: %s\n", resp.Status)
		return "", errors.New("error while getting token")
	}
}

// getKey returns the key of a json as string
func getKey(json string, key string) string {
	return gjson.Get(json, key).String()
}
