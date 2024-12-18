package services

import (
	"encoding/json"
	"log"
	"oestrada1001/lp-chatgpt-integration/models"
)

func replaceProficiencyLevels(proficiencyLevels []models.ProficiencyLevel) ([]models.ProficiencyLevel, error) {
	if len(proficiencyLevels) == 0 {
		return nil, nil
	}

	_, err := CreateAndExecuteReplaceQuery("proficiency_levels", proficiencyLevels)
	if err != nil {
		return nil, err
	}

	return proficiencyLevels, nil
}

func fetchProficiencyLevels(proficiencyLevels []models.ProficiencyLevel) ([]models.ProficiencyLevel, error) {
	if len(proficiencyLevels) == 0 {
		return nil, nil
	}

	rows, _ := CreateAndExecuteReadQuery("proficiency_levels", proficiencyLevels)
	defer rows.Close()
	var updatedProficiencyLevels []models.ProficiencyLevel
	for rows.Next() {
		var proficiencyLevel models.ProficiencyLevel
		err := rows.Scan(
			&proficiencyLevel.Id,
			&proficiencyLevel.Label,
			&proficiencyLevel.Value,
			&proficiencyLevel.Description,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		updatedProficiencyLevels = append(updatedProficiencyLevels, proficiencyLevel)
	}
	return updatedProficiencyLevels, nil
}

func ReplaceAndFetchProficiencyLevels(proficiencyLevels []models.ProficiencyLevel) ([]models.ProficiencyLevel, error) {
	updateContextTypes, err := replaceProficiencyLevels(proficiencyLevels)
	if err != nil {
		return nil, err
	}
	return fetchProficiencyLevels(updateContextTypes)
}

func CreateOrGetProficiencyLevels(proficiencyLevels []models.ProficiencyLevel) (string, error) {
	proficiencyLevels, err := ReplaceAndFetchProficiencyLevels(proficiencyLevels)
	if err != nil {
		return "", err
	}
	jsonProficiencyLevels, err := json.Marshal(proficiencyLevels)
	if err != nil {
		return "", err
	}
	return string(jsonProficiencyLevels), nil
}
