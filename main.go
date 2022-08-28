package main

import (
	"context"
	"log"
	"tags/client"
	"tags/database"
	"tags/model"
)

func main() {
	log.Println("start app...")

	ctx := context.Background()

	owner := "golang"

	repository := "go"

	channel := "golang-tag"

	src, err := database.SelectAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dst, err := client.ListTags(ctx, owner, repository)
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

	if err := client.PostMessage(ctx, channel, owner, repository, tags); err != nil {
		log.Fatal(err)
	}
}
