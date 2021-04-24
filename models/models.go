package models

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func InitDB() error {
	userDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	err = os.MkdirAll(userDir+"/lakitu", os.ModeSticky|os.ModePerm)
	if err != nil {
		return err
	}

	engine, err = xorm.NewEngine("sqlite3", userDir+"/lakitu/lakitu.db")
	
	if err != nil {
		return err
	}

	err = engine.Sync2(new(Settings))
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() {
	engine.Close()
}