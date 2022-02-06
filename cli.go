package main

import (
	"fmt"
)

type CLI struct {
	Version versionCmd `cmd:"" help:"Show version information"`
	Rate    rateCmd    `cmd:"" default:"1" help:"Get current Dockerhub rate (default command)"`
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
