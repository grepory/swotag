package main

import (
	"github.com/grepory/swotag/internal"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "swotag"
	app.Usage = "Fetch entities from Solarwinds, filter them by attributes, and update them with new tags"

	app.Commands = []*cli.Command{
		{
			Name:  "tag",
			Usage: "tag entities given a basic filter",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "entity-type",
					Aliases: []string{"e"},
					Usage:   "Set the type of entities to tag",
				},
				&cli.StringFlag{
					Name:    "filter",
					Aliases: []string{"f"},
					Usage:   "Set the filter to apply (attribute:text)",
				},
				&cli.StringSliceFlag{
					Name:    "tag",
					Aliases: []string{"t"},
					Usage:   "Set tag, can be applied multiple times (tag:value)",
				},
			},
			Action: func(ctx *cli.Context) error {
				return internal.TagEntities(
					ctx.String("entity-type"),
					ctx.String("filter"),
					ctx.StringSlice("tag"))
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
