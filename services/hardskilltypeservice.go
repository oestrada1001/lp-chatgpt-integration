package services

import (
	"encoding/json"
	"log"
	"oestrada1001/lp-chatgpt-integration/models"
)

func replaceHardSkillTypes(hardSkillTypes []models.HardSkillType) ([]models.HardSkillType, error) {
	if len(hardSkillTypes) == 0 {
		return nil, nil
	}

	_, err := CreateAndExecuteReplaceQuery("hard_skill_types", hardSkillTypes)
	if err != nil {
		return nil, err
	}

	return hardSkillTypes, nil
}

func fetchHardSkillTypes(hardSkillTypes []models.HardSkillType) ([]models.HardSkillType, error) {
	if len(hardSkillTypes) == 0 {
		return nil, nil
	}

	rows, _ := CreateAndExecuteReadQuery("hard_skill_types", hardSkillTypes)
	defer rows.Close()

	var updatedHardSkillTypes []models.HardSkillType
	for rows.Next() {
		var hardSkillType models.HardSkillType
		err := rows.Scan(
			&hardSkillType.Id,
			&hardSkillType.Label,
			&hardSkillType.Value,
			&hardSkillType.Description,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		updatedHardSkillTypes = append(updatedHardSkillTypes, hardSkillType)
	}
	return updatedHardSkillTypes, nil
}

func ReplaceAndFetchHardSkillTypes(hardSkillTypes []models.HardSkillType) ([]models.HardSkillType, error) {
	updatedHardSkillTypes, err := replaceHardSkillTypes(hardSkillTypes)
	if err != nil {
		return nil, err
	}
	return fetchHardSkillTypes(updatedHardSkillTypes)
}

func CreateOrGetHardSkillTypes(hardSkillTypes []models.HardSkillType) (string, error) {
	hardSkillTypes, err := ReplaceAndFetchHardSkillTypes(hardSkillTypes)
	if err != nil {
		return "", err
	}

	jsonHardSkillTypes, err := json.Marshal(hardSkillTypes)
	if err != nil {
		return "", err
	}
	return string(jsonHardSkillTypes), nil
}
