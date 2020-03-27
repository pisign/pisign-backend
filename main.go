package main

import (
	"log"

	"github.com/pisign/pisign-backend/api"
	"github.com/pisign/pisign-backend/server"

	"github.com/alecthomas/kong"
)

// Command line inspired by https://github.com/alecthomas/kong#attach-a-run-error-method-to-each-command
type Context struct {
}

type RunCmd struct {
}

func (c *RunCmd) Run(ctx *Context) error {
	log.Printf("Run subcommand!\n")
	server.StartLocalServer(9000)
	return nil
}

type CreateCmd struct {
	Name string `arg,help:"New api name"`
}

func (c *CreateCmd) Run(ctx *Context) error {
	log.Printf("Create subcommand!\n")
	api.AutoCreate(c.Name)
	return nil
}

var cli struct {
	Run    RunCmd    `cmd,help:"Run the go server" default:"1"`
	Create CreateCmd `cmd,help:"Create a new api type"`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
