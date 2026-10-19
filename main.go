package main

import (
	"context"
	"large-file-processing/collector"
	"large-file-processing/pipeline"
	"log"
)

func main() {
	stats := collector.NewStatisticsCollector()

	err := pipeline.Run(context.Background(), "./assets/employees.csv", "./output/rejected.jsonl", stats)
	if err != nil {
		log.Printf("Processing failed: %v", err)
		log.Println("Collected statistics may be incomplete.")
		return
	}

	log.Println("Processing completed successfully")
	for department, departmentStats := range stats.Departments {
		var averageSalary float64
		if departmentStats.EmployeeCount > 0 {
			averageSalary = departmentStats.TotalSalary /
				float64(departmentStats.EmployeeCount)
		}

		log.Printf(
			"Department: %s | Employees: %d | Average salary: %.2f",
			department,
			departmentStats.EmployeeCount,
			averageSalary,
		)
	}
}