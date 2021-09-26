package models

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"

	"github.com/sereneblue/lakitu/models/awsclient"
)

var Engine *xorm.Engine
var initError error

func init() {
	userDir, err := os.UserConfigDir()
	if err != nil {
		initError = err
		return
	}

	err = os.MkdirAll(userDir+"/lakitu", os.ModeSticky|os.ModePerm)
	if err != nil {
		initError = err
		return
	}

	Engine, err = xorm.NewEngine("sqlite3", userDir+"/lakitu/lakitu.db")

	if err != nil {
		initError = err
		return
	}

	err = Engine.Sync2(new(Settings))
	if err != nil {
		initError = err
		return
	}

	err = Engine.Sync2(new(awsclient.AWSRole))
	if err != nil {
		initError = err
		return
	}

	err = Engine.Sync2(new(awsclient.AWSSecurityGroup))
	if err != nil {
		initError = err
		return
	}

	// create triggers for job queue
	Engine.Exec(`
		CREATE TRIGGER IF NOT EXISTS update_job_status_success AFTER UPDATE OF status ON task
		  WHEN new.status = 1
		  BEGIN
		    UPDATE job
		    SET status = 1, complete = 1
		    WHERE job.id = new.job_id
		    AND (
			SELECT COUNT(*)
			FROM task
			WHERE job_id = new.job_id
			AND status = 0	    
		    ) = 0;
		  END;

		CREATE TRIGGER IF NOT EXISTS update_job_status_error AFTER UPDATE OF status ON task
		  WHEN new.status = 2
		  BEGIN
		    UPDATE job
		    SET status = 2, complete = 1
		    WHERE job.id = new.job_id;
		  END;
	`)
}

func IsInit() error {
	return initError
}

func CloseDB() {
	Engine.Close()
}
