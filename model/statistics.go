package model

type DepartmentStats struct {
	EmployeeCount int
	TotalSalary float64
}

func (ds *DepartmentStats) Add(salary float64) {
	ds.EmployeeCount++
	ds.TotalSalary += salary
}