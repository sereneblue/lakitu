package taskrunner

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
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
	TaskCreateNewVolume
	TaskDeleteImage
	TaskDeleteSnapshot
	TaskRequestSpotInstance
	TaskResizeVolume
	TaskTransferImage
	TaskTransferSnapshot
	TaskSaveInstance
	TaskStartInstance
	TaskStopInstance
)

var TASK_NAME map[TaskType]string = map[TaskType]string{
	TaskNone:                "Unknown task",
	TaskComplete:            "Task complete",
	TaskCreateRole:          "Creating role",
	TaskCreateSecurityGroup: "Creating security group",
	TaskCreateInstance:      "Creating instance",
	TaskCreateNewVolume:     "Creating new volume",
	TaskDeleteImage:         "Deleting image",
	TaskDeleteSnapshot:      "Deleting snapshot",
	TaskResizeVolume:        "Resizing volume",
	TaskTransferImage:       "Transferring image",
	TaskTransferSnapshot:    "Transferring snapshot",
	TaskRequestSpotInstance: "Requesting spot instance",
	TaskSaveInstance:        "Saving instance",
	TaskStartInstance:       "Starting instance",
	TaskStopInstance:        "Stopping instance",
}

type Task struct {
	Id        int64
	JobId     int64
	Type      TaskType
	ErrorInfo string `xorm:"error_info"`
	Status    Status
	Metadata  string
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

type TaskSaveInstanceMetadata struct {
	InstanceId    string
	VolumeId      string
	NewAmiId      string
	NewSnapshotId string
}

func (t *Task) HandleTask(client awsclient.AWSClient, machine models.Machine) {
	switch t.Type {
	case TaskCreateRole:
		t.createRole(client, machine)
	case TaskCreateSecurityGroup:
		t.createSecurityGroup(client, machine)
	case TaskCreateInstance:
		t.createInstance(client, machine)
	case TaskCreateNewVolume:
		t.createNewVolume(client, machine)
	case TaskDeleteImage:
		t.deleteImage(client, machine)
	case TaskDeleteSnapshot:
		t.deleteSnapshot(client, machine)
	case TaskResizeVolume:
		t.resizeVolume(client, machine)
	case TaskSaveInstance:
		t.saveMachine(client, machine)
	case TaskRequestSpotInstance:
		t.requestSpotInstance(client, machine)
	case TaskStartInstance:
		t.startMachine(client, machine)
	case TaskStopInstance:
		t.stopMachine(client, machine)
	case TaskTransferImage:
		t.transferImage(client, machine)
	case TaskTransferSnapshot:
		t.transferSnapshot(client, machine)
	default:
		break
	}
}

func (t *Task) updateStatus(status Status, errorInfo string) {
	models.Engine.ID(t.Id).Cols("status,error_info").Update(Task{
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

func (t *Task) createInstance(client awsclient.AWSClient, m models.Machine) {
	var instanceType types.InstanceType
	var j Job

	models.Engine.ID(t.JobId).Get(&j)

	securityGroupId := models.GetSecurityGroupId(m.StreamSoftware)
	role := models.GetRole()

	if m.InstanceType == string(types.InstanceTypeG3sXlarge) {
		instanceType = types.InstanceTypeG3sXlarge
	} else {
		instanceType = types.InstanceTypeG4dnXlarge
	}

	// get ami id
	amiId, err := client.GetWindowsAMIId()
	if err != nil {
		// delete new machine if task fails
		models.Engine.ID(m.Id).Cols("deleted").Update(models.Machine{
			Deleted: true,
		})

		t.updateStatus(ERROR, err.Error())
		return
	}

	models.Engine.ID(m.Id).Cols("ami_id").Update(models.Machine{
		AmiId: amiId,
	})
	m.AmiId = amiId

	instanceId, err := client.CreateInstance(m.AmiId, instanceType, securityGroupId, m.Region, m.AdminPassword, role.Arn)
	if err != nil {
		// delete new machine if task fails
		models.Engine.ID(m.Id).Cols("deleted").Update(models.Machine{
			Deleted: true,
		})

		t.updateStatus(ERROR, err.Error())
		return
	}

	models.Engine.ID(j.Id).Cols("metadata").Update(Job{
		Metadata: instanceId,
	})
	t.updateStatus(COMPLETE, "")
}

func (t *Task) createNewVolume(client awsclient.AWSClient, m models.Machine) {
	var j Job

	models.Engine.ID(t.JobId).Get(&j)

	err := client.CreateNewVolume(j.Metadata, m.Size, m.Region)

	if err == nil {
		t.updateStatus(COMPLETE, "")
	} else {
		// delete new machine if task fails
		models.Engine.ID(m.Id).Cols("deleted").Update(models.Machine{
			Deleted: true,
		})

		t.updateStatus(ERROR, err.Error())
	}
}

func (t *Task) createRole(client awsclient.AWSClient, m models.Machine) {
	role := models.GetRole()
	roles, err := client.GetRoles()

	if err == nil {
		found := false

		for _, r := range roles {
			if *r.RoleId == role.RoleId {
				found = true
			}
		}

		if !found {
			role, err := client.CreateRole()

			if err == nil {
				models.Engine.InsertOne(&role)
				t.updateStatus(COMPLETE, "")
				return
			}
		} else {
			t.updateStatus(COMPLETE, "")
			return
		}
	}

	// delete new machine if task fails
	models.Engine.ID(m.Id).Cols("deleted").Update(models.Machine{
		Deleted: true,
	})

	t.updateStatus(ERROR, err.Error())
}

func (t *Task) createSecurityGroup(client awsclient.AWSClient, m models.Machine) {
	securiyGroupId := models.GetSecurityGroupId(m.StreamSoftware)
	securityGroups, err := client.GetSecurityGroups()

	if err == nil {
		found := false

		for _, sg := range securityGroups {
			if *sg.GroupId == securiyGroupId {
				found = true
			}
		}

		if !found {
			group, err := client.CreateSecurityGroup(m.StreamSoftware)

			if err == nil {
				models.Engine.InsertOne(&group)
				t.updateStatus(COMPLETE, "")
				return
			}
		} else {
			t.updateStatus(COMPLETE, "")
			return
		}
	}

	// delete new machine if task fails
	models.Engine.ID(m.Id).Cols("deleted").Update(models.Machine{
		Deleted: true,
	})

	t.updateStatus(ERROR, err.Error())
}

func (t *Task) deleteImage(client awsclient.AWSClient, m models.Machine) {
	err := client.DeleteImage(m.AmiId, m.Region)

	if err == nil {
		t.updateStatus(COMPLETE, "")
	} else {
		t.updateStatus(ERROR, err.Error())
	}
}

func (t *Task) deleteSnapshot(client awsclient.AWSClient, m models.Machine) {
	err := client.DeleteSnapshot(m.SnapshotId, m.Region)

	if err == nil {
		t.updateStatus(COMPLETE, "")
	} else {
		t.updateStatus(ERROR, err.Error())
	}

	// delete machine
	models.Engine.ID(m.Id).Cols("deleted").Update(models.Machine{
		Deleted: true,
	})
}

func (t *Task) resizeVolume(client awsclient.AWSClient, m models.Machine) {
	var j Job

	models.Engine.ID(t.JobId).Get(&j)

	newSize, _ := strconv.ParseInt(j.Metadata, 10, 32)
	newSnapshotId, err := client.ResizeSnapshot(m.SnapshotId, int32(newSize), m.Region)

	if err == nil {
		t.updateStatus(COMPLETE, "")

		models.Engine.ID(m.Id).Cols("snapshot_id,size").Update(models.Machine{
			SnapshotId: newSnapshotId,
			Size:       int32(newSize),
		})
	} else {
		t.updateStatus(ERROR, err.Error())
	}
}

func (t *Task) requestSpotInstance(client awsclient.AWSClient, m models.Machine) {
	var j Job

	models.Engine.ID(t.JobId).Get(&j)

	var instanceType types.InstanceType

	securityGroupId := models.GetSecurityGroupId(m.StreamSoftware)
	role := models.GetRole()

	if m.InstanceType == string(types.InstanceTypeG3sXlarge) {
		instanceType = types.InstanceTypeG3sXlarge
	} else {
		instanceType = types.InstanceTypeG4dnXlarge
	}

	instanceId, err := client.StartInstance(m.AmiId, m.SnapshotId, instanceType, securityGroupId, m.Region, m.AdminPassword, role.Arn, role.RoleName)

	if err != nil {
		t.updateStatus(ERROR, err.Error())
		return
	}

	models.Engine.ID(j.Id).Cols("metadata").Update(Job{
		Metadata: instanceId,
	})
	t.updateStatus(COMPLETE, "")
}

func (t *Task) saveMachine(client awsclient.AWSClient, m models.Machine) {
	var j Job

	models.Engine.ID(t.JobId).Get(&j)

	// get current instance
	instanceId, volumeId, err := client.GetMachineData(m.AmiId, m.SnapshotId, m.Region)
	if err != nil {
		t.updateStatus(ERROR, err.Error())
		return
	}

	// create ami
	newAmiId, err := client.CreateImage(instanceId, m.Region)
	if err != nil {
		t.updateStatus(ERROR, err.Error())
		return
	}

	// create snapshot
	newSnapshotId, err := client.CreateSnapshot(volumeId, m.Region)
	if err != nil {
		t.updateStatus(ERROR, err.Error())
		return
	}

	jsonMetadata, _ := json.Marshal(TaskSaveInstanceMetadata{
		InstanceId:    instanceId,
		VolumeId:      volumeId,
		NewAmiId:      newAmiId,
		NewSnapshotId: newSnapshotId,
	})

	models.Engine.ID(j.Id).Cols("metadata").Update(Job{
		Metadata: string(jsonMetadata),
	})
	t.updateStatus(COMPLETE, "")
}

func (t *Task) startMachine(client awsclient.AWSClient, m models.Machine) {
	var j Job

	models.Engine.ID(t.JobId).Get(&j)

	// wait for machine to be ready
	for {
		instanceState, err := client.GetInstanceState(j.Metadata, m.Region)
		if err != nil {
			t.updateStatus(ERROR, err.Error())
			return
		}

		if instanceState == types.InstanceStateNameRunning {
			t.updateStatus(COMPLETE, "")
			return
		} else if instanceState == types.InstanceStateNamePending {
			time.Sleep(30 * time.Second)
		} else {
			t.updateStatus(ERROR, "Invalid state for instance: "+j.Metadata)
			return
		}
	}

}

func (t *Task) stopMachine(client awsclient.AWSClient, m models.Machine) {
	var j Job
	var metadata TaskSaveInstanceMetadata

	models.Engine.ID(t.JobId).Get(&j)

	err := json.Unmarshal([]byte(j.Metadata), &metadata)
	if err != nil {
		t.updateStatus(ERROR, err.Error())
		return
	}

	// terminate instance
	err = client.TerminateInstance(metadata.InstanceId, m.Region)
	if err != nil {
		t.updateStatus(ERROR, err.Error())
		return
	}

	// delete volume
	_, err = client.DeleteVolume(metadata.VolumeId, m.Region)
	if err != nil {
		t.updateStatus(ERROR, err.Error())
		return
	}

	// delete older snapshot and ami
	err = client.DeleteImage(m.AmiId, m.Region)
	if err != nil {
		t.updateStatus(ERROR, err.Error())
		return
	}

	err = client.DeleteSnapshot(m.SnapshotId, m.Region)
	if err != nil {
		t.updateStatus(ERROR, err.Error())
		return
	}

	t.updateStatus(COMPLETE, "")
	models.Engine.ID(m.Id).Cols("ami_id,snapshot_id").Update(models.Machine{
		AmiId:      metadata.NewAmiId,
		SnapshotId: metadata.NewSnapshotId,
	})
}

func (t *Task) transferImage(client awsclient.AWSClient, m models.Machine) {
	var j Job

	models.Engine.ID(t.JobId).Get(&j)

	imageId, err := client.CopyImage(m.AmiId, m.Uuid, m.Region, j.Metadata)

	if imageId != "" {
		// update task metadata with new image id
		models.Engine.ID(t.Id).Cols("metadata").Update(Task{
			Metadata: imageId,
		})

		t.updateStatus(COMPLETE, "")
	} else {
		t.updateStatus(ERROR, err.Error())
	}
}

func (t *Task) transferSnapshot(client awsclient.AWSClient, m models.Machine) {
	var j Job
	var transferImageTask Task

	models.Engine.ID(t.JobId).Get(&j)
	models.Engine.Where("task.job_id = ? AND task.type = ?", j.Id, TaskTransferImage).Get(&transferImageTask)

	snapshotId, err := client.CopySnapshot(m.SnapshotId, m.Region, j.Metadata)

	if snapshotId != "" {
		// update machine id
		models.Engine.ID(m.Id).Cols("ami_id,snapshot_id,region").Update(models.Machine{
			AmiId:      transferImageTask.Metadata,
			SnapshotId: snapshotId,
			Region:     j.Metadata,
		})

		err := client.DeleteImage(m.AmiId, m.Region)

		if err == nil {
			err = client.DeleteSnapshot(m.SnapshotId, m.Region)

			if err == nil {
				t.updateStatus(COMPLETE, "")
			} else {
				t.updateStatus(ERROR, err.Error())
			}
		} else {
			t.updateStatus(ERROR, err.Error())
		}
	} else {
		t.updateStatus(ERROR, err.Error())
	}
}
