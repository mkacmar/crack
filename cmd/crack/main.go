package main

import (
	"os"

	"github.com/mkacmar/crack/internal/cli"
)

func main() {
	app := cli.New()
	os.Exit(app.Run(os.Args))
}
