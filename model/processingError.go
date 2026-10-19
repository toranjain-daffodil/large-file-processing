package model

type ProcessingError struct {
	RecordNumber int
	EmployeeID string
	ErrorCategory string
	ErrorMessage string
}