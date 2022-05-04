package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "etxe",
		Usage: "infrastructure as code done right",
		Action: func(c *cli.Context) error {
			fmt.Println("boom! I say!")

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
