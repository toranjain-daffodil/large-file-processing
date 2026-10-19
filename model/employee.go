package model

type Employee struct {
	EmployeeID string
	Name string
	Department string
	Salary float64
}

type RawEmployee struct {
	RecordNumber int
	EmployeeID string
	Name string
	Department string
	Salary string
}