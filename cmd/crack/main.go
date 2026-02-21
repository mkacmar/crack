package main

import (
	"os"

	"go.kacmar.sk/crack/internal/cli"
)

func main() {
	app := cli.New()
	os.Exit(app.Run(os.Args))
}
