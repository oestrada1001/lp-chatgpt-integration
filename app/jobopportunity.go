package app

import (
	"database/sql"
	"fmt"
	"log"
)

type JobOpportunity struct {
	Id                     int    `json:"id"`
	Title                  string `json:"title"`
	Description            string `json:"description"`
	CompanyName            string `json:"company_name"`
	ErrorMessage           string `json:"error_message"`
	HardSkillProcessStatus string `json:"hard_skill_process_status"`
}

func GetJobOpportunity(query string) (JobOpportunity, error) {
	var jobOpportunity JobOpportunity
	db := DatabaseConnection()
	defer db.Close()

	row := db.QueryRow(query)

	// Scan the result into the JobOpportunity struct
	err := row.Scan(
		&jobOpportunity.Id,
		&jobOpportunity.Title,
		&jobOpportunity.Description,
		&jobOpportunity.CompanyName,
		&jobOpportunity.HardSkillProcessStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return jobOpportunity, fmt.Errorf("no job opportunity found with query: %s", query)
		}
		return jobOpportunity, fmt.Errorf("error querying job opportunity: %v", err)
	}

	return jobOpportunity, nil
}

func FetchJobOpportunities(query string) ([]JobOpportunity, error) {
	db := DatabaseConnection()
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Database query error:", err)
		return nil, err
	}
	defer rows.Close()
	var jobOpportunities []JobOpportunity
	for rows.Next() {
		var job JobOpportunity
		err := rows.Scan(
			&job.Title,
			&job.Description,
			&job.CompanyName,
			&job.ErrorMessage,
			&job.HardSkillProcessStatus,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		jobOpportunities = append(jobOpportunities, job)
	}
	return jobOpportunities, nil
}
