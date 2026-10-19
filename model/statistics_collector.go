package model

import "sync"

type StatisticsCollector struct {
	mu sync.Mutex
	Departments map[string]departmentStats
}

type departmentStats struct {
	employeeCount int
	totalSalary float64
}

func (ds *departmentStats) add(salary float64) {
	ds.employeeCount++
	ds.totalSalary += salary
}

func (sc *StatisticsCollector) Add(deptName string, salary float64) {
	sc.mu.Lock()

	ds, found := sc.Departments[deptName]
	if !found {
		sc.Departments[deptName] = departmentStats{employeeCount: 1, totalSalary: salary}
		return
	}
	ds.add(salary)
	sc.mu.Unlock()
}
