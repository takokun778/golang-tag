package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"tags/client"
	"tags/database"
	"tags/model"

	"github.com/slack-go/slack"
)

func main() {
	log.Println("start app...")

	ctx := context.Background()

	src, err := database.SelectAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dst, err := client.ListTags(ctx, "golang", "go")
	if err != nil {
		log.Fatal(err)
	}

	tags := model.Take(dst, src)

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
