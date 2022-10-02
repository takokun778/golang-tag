package csv_test

import (
	"os"
	"reflect"
	"tags/csv"
	"tags/model"
	"testing"
)

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func setup() {}

func teardown() {
	os.Remove("data/test.csv")
}

func Test(t *testing.T) {
	t.Parallel()

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

	if err := csv.Write("test", tags); err != nil {
		t.Fatal(err)
	}

	got, err := csv.Read("test")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, tags) {
		t.Errorf("got %v, want %v", got, tags)
	}
}
