package main

import (
	"fmt"
	"large-file-processing/model"
	"large-file-processing/reader"
	"log"
)

func main() {
	err := reader.ReadCSV("assets/employees.csv", func(employee model.RawEmployee) error {
		fmt.Printf("%+v\n", employee)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}