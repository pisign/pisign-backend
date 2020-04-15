// Package main is the main point of entry for the pisign backend server
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
	Port int `name:"port" short:"p" help:"Port to run server on" default:"9000"`
}

func (c *RunCmd) Run(ctx *Context) error {
	log.Printf("Run subcommand with port = %v!\n", c.Port)
	return server.StartLocalServer(c.Port)
}

type CreateCmd struct {
	Name string `arg help:"New api name"`
}

func (c *CreateCmd) Run(ctx *Context) error {
	log.Printf("Create subcommand!\n")
	return api.AutoCreate(c.Name)
}

var cli struct {
	Run    RunCmd    `cmd help:"Run the go server" default="1"`
	Create CreateCmd `cmd help:"Create a new api type"`
}

// main is the entry function to the server. It can be run with different sub commands, as follows:
/*
Usage: main <command>

Flags:
  --help    Show context-sensitive help.

Commands:
  run
    Run the go server

  create <name>
    Create a new api type

Run "main <command> --help" for more information on a command.
*/
func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
