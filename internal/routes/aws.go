package routes

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/sereneblue/lakitu/models/awsclient"
)

type AWSCredentials struct {
	AccessKey string `form:"accessKey"`
	SecretKey string `form:"secretKey"`
}

func GetAWSRegions(c echo.Context) error {
	sess, _ := session.Get("session", c)

	client := awsclient.NewAWSClient(sess.Values["accessKey"].(string), sess.Values["secretKey"].(string), sess.Values["defaultRegion"].(string))
	regions := client.GetRegions()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"regions": regions,
		},
	})
}

func GetAWSGPUInstances(c echo.Context) error {
	region := c.FormValue("region")

	sess, _ := session.Get("session", c)

	client := awsclient.NewAWSClient(sess.Values["accessKey"].(string), sess.Values["secretKey"].(string), sess.Values["defaultRegion"].(string))
	instances := client.GetGPUInstances(region)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"instances": instances,
		},
	})
}

func GetAWSPricing(c echo.Context) error {
	region := c.FormValue("region")

	sess, _ := session.Get("session", c)

	client := awsclient.NewAWSClient(sess.Values["accessKey"].(string), sess.Values["secretKey"].(string), sess.Values["defaultRegion"].(string))
	prices := client.GetPrices(region)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    prices,
	})
}

func VerifiyAWSCredentials(c echo.Context) error {
	creds := new(AWSCredentials)

	if err := c.Bind(creds); err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Invalid input",
		})
	}

	client := awsclient.NewAWSClient(creds.AccessKey, creds.SecretKey, "us-east-1")
	success, err := client.IsValidAWSCredentials()

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": success,
			"message": err.Error(),
		})
	}

	regions := client.GetRegions()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": success,
		"data": map[string]interface{}{
			"regions": regions,
		},
	})
}

func ping(region string, latency map[string]int, wg *sync.WaitGroup, m *sync.Mutex) {
	defer wg.Done()

	addr, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("ec2.%s.amazonaws.com:80", region))

	start := time.Now()
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		latency[region] = 0
	}
	defer conn.Close()

	m.Lock()
	latency[region] = int(time.Since(start) / time.Millisecond)
	m.Unlock()
}

func PingAWS(c echo.Context) error {
	form := new(LatencyForm)

	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Invalid input",
		})
	}

	var latency = make(map[string]int)
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	wg.Add(len(form.Regions))

	for _, region := range form.Regions {
		go ping(region, latency, &wg, mutex)
	}

	wg.Wait()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"latency": latency,
		},
	})
}
