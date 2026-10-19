package worker

import (
	"context"
	"strconv"
	"strings"

	"large-file-processing/collector"
	"large-file-processing/model"
)

func ProcessEmployee(
	ctx context.Context,
	raw model.RawEmployee,
	stats *collector.StatisticsCollector,
	errorCollector *collector.ErrorCollector,
) error {
	reject := func(category, message string) error {
		return errorCollector.Add(model.ProcessingError{
			RecordNumber: raw.RecordNumber,
			EmployeeID:   raw.EmployeeID,
			ErrorCategory: category,
			ErrorMessage:  message,
		})
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if strings.TrimSpace(raw.EmployeeID) == "" {
		return reject("validation", "employee ID is empty")
	}

	if strings.TrimSpace(raw.Name) == "" {
		return reject("validation", "name is empty")
	}

	department := strings.TrimSpace(raw.Department)
	if department == "" {
		return reject("validation", "department is empty")
	}

	salary, err := strconv.ParseFloat(strings.TrimSpace(raw.Salary), 64)
	if err != nil {
		return reject("validation", "salary is not numeric")
	}

	if salary < 0 {
		return reject("validation", "salary cannot be negative")
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	stats.Add(department, salary)
	return nil
}
