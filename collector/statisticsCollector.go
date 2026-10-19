package collector

import (
	"large-file-processing/model"
	"sync"
)

type StatisticsCollector struct {
	mu          sync.Mutex
	Departments map[string]model.DepartmentStats
}

func NewStatisticsCollector() *StatisticsCollector {
	return &StatisticsCollector{
		mu: sync.Mutex{},
		Departments: make(map[string]model.DepartmentStats),
	}
}

func (sc *StatisticsCollector) Add(deptName string, salary float64) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	ds, found := sc.Departments[deptName]
	if !found {
		sc.Departments[deptName] = model.DepartmentStats{EmployeeCount: 1, TotalSalary: salary}
		return
	}
	ds.Add(salary)
	sc.Departments[deptName] = ds
}
