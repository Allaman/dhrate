package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

type client struct {
	httpClient   *http.Client
	tokenURL     string
	ratelimitURL string
}

func createClient() client {
	return client{httpClient: &http.Client{Timeout: time.Duration(2) * time.Second}, tokenURL: TOKEN_URL, ratelimitURL: RATELIMIT_URL}
}

// getRateLimit returns the value of the headers 'ratelimit-limit' and 'ratelimit-remaining' as json formatted string
func (c *client) getRateLimit(token string) (string, error) {
	req, err := http.NewRequest("HEAD", c.ratelimitURL, nil)
	if err != nil {
		return "", errors.New("error while creating request to get rate limit")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", errors.New("error while performing http request to get rate limit")
	}
	defer resp.Body.Close()
	if isHTTPResponseCode200(resp) {
		return fmt.Sprintf("{\"rate limit\": \"%s\",\n\"rate limit remaining\": \"%s\"}", resp.Header.Get("ratelimit-limit"), resp.Header.Get("ratelimit-remaining")), nil
	} else {
		log.Printf("Status is: %s\n", resp.Status)
		return "", errors.New("rate limit http response was not 200")
	}
}

// getClientToken returns a token as string for an unauthenticated user
func (c *client) getClientToken() (string, error) {
	resp, err := c.httpClient.Get(c.tokenURL)
	if err != nil {
		return "", errors.New("error while performing http request to get token")
	}
	defer resp.Body.Close()
	if isHTTPResponseCode200(resp) {
		b := readRespBody(resp)
		return getJSONKey(string(b), "token"), nil
	} else {
		log.Printf("Status is: %s\n", resp.Status)
		return "", errors.New("token http response was not 200")
	}
}

func getJSONKey(json string, key string) string {
	return gjson.Get(json, key).String()
}

func isHTTPResponseCode200(resp *http.Response) bool {
	return resp.StatusCode == 200
}

func readRespBody(resp *http.Response) []byte {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error while reading response body: %v", err)
	}
	return b
}
