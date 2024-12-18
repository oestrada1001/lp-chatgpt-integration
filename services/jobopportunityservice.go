package services

import (
	"database/sql"
	"fmt"
	"log"
	"oestrada1001/lp-chatgpt-integration/database"
	"oestrada1001/lp-chatgpt-integration/models"
)

func GetJobOpportunity(query string) (models.JobOpportunity, error) {
	var jobOpportunity models.JobOpportunity
	db := database.Connection()
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

func FetchJobOpportunities(query string) ([]models.JobOpportunity, error) {
	db := database.Connection()
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Database query error:", err)
		return nil, err
	}
	defer rows.Close()
	var jobOpportunities []models.JobOpportunity
	for rows.Next() {
		var job models.JobOpportunity
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
