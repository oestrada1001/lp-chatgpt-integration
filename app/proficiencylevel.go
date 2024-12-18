package app

import "log"

type ProficiencyLevel struct {
	Id          int
	Label       string
	Value       string
	Description string
}

func (h ProficiencyLevel) GetId() int {
	return h.Id
}

func (h ProficiencyLevel) GetLabel() string {
	return h.Label
}

func (h ProficiencyLevel) GetValue() string {
	return h.Value
}

func (h ProficiencyLevel) GetDescription() string {
	return h.Description
}

func replaceProficiencyLevels(proficiencyLevels []ProficiencyLevel) ([]ProficiencyLevel, error) {
	if len(proficiencyLevels) == 0 {
		return nil, nil
	}

	_, err := CreateAndExecuteReplaceQuery("proficiency_levels", proficiencyLevels)
	if err != nil {
		return nil, err
	}

	return proficiencyLevels, nil
}

func fetchProficiencyLevels(proficiencyLevels []ProficiencyLevel) ([]ProficiencyLevel, error) {
	if len(proficiencyLevels) == 0 {
		return nil, nil
	}

	rows, _ := CreateAndExecuteReadQuery("proficiency_levels", proficiencyLevels)
	defer rows.Close()
	var updatedProficiencyLevels []ProficiencyLevel
	for rows.Next() {
		var proficiencyLevel ProficiencyLevel
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

func ReplaceAndFetchProficiencyLevels(proficiencyLevels []ProficiencyLevel) ([]ProficiencyLevel, error) {
	updateContextTypes, err := replaceProficiencyLevels(proficiencyLevels)
	if err != nil {
		return nil, err
	}
	return fetchProficiencyLevels(updateContextTypes)
}
