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

func GetAWSSettings() (string, string) {
	var accessKey, secretKey Settings

	has, err := engine.Where("key = 'awsAccessKeyId'").Get(&accessKey)
	if err != nil || !has {
		return "", ""
	}

	has, err = engine.Where("key = 'awsSecretKey'").Get(&secretKey)
	if err != nil || !has {
		return "", ""
	}

	return accessKey.Value, secretKey.Value
}

func GetEncryptedData() (string, string) {
	var key, salt Settings

	has, err := engine.Where("key = 'encKey'").Get(&key)
	if err != nil || !has {
		return "", ""
	}

	has, err = engine.Where("key = 'encSalt'").Get(&salt)
	if err != nil || !has {
		return "", ""
	}

	return key.Value, salt.Value
}

func GetPasswordHash() string {
	var s Settings

	has, err := engine.Where("key = 'password'").Get(&s)
	if err != nil || !has {
		return ""
	}

	return s.Value
}

func IsFirstRun() bool {
	has, _ := engine.Table("settings").Exist()
	if has {
		return false
	}

	return true
}