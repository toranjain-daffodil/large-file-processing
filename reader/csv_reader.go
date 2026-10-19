package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"large-file-processing/model"
	"os"
)

type RecordHandler func(model.RawEmployee) error

func ReadCSV(path string, handle RecordHandler) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open csv: %w", err)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)

	header, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("read csv header: %w", err)
	}

	expected := []string{
		"employee_id",
		"name",
		"department",
		"salary",
	}

	if len(header) != len(expected) {
		return fmt.Errorf("invalid csv header: expected %d columns got %d", len(expected), len(header))
	}

	for i, column := range expected {
		if header[i] != column {
			return fmt.Errorf("invalid csv header at column %d: expected %q, got %q", i+1, column, header[i])
		}
	}

	recordNumber := 1

	for {
		fields, err := csvReader.Read()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return fmt.Errorf("read csv record %d: %w", recordNumber+1, err)
		}

		recordNumber++
		if len(fields) != len(expected) {
			return fmt.Errorf("record %d, expected %d columns, got %d", recordNumber, len(expected), len(fields))
		}

		employee := model.RawEmployee{
			RecordNumber: recordNumber,
			EmployeeID: fields[0],
			Name: fields[1],
			Department: fields[2],
			Salary: fields[3],
		}

		if err := handle(employee); err != nil {
			return fmt.Errorf("process record %d: %w", recordNumber, err)
		}
	}


}