package models

import (
	_ "oestrada1001/lp-chatgpt-integration/database"
)

type Labellable interface {
	GetId() int
	GetLabel() string
	GetValue() string
	GetDescription() string
}
