package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jvardilos/ccbserving/app/commands/serving"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "ccbserving",
		Usage: "Utility to get a serving trend of each member of a given group",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "client_id",
				Aliases:  []string{"c"},
				Usage:    "CHMS client-id",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "secret",
				Aliases:  []string{"s"},
				Usage:    "CHMS secret",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "subdomain",
				Aliases:  []string{"d"},
				Usage:    "CHMS subdomain",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "group",
				Aliases:  []string{"g"},
				Usage:    "Group name",
				Required: true,
			},
		},
		Action: serving.GetServing,
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
	}
}
