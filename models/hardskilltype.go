package models

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
