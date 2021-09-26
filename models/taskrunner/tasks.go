package taskrunner

import (
	"strconv"
	"time"

	"github.com/sereneblue/lakitu/models"
	"github.com/sereneblue/lakitu/models/awsclient"
)

type TaskType = int64

const (
	TaskNone JobType = iota
	TaskComplete
	TaskCreateRole
	TaskCreateSecurityGroup
	TaskCreateInstance
	TaskDeleteImage
	TaskDeleteSnapshot
	TaskResizeVolume
	TaskTransferZone
	TaskStartInstance
	TaskStopInstance
)

var TASK_NAME map[TaskType]string = map[TaskType]string{
	TaskNone:                "Unknown task",
	TaskComplete:            "Task complete",
	TaskCreateRole:          "Creating role",
	TaskCreateSecurityGroup: "Creating security group",
	TaskCreateInstance:      "Creating instance",
	TaskDeleteImage:         "Deleting image",
	TaskDeleteSnapshot:      "Deleting snapshot",
	TaskResizeVolume:        "Resizing volume",
	TaskTransferZone:        "Transferring zone",
	TaskStartInstance:       "Starting instance",
	TaskStopInstance:        "Stopping instance",
}

type Task struct {
	Id        int64
	JobId     int64
	Type      TaskType
	ErrorInfo string `xorm:"error_info"`
	Status    Status
	Created   time.Time `xorm:"created"`
	Updated   time.Time `xorm:"updated"`
}

type TaskLog struct {
	Id        int64  `json:"id"`
	Event     string `json:"event"`
	ErrorInfo string `json:"errorInfo"`
	Timestamp int64  `json:"timestamp"`
	Status    Status `json:"status"`
}

func (t *Task) HandleTask(client awsclient.AWSClient, machine models.Machine) {
	switch t.Type {
	}
}

func (t *Task) updateStatus(status Status, errorInfo string) {
	models.Engine.ID(t.Id).Cols("status,error").Update(Task{
		Status:    status,
		ErrorInfo: errorInfo,
	})
}

func GetTaskLog() []TaskLog {
	log := []TaskLog{}

	entries, err := models.Engine.Query(`
		select task.id,
			   task.error_info,
			   task.type,
			   task.status,
			   task.updated
		  from task
	  order by id desc
		 limit 10
	`)

	if err != nil {
		return log
	}

	for _, logEntry := range entries {
		id, _ := strconv.ParseInt(string(logEntry["id"]), 10, 64)
		status, _ := strconv.ParseInt(string(logEntry["status"]), 10, 64)
		taskType, _ := strconv.ParseInt(string(logEntry["type"]), 10, 64)

		updatedTime, _ := time.Parse("2006-01-02T15:04:05Z", string(logEntry["updated"]))

		log = append(log, TaskLog{
			Id:        id,
			Event:     TASK_NAME[taskType],
			ErrorInfo: string(logEntry["error_info"]),
			Timestamp: updatedTime.UnixMilli(),
			Status:    Status(status),
		})
	}

	return log
}

func NewTask(jobId int64, taskType JobType) Task {
	return Task{
		JobId:  jobId,
		Type:   taskType,
		Status: PENDING,
	}
}
