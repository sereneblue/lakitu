package models

import (
	"github.com/google/uuid"
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
	StreamSoftware awsclient.StreamSoftware
	InstanceType   string
	Size           int64
	Deleted        bool
}

type MachineDetail struct {
	Status       MachineStatus `json:"status"`
	Uuid         string        `json:"uuid"`
	Name         string        `json:"name"`
	Region       string        `json:"region"`
	InstanceType string        `json:"instanceType"`
	Size         int64         `json:"size"`
}

const (
	MachineOnline      MachineStatus = "online"
	MachineOffline     MachineStatus = "offline"
	MachineUnavailable MachineStatus = "unavailable"
)

func GetMachineId(uuid string) int64 {
	var m Machine

	has, err := Engine.Where("uuid = ?", uuid).Get(&m)

	if err != nil || !has {
		return 0
	}

	return m.Id
}

func NewMachine(name string, region string, streamSw awsclient.StreamSoftware, instanceType string, size int64) Machine {
	m := Machine{
		Name:           name,
		Uuid:           uuid.NewString(),
		Region:         region,
		StreamSoftware: streamSw,
		InstanceType:   instanceType,
		Size:           size,
	}

	return m
}
