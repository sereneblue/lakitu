package models

import (
	"github.com/google/uuid"
	"github.com/sereneblue/lakitu/internal/util"
	"github.com/sereneblue/lakitu/models/awsclient"
)

type MachineStatus = string

type Machine struct {
	Id             int64
	Name           string
	Uuid           string
	AmiId          string
	SnapshotId     string
	Region         string
	AdminPassword  string
	StreamSoftware awsclient.StreamSoftware
	InstanceType   string
	Size           int32
	Deleted        bool
}

type MachineDetail struct {
	Status       MachineStatus `json:"status"`
	Uuid         string        `json:"uuid"`
	Name         string        `json:"name"`
	Region       string        `json:"region"`
	InstanceType string        `json:"instanceType"`
	AmiId        string        `json:"amiId"`
	SnapshotId   string        `json:"snapshotId"`
	Size         int32         `json:"size"`
}

const (
	MachineStatusOnline      MachineStatus = "online"
	MachineStatusOffline     MachineStatus = "offline"
	MachineStatusUnavailable MachineStatus = "unavailable"
	MachineStatusUnknown     MachineStatus = "unknown"
)

func GetMachineId(uuid string) int64 {
	var m Machine

	has, err := Engine.Where("uuid = ?", uuid).Get(&m)

	if err != nil || !has {
		return 0
	}

	return m.Id
}

func GetMachines() []MachineDetail {
	var machines []Machine
	machineDetails := []MachineDetail{}

	Engine.Where("deleted = 0").Find(&machines)

	for _, m := range machines {
		machineDetails = append(machineDetails, MachineDetail{
			Status:       MachineStatusUnknown,
			Uuid:         m.Uuid,
			Name:         m.Name,
			Region:       m.Region,
			InstanceType: m.InstanceType,
			Size:         m.Size,
			AmiId:        m.AmiId,
			SnapshotId:   m.SnapshotId,
		})
	}

	return machineDetails
}

func NewMachine(name string, region string, streamSw awsclient.StreamSoftware, instanceType string, size int32) Machine {
	// shouldn't fail, but who knows..
	pwd, _ := util.GeneratePassword()

	m := Machine{
		Name:           name,
		Uuid:           uuid.NewString(),
		Region:         region,
		AdminPassword:  pwd,
		StreamSoftware: streamSw,
		InstanceType:   instanceType,
		Size:           size,
	}

	return m
}
