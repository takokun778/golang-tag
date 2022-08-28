package client

import (
	"context"
	"tags/model"

	"github.com/google/go-github/github"
)

var client *github.Client

func init() {
	client = github.NewClient(nil)
}

func ListTags(ctx context.Context, owner, repository string) ([]model.Tag, error) {
	var res []model.Tag

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
			})
		}

		if len(tags) < 100 {
			break
		}

		opts.Page += 1
	}

	return res, nil
}
