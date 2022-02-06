package main

import (
	"fmt"
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
	return rate()
}
