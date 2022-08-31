package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"tags/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var database *bun.DB

func init() {
	dsn := os.Getenv("DATABASE_URL")

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
	}

	config.PreferSimpleProtocol = true

	db := stdlib.OpenDB(*config)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	database = bun.NewDB(db, pgdialect.New())
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

func CreateIndex(ctx context.Context) error {
	_, err := database.NewCreateIndex().
		Model((*model.Tag)(nil)).
		Index("name_idx").
		Column("name").
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = database.NewCreateIndex().
		Model((*model.Tag)(nil)).
		Index("repo_idx").
		Column("repo").
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func DropTable(ctx context.Context) error {
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
	_, err := database.NewInsert().
		Model(&tags).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, name string) error {
	_, err := database.NewDelete().
		Model((*model.Tag)(nil)).
		Where("name = ?", name).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func SelectAll(ctx context.Context, owner, repository string) ([]model.Tag, error) {
	var tags []model.Tag

	repo := fmt.Sprintf(model.RepoFormat, owner, repository)

	err := database.NewSelect().
		Model(&tags).
		Where("repo = ?", repo).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func Close() error {
	return database.Close()
}
