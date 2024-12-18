package models

type JobOpportunity struct {
	Id                     int    `json:"id"`
	Title                  string `json:"title"`
	Description            string `json:"description"`
	CompanyName            string `json:"company_name"`
	ErrorMessage           string `json:"error_message"`
	HardSkillProcessStatus string `json:"hard_skill_process_status"`
}
