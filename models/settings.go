package models

import "github.com/alexedwards/argon2id"

type Settings struct {
	Key   string `xorm:"not null unique index"`
	Value string
}

var KEKParams = &argon2id.Params{
	Memory:      32 * 1024,
	Iterations:  8,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

func (s *Settings) Insert() (bool, error) {
	_, err := engine.Insert(s)

	if err != nil {
		return false, err
	}

	return true, nil
}

func IsFirstRun() bool {
	has, _ := engine.Table("settings").Exist()
	if has {
		return false
	}

	return true
}