package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/slack-go/slack"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

type Tag struct {
	Name string `bun:"name"`
}

func init() {
	dsn := os.Getenv("DATABASE_URL")

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	if _, err := db.NewCreateTable().Model((*Tag)(nil)).Exec(context.Background()); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Println(err.Error())
		} else {
			log.Fatal(err.Error())
		}
	}

	DB = db
}

func main() {
	log.Println("start app...")

	ctx := context.Background()

	var src []Tag

	if err := DB.NewSelect().Model(&src).Scan(ctx); err != nil {
		log.Fatal(err.Error())
	}

	var dst []Tag

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
			dst = append(dst, Tag{
				Name: *tag.Name,
			})
		}

		if len(tags) < 100 {
			break
		}

		opts.Page += 1
	}

	tags := Take(dst, src)

	if len(tags) != 0 {
		if _, err := DB.NewInsert().Model(&tags).Exec(ctx); err != nil {
			log.Fatal(err.Error())
		}
	} else {
		log.Println("Not found new tag")
	}

	tkn := os.Getenv("SLACK_TOKEN")

	if tkn == "" {
		log.Fatal("slack token is not set")
	}

	channel := "golang-tag"

	sc := slack.New(tkn)

	for _, tag := range tags {
		msg := fmt.Sprintf("Added %s tag! \n\n %s", tag.Name, "https://github.com/golang/go/releases/tag/"+tag.Name)
		if _, _, err := sc.PostMessage(channel, slack.MsgOptionText(msg, true)); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func Take(from, target []Tag) []Tag {
	result := from
	for _, i := range target {
		list := make([]Tag, 0)
		for _, j := range result {
			if i != j {
				list = append(list, j)
			}
		}
		result = list
	}
	return result
}
