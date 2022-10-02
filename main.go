package main

import (
	"log"
	"os"
	"tags/client"
	"tags/csv"
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

	src, err := csv.Read(repository)
	if err != nil {
		return err
	}

	dst, err := client.ListTags(ctx.Context, owner, repository)
	if err != nil {
		return err
	}

	diff := model.Take(dst, src)

	if err := csv.Write(repository, dst); err != nil {
		return err
	}

	if len(diff) == 0 {
		log.Println("not found new tag")

		return nil
	}

	if err := client.PostMessage(ctx.Context, owner, repository, diff); err != nil {
		return err
	}

	return nil
}
