package csv

import (
	"fmt"
	"os"
	"tags/model"

	"github.com/gocarina/gocsv"
)

const src = "data/src.csv"

func Read(filename string) ([]model.Tag, error) {
	os.Rename(fmt.Sprintf("data/%s.csv", filename), src)

	file, err := os.OpenFile(src, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tags := []model.Tag{}

	if err := gocsv.UnmarshalFile(file, &tags); err != nil {
		return nil, err
	}

	if err := os.Remove(src); err != nil {
		return nil, err
	}

	return tags, nil
}

func Write(filename string, tags []model.Tag) error {
	file, err := os.OpenFile(fmt.Sprintf("data/%s.csv", filename), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&tags, file); err != nil {
		return err
	}

	return nil
}
