package app

import (
	"fmt"
	"log"
	"strings"
)

type HardSkillable struct {
	HardSkillId        int
	SkillableId        int
	SkillableType      string
	ProficiencyLevelId int
	HardSkillContextId int
}

func replaceHardSkillable(hardSkillables []HardSkillable) ([]HardSkillable, error) {
	if len(hardSkillables) == 0 {
		return nil, nil
	}

	valuePlaceholder := make([]string, len(hardSkillables))
	for i, item := range hardSkillables {
		valuePlaceholder[i] = fmt.Sprintf("('%s','%s','%s','%s','%s')",
			item.SkillableId,
			item.SkillableType,
			item.HardSkillId,
			item.ProficiencyLevelId,
			item.HardSkillContextId,
		)
	}
	queryValues := strings.ReplaceAll(strings.Trim(strings.Join(valuePlaceholder, ","), ""), " ", "")
	query := fmt.Sprintf("REPLACE INTO hard_skillables (skillable_id, skillable_type, hard_skill_id, proficiency_level_id, hard_skill_context_id) VALUES %s", queryValues)

	db := DatabaseConnection()
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

func fetchHardSkillables(hardSkillables []HardSkillable) ([]HardSkillable, error) {
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

	rows, err := DatabaseConnection().Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var updatedHardSkillables []HardSkillable
	for rows.Next() {
		var skillable HardSkillable
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

func ReplaceAndFetchHardSkillables(skillables []HardSkillable) ([]HardSkillable, error) {
	updateSkillables, err := replaceHardSkillable(skillables)
	if err != nil {
		return nil, err
	}
	return fetchHardSkillables(updateSkillables)
}
