package services

import (
	"log"
	"oestrada1001/lp-chatgpt-integration/models"
)

func replaceHardSkillContexts(hardSkillContexts []models.HardSkillContext) ([]models.HardSkillContext, error) {
	if len(hardSkillContexts) == 0 {
		return nil, nil
	}

	_, err := CreateAndExecuteReplaceQuery("hard_skill_contexts", hardSkillContexts)
	if err != nil {
		return nil, err
	}

	return hardSkillContexts, nil
}

func fetchHardSkillContexts(hardSkillContexts []models.HardSkillContext) ([]models.HardSkillContext, error) {
	if len(hardSkillContexts) == 0 {
		return nil, nil
	}

	rows, _ := CreateAndExecuteReadQuery("hard_skill_contexts", hardSkillContexts)
	defer rows.Close()
	var updatedHardSkillContexts []models.HardSkillContext
	for rows.Next() {
		var hardSkillContext models.HardSkillContext
		err := rows.Scan(
			&hardSkillContext.Id,
			&hardSkillContext.Label,
			&hardSkillContext.Value,
			&hardSkillContext.Description,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		updatedHardSkillContexts = append(updatedHardSkillContexts, hardSkillContext)
	}
	return updatedHardSkillContexts, nil
}

func ReplaceAndFetchHardSkillContexts(hardSkillContexts []models.HardSkillContext) ([]models.HardSkillContext, error) {
	updatedHardSkillContexts, err := replaceHardSkillContexts(hardSkillContexts)
	if err != nil {
		return nil, err
	}
	return fetchHardSkillContexts(updatedHardSkillContexts)
}
