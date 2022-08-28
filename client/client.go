package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"tags/model"

	"github.com/google/go-github/github"
	"github.com/slack-go/slack"
)

func ListTags(ctx context.Context, owner, repository string) ([]model.Tag, error) {
	var res []model.Tag

	client := github.NewClient(nil)

	repo := fmt.Sprintf(model.RepoFormat, owner, repository)

	opts := &github.ListOptions{
		PerPage: 100,
		Page:    1,
	}

	for {
		tags, _, err := client.Repositories.ListTags(ctx, owner, repository, opts)
		if err != nil {
			return nil, err
		}

		for _, tag := range tags {
			res = append(res, model.Tag{
				Name: *tag.Name,
				Repo: repo,
			})
		}

		if len(tags) < 100 {
			break
		}

		opts.Page += 1
	}

	return res, nil
}

func PostMessage(ctx context.Context, owner, repository string, tags []model.Tag) error {
	tkn := os.Getenv("SLACK_TOKEN")

	if tkn == "" {
		log.Println("slack token is not set")
		return nil
	}

	channel := os.Getenv("SLACK_CHANNEL")

	if channel == "" {
		log.Println("slack channel is not set")
		return nil
	}

	client := slack.New(tkn)

	for _, tag := range tags {
		url := fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", owner, repository, tag.Name)

		msg := fmt.Sprintf("released %s\n\n %s", tag.Name, url)

		if _, _, err := client.PostMessage(channel, slack.MsgOptionText(msg, true)); err != nil {
			return err
		}
	}

	return nil
}
