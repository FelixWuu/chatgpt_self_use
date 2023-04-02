package main

import (
	"github.com/alecthomas/kong"
)

var CLI struct {
	Verbose bool   `help:"Verbose mode."`
	Config  string `help:"Config file." name:"config" type:"file" default:"config.toml"`
}

func main() {
	kong.Parse(&CLI)
}
