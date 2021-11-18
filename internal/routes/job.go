package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sereneblue/lakitu/models"
	"github.com/sereneblue/lakitu/models/taskrunner"
)

func GetCurrentJobStatus(c echo.Context) error {
	jobId := runner.GetCurrentJob()

	if jobId > 0 {
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/jobs/%d", jobId))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "No jobs in queue",
	})
}

func GetJobStatus(c echo.Context) error {
	id := c.Param("id")

	results, err := models.Engine.Query(`
	    	WITH tasks as (
				SELECT *
			  	  FROM task
			 	 WHERE job_id = ?
	    	)
		    SELECT j.status,
		    	   IFNULL(MIN(CASE WHEN tasks.status = ? THEN tasks.type END), 0) as currentTask,
			   	   SUM(CASE WHEN tasks.status = ? THEN 1 ELSE 0 END) as completed, 
			   	   COUNT(tasks.id) as total,
			   	   MAX(tasks.error_info) as error,
			   	   j.type as jobType
		      FROM job j
   LEFT OUTER JOIN tasks
                   ON j.id = tasks.job_id
		   	 WHERE j.id =?
	`, id, taskrunner.PENDING, taskrunner.COMPLETE, id)

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	if len(results) > 0 {
		jobType, _ := strconv.ParseInt(string(results[0]["jobType"]), 10, 64)
		completed, _ := strconv.ParseInt(string(results[0]["completed"]), 10, 64)
		status, _ := strconv.ParseInt(string(results[0]["status"]), 10, 64)
		currentTask, _ := strconv.ParseInt(string(results[0]["currentTask"]), 10, 64)
		total, _ := strconv.ParseInt(string(results[0]["total"]), 10, 64)
		errorInfo := string(results[0]["error"])

		if completed == total {
			currentTask = taskrunner.TaskComplete
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"status":      status,
				"name":        taskrunner.JOB_NAME[jobType],
				"error": 	   errorInfo,
				"isComplete":  completed == total,
				"completed":   completed,
				"total":       total,
				"currentTask": taskrunner.TASK_NAME[currentTask],
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": false,
		"message": "Job does not exist",
	})
}
