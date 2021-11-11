package main

import (
	"log"
	"os"

	lakituCLI "github.com/sereneblue/lakitu/internal/cli"
	"github.com/urfave/cli/v2"
)

var version string

func main() {
	app := &cli.App{
		Name:      "lakitu-cli",
		Version:   version,
		Usage:     "A small utility to manage your cloud gaming instance",
		UsageText: "lakitu-cli [global options] command [command options] [arguments...]",
		Commands: []*cli.Command{
			{
				Name:  "bootstrap",
				Usage: "Download and install dependencies on new server",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:  "mount",
				Usage: "Manage storage (Format new volume/instance stores and attach volume from snapshot ID )",
				Action: func(c *cli.Context) error {
					snapshotId := c.Args().First()

					err := lakituCLI.MountSnapshot(snapshotId)

					if err != nil {
						log.Fatal(err)
					}

					return err
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
