package database

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var database *bun.DB

type Tag struct {
	Name string `bun:"name"`
}

func init() {
	dsn := os.Getenv("DATABASE_URL")

	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	database = db
}

func CreateTable(ctx context.Context) error {
	_, err := database.NewCreateTable().
		Model((*Tag)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTable(ctx context.Context) error {
	_, err := database.NewDropTable().
		Model((*Tag)(nil)).
		IfExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func BulkInsert(ctx context.Context, tags []Tag) error {
	if _, err := database.NewInsert().Model(&tags).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func SelectAll(ctx context.Context) ([]Tag, error) {
	var tags []Tag

	if err := database.NewSelect().Model(&tags).Scan(ctx); err != nil {
		return nil, err
	}

	return tags, nil
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
