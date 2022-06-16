package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/alecthomas/repr"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/hexbee-net/etxe/pkg/etx"
)

const (
	EtxExtension       = ".etx"
	TerraformExtension = ".tf"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}

	app := &cli.App{
		Name:                   "etxe",
		Version:                "v0.0.1",
		Copyright:              "(c) 2022 Hexbee",
		Usage:                  "infrastructure as code done right",
		UseShortOptionHandling: true,

		Commands: []*cli.Command{
			{
				Name:     "init",
				Category: "Main commands",
				Usage:    "Prepare your working directory for other commands",
				Action:   commandNotImplemented,
			},
			{
				Name:     "check",
				Category: "Main commands",
				Usage:    "Check whether the configuration is valid",
				Action:   commandCheck(logger),
			},
			{
				Name:     "plan",
				Category: "Main commands",
				Usage:    "Show changes required by the current configuration",
				Action:   commandNotImplemented,
			},
			{
				Name:     "apply",
				Category: "Main commands",
				Usage:    "Create or update infrastructure",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "force",
						Aliases:     []string{"f"},
						Usage:       "Don't ask for confirmation before applying",
						Destination: nil,
					},
				},
				Action: commandNotImplemented,
			},
			{
				Name:     "destroy",
				Category: "Main commands",
				Usage:    "Destroy previously-created infrastructure",
				Action:   commandNotImplemented,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
	}
}

func commandNotImplemented(c *cli.Context) error {
	panic("not implemented")
}

func commandCheck(logger *zap.Logger) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		var err error

		targetDir := c.Args().Get(0)
		if targetDir == "" {
			targetDir = "."
		}

		targetDir, err = filepath.Abs(targetDir)
		if err != nil {
			return fmt.Errorf("failed to resolve target path: %w", err)
		}

		logger.Debug("checking directory configuration", zap.String("target", targetDir))

		files, err := ioutil.ReadDir(targetDir)
		if err != nil {
			return fmt.Errorf("failed to list files in directory %q: %w", targetDir, err)
		}

		for _, file := range files {
			ext := filepath.Ext(file.Name())
			if !file.IsDir() && (ext == EtxExtension || ext == TerraformExtension) {
				targetPath := path.Join(targetDir, file.Name())

				reader, err := os.Open(targetPath)
				if err != nil {
					return fmt.Errorf("failed to open file %q: %w", targetPath, err)
				}

				ast, err := etx.Parse(reader)
				if err != nil {
					return fmt.Errorf("failed to parse file %q: %w", targetPath, err)
				}

				repr.Println(ast)
			}
		}

		return nil
	}
}
