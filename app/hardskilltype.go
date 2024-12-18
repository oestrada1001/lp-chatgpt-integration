package app

import (
	"encoding/json"
	"log"
	"oestrada1001/lp-chatgpt-integration/chatgpt"
)

type HardSkillType struct {
	Id          int
	Label       string
	Value       string
	Description string
}

func (h HardSkillType) GetId() int {
	return h.Id
}

func (h HardSkillType) GetLabel() string {
	return h.Label
}

func (h HardSkillType) GetValue() string {
	return h.Value
}

func (h HardSkillType) GetDescription() string {
	return h.Description
}

func replaceHardSkillTypes(hardSkillTypes []HardSkillType) ([]HardSkillType, error) {
	if len(hardSkillTypes) == 0 {
		return nil, nil
	}

	_, err := CreateAndExecuteReplaceQuery("hard_skill_types", hardSkillTypes)
	if err != nil {
		return nil, err
	}

	return hardSkillTypes, nil
}

func fetchHardSkillTypes(hardSkillTypes []HardSkillType) ([]HardSkillType, error) {
	if len(hardSkillTypes) == 0 {
		return nil, nil
	}

	rows, _ := CreateAndExecuteReadQuery("hard_skill_types", hardSkillTypes)
	defer rows.Close()

	var updatedHardSkillTypes []HardSkillType
	for rows.Next() {
		var hardSkillType HardSkillType
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

func ReplaceAndFetchHardSkillTypes(hardSkillTypes []HardSkillType) ([]HardSkillType, error) {
	updatedHardSkillTypes, err := replaceHardSkillTypes(hardSkillTypes)
	if err != nil {
		return nil, err
	}
	return fetchHardSkillTypes(updatedHardSkillTypes)
}

func CreateOrGetHardSkillTypes(hardSkillTypes []HardSkillType) (chatgpt.FunctionResponse, error) {
	hardSkillTypes, err := ReplaceAndFetchHardSkillTypes(hardSkillTypes)
	if err != nil {
		return chatgpt.FunctionResponse{}, err
	}

	jsonHardSkillTypes, err := json.Marshal(hardSkillTypes)
	if err != nil {
		return chatgpt.FunctionResponse{}, err
	}
	return chatgpt.FunctionResponse{
		Message: "Hard Skill Types",
		Data:    string(jsonHardSkillTypes),
	}, nil
}
