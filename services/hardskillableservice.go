package services

import (
	"encoding/json"
	"fmt"
	"log"
	"oestrada1001/lp-chatgpt-integration/database"
	"oestrada1001/lp-chatgpt-integration/models"
	"strings"
)

func replaceHardSkillable(hardSkillables []models.HardSkillable) ([]models.HardSkillable, error) {
	if len(hardSkillables) == 0 {
		return nil, nil
	}

	valuePlaceholder := make([]string, len(hardSkillables))
	for i, item := range hardSkillables {
		valuePlaceholder[i] = fmt.Sprintf("('%d','%s','%d','%d','%d')",
			item.SkillableId,
			item.SkillableType,
			item.HardSkillId,
			item.ProficiencyLevelId,
			item.HardSkillContextId,
		)
	}
	queryValues := strings.Trim(strings.Join(valuePlaceholder, ","), "")
	query := fmt.Sprintf("REPLACE INTO hard_skillables (skillable_id, skillable_type, hard_skill_id, proficiency_level_id, hard_skill_context_id) VALUES %s", queryValues)

	db := database.Connection()
	defer db.Close()

	results, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	log.Print("RowsAffected:")
	log.Print(results.RowsAffected())

	return hardSkillables, nil
}

func fetchHardSkillables(hardSkillables []models.HardSkillable) ([]models.HardSkillable, error) {
	if len(hardSkillables) == 0 {
		return nil, nil
	}

	valuePlaceholder := make([]string, len(hardSkillables))
	for i, item := range hardSkillables {
		valuePlaceholder[i] = fmt.Sprintf("'%s'", item.SkillableId)
	}
	queryValue := strings.Trim(strings.Join(valuePlaceholder, ","), "")
	query := fmt.Sprintf("Select * from hard_skillables where skillable_id IN (%s)", queryValue)
	fmt.Println("query", query)

	rows, err := database.Connection().Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var updatedHardSkillables []models.HardSkillable
	for rows.Next() {
		var skillable models.HardSkillable
		err := rows.Scan(
			&skillable.SkillableId,
			&skillable.SkillableType,
			&skillable.HardSkillId,
			&skillable.ProficiencyLevelId,
			&skillable.HardSkillContextId,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		updatedHardSkillables = append(updatedHardSkillables, skillable)
	}
	return updatedHardSkillables, nil
}

func ReplaceAndFetchHardSkillables(skillables []models.HardSkillable) ([]models.HardSkillable, error) {
	updateSkillables, err := replaceHardSkillable(skillables)
	if err != nil {
		return nil, err
	}
	return fetchHardSkillables(updateSkillables)
}

func CreateOrGetHardSkillables(skillables []models.HardSkillable) (string, error) {
	skillables, err := ReplaceAndFetchHardSkillables(skillables)
	if err != nil {
		return "", err
	}
	jsonSkillables, err := json.Marshal(skillables)
	if err != nil {
		return "", err
	}
	return string(jsonSkillables), nil
}
