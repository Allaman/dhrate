package main

import (
	"fmt"
	"net/http"
	"time"
)

type CLI struct {
	Version versionCmd `cmd:"" help:"Show version information"`
	Rate    rateCmd    `cmd:"" help:"Get current Dockerhub rate"`
}

type versionCmd struct {
	Version string
}

func (c *versionCmd) Run() error {
	fmt.Println(Version)
	return nil
}

type rateCmd struct{}

func (c *rateCmd) Run() error {
	client := http.Client{Timeout: time.Duration(2) * time.Second}
	token, err := getToken(&client)
	if err != nil {
		return err
	}
	limit, err := getRateLimit(&client, token)
	if err != nil {
		return err
	}
	fmt.Println(limit)
	return nil
}
