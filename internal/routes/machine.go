package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/sereneblue/lakitu/models"
	"github.com/sereneblue/lakitu/models/awsclient"
	"github.com/sereneblue/lakitu/models/taskrunner"
)

func CreateMachine(c echo.Context) error {
	region := c.FormValue("region")
	size := c.FormValue("size")

	storageSize, _ := strconv.ParseInt(size, 10, 32)

	name := c.FormValue("name")
	streamSW := c.FormValue("streamSW")
	instanceType := c.FormValue("instanceType")

	var streamSoftware awsclient.StreamSoftware

	if streamSW == awsclient.PARSEC.String() {
		streamSoftware = awsclient.PARSEC
	} else {
		streamSoftware = awsclient.MOONLIGHT
	}

	newMachine := models.NewMachine(name, region, streamSoftware, instanceType, int32(storageSize))
	models.Engine.InsertOne(&newMachine)

	jobId, err := runner.Queue(taskrunner.JobCreateMachine, taskrunner.JobData{
		MachineId: newMachine.Id,
	})

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"jobId": jobId,
		},
	})
}

func DeleteMachine(c echo.Context) error {
	machineUuid := c.FormValue("uuid")

	if machineUuid == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Machine UUID not provided",
		})
	}

	machineId := models.GetMachineId(machineUuid)
	jobId, err := runner.Queue(taskrunner.JobDeleteMachine, taskrunner.JobData{
		MachineId: machineId,
	})

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"jobId": jobId,
		},
	})
}

func ListMachines(c echo.Context) error {
	sess, _ := session.Get("session", c)

	machines := models.GetMachines()

	client := awsclient.NewAWSClient(sess.Values["accessKey"].(string), sess.Values["secretKey"].(string), sess.Values["defaultRegion"].(string))
	instances := client.GetImagesAndIntances()

	for _, i := range instances {
		for _, m := range machines {
			fmt.Println(i, m)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string][]models.MachineDetail{
			"machines": machines,
		},
	})
}

func ResizeMachine(c echo.Context) error {
	machineUuid := c.FormValue("uuid")
	newSize := c.FormValue("size")

	if machineUuid == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Machine UUID not provided",
		})
	}

	machineId := models.GetMachineId(machineUuid)
	jobId, err := runner.Queue(taskrunner.JobResizeMachine, taskrunner.JobData{
		MachineId: machineId,
		Metadata:  newSize,
	})

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"jobId": jobId,
		},
	})
}

func StartMachine(c echo.Context) error {
	machineUuid := c.FormValue("uuid")

	if machineUuid == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Machine UUID not provided",
		})
	}

	machineId := models.GetMachineId(machineUuid)
	jobId, err := runner.Queue(taskrunner.JobStartMachine, taskrunner.JobData{
		MachineId: machineId,
	})

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"jobId": jobId,
		},
	})
}

func StopMachine(c echo.Context) error {
	machineUuid := c.FormValue("uuid")

	if machineUuid == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Machine UUID not provided",
		})
	}

	machineId := models.GetMachineId(machineUuid)
	jobId, err := runner.Queue(taskrunner.JobStopMachine, taskrunner.JobData{
		MachineId: machineId,
	})

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"jobId": jobId,
		},
	})
}

func TransferMachine(c echo.Context) error {
	machineUuid := c.FormValue("uuid")
	newRegion := c.FormValue("region")

	if machineUuid == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Machine UUID not provided",
		})
	}

	machineId := models.GetMachineId(machineUuid)
	jobId, err := runner.Queue(taskrunner.JobTransferMachine, taskrunner.JobData{
		MachineId: machineId,
		Metadata:  newRegion,
	})

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"jobId": jobId,
		},
	})
}
