package taskrunner

import (
	"sync"
	"sync/atomic"
	"time"

	"xorm.io/xorm"

	"github.com/sereneblue/lakitu/models"
	"github.com/sereneblue/lakitu/models/awsclient"
)

func init() {
	models.Engine.Sync2(new(Job))
	models.Engine.Sync2(new(Task))
}

type Status int64

const (
	PENDING Status = iota
	COMPLETE
	ERROR
	CANCELED
)

type TaskRunner struct {
	client  awsclient.AWSClient
	mu      *sync.Mutex
	running int32
}

func (t *TaskRunner) GetCurrentJob() int64 {
	var j Job

	found, err := models.Engine.Where("complete = 0").Get(&j)

	if err != nil {
		return 0
	}

	if found {
		return j.Id
	}

	return 0
}

func (t *TaskRunner) IsRunning() bool {
	return atomic.LoadInt32(&(t.running)) == 1
}

func (t *TaskRunner) Queue(jobType JobType, jobData JobData) (int64, error) {
	res, err := models.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		job := NewJob(jobType, jobData)
		_, err := models.Engine.InsertOne(&job)
		if err != nil {
			return nil, err
		}

		_, err = models.Engine.Insert(job.GetTasks())
		if err != nil {
			return nil, err
		}

		return job.Id, nil
	})

	if err != nil {
		return 0, err
	}

	return res.(int64), nil
}

func (t *TaskRunner) Start(client awsclient.AWSClient) {
	var nextTask Task

	t.mu.Lock()
	defer t.mu.Unlock()

	atomic.StoreInt32(&(t.running), int32(1))

	t.client = client

	for range time.Tick(time.Second * 5) {
		if atomic.LoadInt32(&(t.running)) == 0 {
			break
		}

		found, err := models.Engine.
			Table("task").
			Select("task.*").
			Join("INNER", "job", "job.id = task.job_id").
			Where("task.status = ? AND job.status = ?", PENDING, PENDING).
			Get(&nextTask)

		if err == nil && found {
			var m models.Machine

			models.Engine.
				Table("machine").
				Select("machine.*").
				Join("INNER", "job", "job.machine_id = machine.id").
				Where("job.id = ?", nextTask.JobId).Get(&m)

			nextTask.HandleTask(t.client, m)
			nextTask = Task{}
		}
	}
}

func (t *TaskRunner) Stop() {
	atomic.StoreInt32(&(t.running), int32(0))
}

func NewTaskRunner() TaskRunner {
	var mu sync.Mutex

	return TaskRunner{
		mu: &mu,
	}
}
