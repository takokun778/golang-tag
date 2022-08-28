package main

import (
	"log"
	"os"
	"tags/client"
	"tags/database"
	"tags/model"

	"github.com/urfave/cli/v2"
)

func main() {
	log.Println("start app...")

	app := &cli.App{
		Name:  "tag",
		Usage: "github tags notification app",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "owner",
				Aliases:  []string{"o"},
				Required: true,
				Usage:    "github repository owner",
			},
			&cli.StringFlag{
				Name:     "repository",
				Aliases:  []string{"r"},
				Required: true,
				Usage:    "github repository name",
			},
			&cli.StringFlag{
				Name:     "channel",
				Aliases:  []string{"c"},
				Required: true,
				Usage:    "slack channel name",
			},
		},
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx *cli.Context) error {
	owner := ctx.String("owner")

	repository := ctx.String("repository")

	channel := ctx.String("channel")

	src, err := database.SelectAll(ctx.Context, owner, repository)
	if err != nil {
		return err
	}

	dst, err := client.ListTags(ctx.Context, owner, repository)
	if err != nil {
		return err
	}

	tags := model.Take(dst, src)

	if len(tags) != 0 {
		if err := database.BulkInsert(ctx.Context, tags); err != nil {
			return err
		}
	} else {
		log.Println("not found new tag")
	}

	if err := client.PostMessage(ctx.Context, channel, owner, repository, tags); err != nil {
		return err
	}

	return nil
}
