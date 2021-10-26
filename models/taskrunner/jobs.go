package taskrunner

import (
	"time"
)

type JobType = int64

const (
	JobNone JobType = iota
	JobCreateMachine
	JobDeleteMachine
	JobStartMachine
	JobStopMachine
	JobTransferMachine
	JobResizeMachine
)

var JOB_NAME map[JobType]string = map[JobType]string{
	JobNone:            "Unknown job",
	JobCreateMachine:   "Creating machine",
	JobDeleteMachine:   "Deleting machine",
	JobStartMachine:    "Starting machine",
	JobStopMachine:     "Stopping machine",
	JobTransferMachine: "Transferring machine",
	JobResizeMachine:   "Resizing machine",
}

type Job struct {
	Id        int64
	MachineId int64
	Metadata  string
	Type      JobType
	Status    Status
	Complete  bool
	Created   time.Time `xorm:"created"`
	Updated   time.Time `xorm:"updated"`
}

type JobData struct {
	MachineId int64
	Metadata  string
}

func (j *Job) GetTasks() []Task {
	switch j.Type {
	case JobCreateMachine:
		return []Task{
			NewTask(j.Id, TaskCreateRole),
			NewTask(j.Id, TaskCreateSecurityGroup),
			NewTask(j.Id, TaskCreateInstance),
		}
	case JobDeleteMachine:
		return []Task{
			NewTask(j.Id, TaskDeleteImage),
			NewTask(j.Id, TaskDeleteSnapshot),
		}
	case JobStartMachine:
		return []Task{
			NewTask(j.Id, TaskCreateInstance),
		}
	case JobStopMachine:
		return []Task{
			NewTask(j.Id, TaskSaveInstance),
			NewTask(j.Id, TaskStopInstance),
		}
	case JobTransferMachine:
		return []Task{
			NewTask(j.Id, TaskTransferImage),
			NewTask(j.Id, TaskTransferSnapshot),
		}
	case JobResizeMachine:
		return []Task{
			NewTask(j.Id, TaskResizeVolume),
		}
	default:
		break
	}

	return []Task{}
}

func NewJob(jobType JobType, data JobData) Job {
	return Job{
		MachineId: data.MachineId,
		Type:      jobType,
		Status:    PENDING,
		Complete:  false,
		Metadata:  data.Metadata,
	}
}
