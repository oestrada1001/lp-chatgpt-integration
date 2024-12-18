package models

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
