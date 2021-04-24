package models

type Settings struct {
	Key   string `xorm:"not null unique index"`
	Value string
}

func IsFirstRun() bool {
	has, _ := engine.Table("settings").Exist()
	if has {
		return false
	}

	return true
}
