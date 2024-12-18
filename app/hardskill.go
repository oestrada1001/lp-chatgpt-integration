package app

import (
	_ "database/sql"
	"fmt"
	"log"
	"strings"
)

type HardSkill struct {
	Id              int
	Name            string
	Link            string
	Logo            string
	HardSkillTypeId int
}

func replaceHardSkills(hardSkills []HardSkill) ([]HardSkill, error) {
	if len(hardSkills) == 0 {
		return nil, nil
	}

	valuePlaceholders := make([]string, len(hardSkills))
	for i, hardSkill := range hardSkills {
		valuePlaceholders[i] = fmt.Sprintf("('%s','%s','%s','%s')",
			hardSkill.Name,
			hardSkill.Link,
			hardSkill.Logo,
			hardSkill.HardSkillTypeId,
		)
	}

	queryValues := strings.ReplaceAll(strings.Trim(strings.Join(valuePlaceholders, ","), ""), " ", "")
	query := fmt.Sprintf("REPLACE INTO hard_skills (name, link, logo, hard_skill_type_id) VALUES %s", queryValues)

	db := DatabaseConnection()
	defer db.Close()

	results, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	log.Print("RowsAffected:")
	log.Print(results.RowsAffected())

	return hardSkills, nil
}

func fetchHardSkills(hardSkills []HardSkill) ([]HardSkill, error) {
	if len(hardSkills) == 0 {
		return nil, nil
	}

	valuePlaceholders := make([]string, len(hardSkills))
	for i, hardSkill := range hardSkills {
		valuePlaceholders[i] = fmt.Sprintf("'%s'", hardSkill.Name)
	}
	queryValue := strings.Trim(strings.Join(valuePlaceholders, ","), "")
	query := fmt.Sprintf("SELECT id, name, link, logo, hard_skill_type_id FROM hard_skills WHERE name IN (%s)", queryValue)
	fmt.Println("query", query)

	rows, err := DatabaseConnection().Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var updatedHardSkills []HardSkill
	for rows.Next() {
		var hardSkill HardSkill
		err := rows.Scan(
			&hardSkill.Id,
			&hardSkill.Name,
			&hardSkill.Link,
			&hardSkill.Logo,
			&hardSkill.HardSkillTypeId,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		updatedHardSkills = append(updatedHardSkills, hardSkill)
	}
	return updatedHardSkills, nil
}

func ReplaceAndFetchHardSkills(hardSkills []HardSkill) ([]HardSkill, error) {
	updatedHardSkills, err := replaceHardSkills(hardSkills)
	if err != nil {
		return nil, err
	}
	return fetchHardSkills(updatedHardSkills)
}
