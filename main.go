package main

import (
	"fmt"
	"large-file-processing/model"
	"large-file-processing/reader"
	"log"
)

func main() {
	err := reader.ReadCSV("assets/employees.csv", recordHandler)

	if err != nil {
		log.Fatal(err)
	}
}

func recordHandler(employee model.RawEmployee) error {
	fmt.Printf("%+v\n", employee)
	return nil
}