package main

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	_ "log"
	"net/http"
	"oestrada1001/lp-chatgpt-integration/app"
	"oestrada1001/lp-chatgpt-integration/chatgpt"
)

func main() {
	r := httprouter.New()

	r.GET("/fetch-job-opportunities", FetchJobOpportunities)
	r.GET("/fetch-or-create-proficiency-levels", RunFetchOrCreateProficiencyLevels)
	r.GET("/fetch-or-create-hard-skill-types", RunFetchOrCreateHardSkillTypes)
	r.GET("/process-job-opportunities", ProcessJobOpportunities)

	fmt.Println("Server started on port 8082")
	err := http.ListenAndServe(":8082", r)
	if err != nil {
		fmt.Println("Error starting server:")
		fmt.Println(err)
		return
	}
}

func ProcessJobOpportunities(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	jobOpportunities, err := app.FetchJobOpportunities("SELECT title, description, company_name, error_message, hard_skill_process_status FROM job_opportunities")
	if err != nil {
		http.Error(rw, "Failed to query database", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(rw, jobOpportunities)
}

func RunFetchOrCreateHardSkillTypes(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	hardSkills := []app.HardSkillType{
		{Label: "Go", Value: "go", Description: "Go is a programming language"},
		{Label: "Gl1o", Value: "gkjlo", Description: "adkGo is a programming language"},
		{Label: "Goi", Value: "jklgo", Description: "Gasdko is a programming language"},
	}
	hardSkillTypes, err := app.ReplaceAndFetchHardSkillTypes(hardSkills)
	if err != nil {
		http.Error(rw, "Failed to query database", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(rw, hardSkillTypes)
}

func RunFetchOrCreateProficiencyLevels(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	proficiencyLevels := []app.ProficiencyLevel{
		{Label: "2o", Value: "go", Description: "Go is a programming language"},
		{Label: "dl1o", Value: "gkjlo", Description: "adkGo is a programming language"},
		{Label: "doi", Value: "jklgo", Description: "Gasdko is a programming language"},
	}

	proficiencyLevels, err := app.ReplaceAndFetchProficiencyLevels(proficiencyLevels)
	if err != nil {
		http.Error(rw, "Failed to query database", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(rw, proficiencyLevels)
}

func FetchJobOpportunities(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	jobOpportunities, err := app.FetchJobOpportunities("SELECT title, description, company_name, COALESCE(error_message, '') AS error_message, hard_skill_process_status FROM job_opportunities")
	if err != nil {
		http.Error(rw, "Failed to query database", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(rw).Encode(jobOpportunities)
	if err != nil {
		return
	}

	chatgpt.JobAssistant(jobOpportunities[0])

	fmt.Fprintln(rw, jobOpportunities)
}
