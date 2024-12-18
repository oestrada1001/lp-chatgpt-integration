package app

import "log"

type HardSkillContext struct {
	Id          int
	Label       string
	Value       string
	Description string
}

func (h HardSkillContext) GetId() int {
	return h.Id
}

func (h HardSkillContext) GetLabel() string {
	return h.Label
}

func (h HardSkillContext) GetValue() string {
	return h.Value
}

func (h HardSkillContext) GetDescription() string {
	return h.Description
}

func replaceHardSkillContexts(hardSkillContexts []HardSkillContext) ([]HardSkillContext, error) {
	if len(hardSkillContexts) == 0 {
		return nil, nil
	}

	_, err := CreateAndExecuteReplaceQuery("hard_skill_contexts", hardSkillContexts)
	if err != nil {
		return nil, err
	}

	return hardSkillContexts, nil
}

func fetchHardSkillContexts(hardSkillContexts []HardSkillContext) ([]HardSkillContext, error) {
	if len(hardSkillContexts) == 0 {
		return nil, nil
	}

	rows, _ := CreateAndExecuteReadQuery("hard_skill_contexts", hardSkillContexts)
	defer rows.Close()
	var updatedHardSkillContexts []HardSkillContext
	for rows.Next() {
		var hardSkillContext HardSkillContext
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

func ReplaceAndFetchHardSkillContexts(hardSkillContexts []HardSkillContext) ([]HardSkillContext, error) {
	updatedHardSkillContexts, err := replaceHardSkillContexts(hardSkillContexts)
	if err != nil {
		return nil, err
	}
	return fetchHardSkillContexts(updatedHardSkillContexts)
}
