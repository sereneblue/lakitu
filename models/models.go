package models

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"

	"github.com/sereneblue/lakitu/models/awsclient"
)

var Engine *xorm.Engine
var initError error

const (
	PENDING = iota
	COMPLETE
	ERROR
	CANCELED
)

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

	err = Engine.Sync2(new(Machine))
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
	Engine.Exec(fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS update_job_status_success AFTER UPDATE OF status ON task
		  WHEN new.status = %d
		  BEGIN
		    UPDATE job
		       SET status = %d, 
		       	   complete = 1
		     WHERE job.id = new.job_id
		           AND (
						SELECT COUNT(*)
						FROM task
						WHERE job_id = new.job_id
						AND status = %d    
				   ) = 0;
		  END;

		CREATE TRIGGER IF NOT EXISTS update_job_status_error AFTER UPDATE OF status ON task
		  WHEN new.status = %d
		  BEGIN
		    UPDATE job
		       SET status = %d,
		           complete = 1
		     WHERE job.id = new.job_id;

		    UPDATE task
		       SET status = %d
		     WHERE task.job_id = new.job_id
		     	   AND task.status = %d;
		  END;
	`, 
	COMPLETE, COMPLETE, COMPLETE, 
	ERROR, ERROR, CANCELED, PENDING))
}

func IsInit() error {
	return initError
}

func CloseDB() {
	Engine.Close()
}
