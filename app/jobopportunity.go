package app

import (
	"database/sql"
	"log"
)

type JobOpportunity struct {
	Title                  string
	Description            string
	CompanyName            string
	ErrorMessage           sql.NullString
	HardSkillProcessStatus string
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
