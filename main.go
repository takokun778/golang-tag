package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"tags/database"

	"github.com/google/go-github/github"
	"github.com/slack-go/slack"
)

func main() {
	log.Println("start app...")

	ctx := context.Background()

	src, err := database.SelectAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var dst []database.Tag

	client := github.NewClient(nil)

	owner := "golang"

	repo := "go"

	opts := &github.ListOptions{
		PerPage: 100,
		Page:    1,
	}

	for {
		tags, _, err := client.Repositories.ListTags(ctx, owner, repo, opts)
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, tag := range tags {
			dst = append(dst, database.Tag{
				Name: *tag.Name,
			})
		}

		if len(tags) < 100 {
			break
		}

		opts.Page += 1
	}

	tags := database.Take(dst, src)

	if len(tags) != 0 {
		if err := database.BulkInsert(ctx, tags); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("not found new tag")
	}

	channel := os.Getenv("SLACK_CHANNEL")

	if channel == "" {
		log.Println("slack channel is not set")
		return
	}

	tkn := os.Getenv("SLACK_TOKEN")

	if tkn == "" {
		log.Println("slack token is not set")
		return
	}

	sc := slack.New(tkn)

	for _, tag := range tags {
		msg := fmt.Sprintf("added %s tag! \n\n %s", tag.Name, "https://github.com/golang/go/releases/tag/"+tag.Name)
		if _, _, err := sc.PostMessage(channel, slack.MsgOptionText(msg, true)); err != nil {
			log.Fatal(err.Error())
		}
	}
}
