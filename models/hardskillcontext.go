package models

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
