package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"tags/model"

	_ "github.com/lib/pq"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var database *bun.DB

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
		Model((*model.Tag)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTable(ctx context.Context) error {
	_, err := database.NewDropTable().
		Model((*model.Tag)(nil)).
		IfExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func BulkInsert(ctx context.Context, tags []model.Tag) error {
	if _, err := database.NewInsert().Model(&tags).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func SelectAll(ctx context.Context) ([]model.Tag, error) {
	var tags []model.Tag

	if err := database.NewSelect().Model(&tags).Scan(ctx); err != nil {
		return nil, err
	}

	return tags, nil
}
