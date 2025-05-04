package main

import (
	"context"
	"fmt"
	"os"

	"github.com/geokaralis/apigen/pkg/cmd/create"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "apigen",
		Usage: "a cli tool for generating type-safe typescript api clients from openapi specifications",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create resources",
				Commands: []*cli.Command{
					{
						Name:      "client",
						Usage:     "generate a typescript api client",
						Action:    create.Client,
						ArgsUsage: "<schema>",
						Category:  "create",
						Arguments: []cli.Argument{
							&cli.StringArgs{
								Name:      "schema",
								Min:       1,
								Max:       1,
								UsageText: "path to openapi schema file",
							},
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "output",
								Aliases:  []string{"o"},
								Usage:    "output directory for generated client",
								Required: true,
							},
							&cli.StringFlag{
								Name:  "features",
								Usage: "comma-separated list of features to enable (msw,zod)",
								Value: "",
							},
						},
					},
				},
			},
		},
	}

	err := app.Run(context.Background(), os.Args)

	if err != nil {
		fmt.Println()
		fmt.Println()
		fmt.Println(err.Error())
		fmt.Println()
		os.Exit(1)
	}
}
