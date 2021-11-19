package models

import (
	"github.com/alexedwards/argon2id"
	"github.com/sereneblue/lakitu/models/awsclient"
)

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
	_, err := Engine.Insert(s)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Settings) Update() (bool, error) {
	_, err := Engine.Update(s, &Settings{Key: s.Key})

	if err != nil {
		return false, err
	}

	return true, nil
}

func GetAWSSettings() (string, string) {
	var accessKey, secretKey Settings

	has, err := Engine.Where("key = 'awsAccessKeyId'").Get(&accessKey)
	if err != nil || !has {
		return "", ""
	}

	has, err = Engine.Where("key = 'awsSecretKey'").Get(&secretKey)
	if err != nil || !has {
		return "", ""
	}

	return accessKey.Value, secretKey.Value
}

func GetDefaultRegion() string {
	var s Settings

	has, err := Engine.Where("key = 'defaultRegion'").Get(&s)
	if err != nil || !has {
		return ""
	}

	return s.Value
}

func GetEncryptedData() (string, string) {
	var key, salt Settings

	has, err := Engine.Where("key = 'encKey'").Get(&key)
	if err != nil || !has {
		return "", ""
	}

	has, err = Engine.Where("key = 'encSalt'").Get(&salt)
	if err != nil || !has {
		return "", ""
	}

	return key.Value, salt.Value
}

func GetPasswordHash() string {
	var s Settings

	has, err := Engine.Where("key = 'password'").Get(&s)
	if err != nil || !has {
		return ""
	}

	return s.Value
}

func GetRole() awsclient.AWSRole {
	var r awsclient.AWSRole

	has, err := Engine.Desc("id").Get(&r)

	if err != nil || !has {
		return r
	}

	return r
}

func GetSecurityGroupId(streamSW awsclient.StreamSoftware, region string) string {
	var sg awsclient.AWSSecurityGroup

	has, err := Engine.Desc("id").Where("stream_software = ? AND region = ?", streamSW, region).Get(&sg)

	if err != nil || !has {
		return ""
	}

	return sg.GroupId
}

func IsFirstRun() bool {
	has, _ := Engine.Table("settings").Where("key = 'awsSecretKey'").Exist()
	if has {
		return false
	}

	return true
}
