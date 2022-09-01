package database_test

import (
	"context"
	"log"
	"os"
	"reflect"
	"tags/database"
	"tags/model"
	"testing"
)

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func setup() {
	ctx := context.Background()
	if err := database.DropTable(ctx); err != nil {
		log.Fatal(err)
	}
	if err := database.CreateTable(ctx); err != nil {
		log.Fatal(err)
	}
	if err := database.CreateIndex(ctx); err != nil {
		log.Fatal(err)
	}
}

func teardown() {
	ctx := context.Background()
	if err := database.DropTable(ctx); err != nil {
		log.Fatal(err)
	}
}

func Test(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tags := []model.Tag{
		{
			Name: "1",
			Repo: "test/test",
		},
		{
			Name: "2",
			Repo: "test/test",
		},
	}

	if err := database.BulkInsert(ctx, tags); err != nil {
		t.Fatal(err)
	}

	got, err := database.SelectAll(ctx, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, tags) {
		t.Errorf("SelectAll = %v, want %v", got, tags)
	}

	if err := database.Delete(ctx, "1"); err != nil {
		t.Fatal(err)
	}

	if err := database.Delete(ctx, "2"); err != nil {
		t.Fatal(err)
	}

	got, err = database.SelectAll(ctx, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(len(got), 0) {
		t.Errorf("SelectAll length = %v, want %v", len(got), 0)
	}
}
