package models

import (
	_ "database/sql"
)

type HardSkill struct {
	Id              int
	Name            string
	Link            string
	Logo            string
	HardSkillTypeId int
}
