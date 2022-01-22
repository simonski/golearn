package main

import (
	"fmt"
	"os"

	"github.com/simonski/golearn/learn/grpc"
	"github.com/simonski/golearn/learn/http"
	"github.com/simonski/golearn/learn/sqlite"
	"github.com/simonski/goutils"
)

func main() {
	cli := goutils.NewCLI(os.Args)
	command := cli.GetCommand()
	if command == "db" {
		app := sqlite.NewApp()
		app.HandleInput(command, cli)
	} else if command == "grpc" {
		app := grpc.NewApp()
		app.HandleInput(command, cli)
	} else if command == "http" {
		app := http.NewApp()
		app.HandleInput(command, cli)
	} else {
		fmt.Println("Error, usage: ./learn <COMMAND> (where command is db, grpc, http)")
		os.Exit(1)
	}
}
