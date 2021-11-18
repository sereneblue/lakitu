package routes

import (
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/sereneblue/lakitu/internal/util"
	"github.com/sereneblue/lakitu/models"
	"github.com/sereneblue/lakitu/models/awsclient"
	"github.com/sereneblue/lakitu/models/taskrunner"
)

var runner taskrunner.TaskRunner

func init() {
	runner = taskrunner.NewTaskRunner()
}

func ChangePassword(c echo.Context) error {
	oldPwd := c.FormValue("oldPwd")
	newPwd := c.FormValue("newPwd")
	confirmNewPwd := c.FormValue("confirmNewPwd")

	if newPwd != confirmNewPwd {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "New password does not match",
		})
	}

	match, err := argon2id.ComparePasswordAndHash(oldPwd, models.GetPasswordHash())
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "There was an error changing your password",
		})
	}

	if !match {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Invalid password",
		})
	}

	encKey, salt := models.GetEncryptedData()

	KEK := util.CreateHashWithSalt(oldPwd, []byte(salt), models.KEKParams)

	key, err := util.Decrypt(encKey, KEK)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not decrypt encryption key",
		})
	}

	var s models.Settings

	pwdHash, err := argon2id.CreateHash(newPwd, argon2id.DefaultParams)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not hash password",
		})
	}
	s.Key = "password"
	s.Value = pwdHash
	s.Update()

	newKEK, err := argon2id.CreateHash(newPwd, models.KEKParams)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not generate new KEK",
		})
	}

	_, newSalt, _, _ := argon2id.DecodeHash(newKEK)
	s.Key = "encSalt"
	s.Value = string(newSalt)
	s.Update()

	// save encrypted key
	newEncKey, err := util.Encrypt(key, newKEK)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not generate new encrypted key",
		})
	}

	s.Key = "encKey"
	s.Value = newEncKey
	s.Update()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password was successfully updated",
	})
}

func ChangePreferences(c echo.Context) error {
	sess, _ := session.Get("session", c)

	accessKey := c.FormValue("accessKey")
	secretKey := c.FormValue("secretKey")
	defaultRegion := c.FormValue("defaultRegion")

	if defaultRegion != sess.Values["defaultRegion"].(string) {
		if _, ok := awsclient.AWS_REGIONS[defaultRegion]; !ok {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "Invalid AWS region",
			})
		}

		sess.Values["defaultRegion"] = defaultRegion
	}

	// check if aws keys changed
	if accessKey != sess.Values["accessKey"].(string) || secretKey != sess.Values["secretKey"].(string) {
		client := awsclient.NewAWSClient(accessKey, secretKey, defaultRegion)
		success, err := client.IsValidAWSCredentials()

		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": success,
				"message": err.Error(),
			})
		}

		encKey, _ := models.GetEncryptedData()

		key, err := util.Decrypt(encKey, sess.Values["KEK"].(string))
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "Could not decrypt encryption key",
			})
		}

		var s models.Settings

		encAWSAccessKey, err := util.Encrypt(accessKey, key)
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "Could not encrypt AWS access key",
			})
		}
		s.Key = "awsAccessKeyId"
		s.Value = encAWSAccessKey
		s.Update()

		encAWSSecretKey, err := util.Encrypt(secretKey, key)
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "Could not encrypt AWS secret key",
			})
		}
		s.Key = "awsSecretKey"
		s.Value = encAWSSecretKey
		s.Update()

		sess.Values["accessKey"] = accessKey
		sess.Values["secretKey"] = secretKey
	}

	sess.Save(c.Request(), c.Response())

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Preferences were successfully updated",
	})
}

func IsLoggedIn(c echo.Context) error {
	if !runner.IsRunning() {
		sess, _ := session.Get("session", c)

		client := awsclient.NewAWSClient(sess.Values["accessKey"].(string), sess.Values["secretKey"].(string), sess.Values["defaultRegion"].(string))

		go runner.Start(client)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func Login(c echo.Context) error {
	loginPwd := c.FormValue("password")

	match, err := argon2id.ComparePasswordAndHash(loginPwd, models.GetPasswordHash())
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "There was an error during login",
		})
	}

	if !match {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Invalid password",
		})
	}

	encKey, salt := models.GetEncryptedData()
	accessKey, secretKey := models.GetAWSSettings()

	// generate KEK
	KEK := util.CreateHashWithSalt(loginPwd, []byte(salt), models.KEKParams)

	key, err := util.Decrypt(encKey, KEK)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not decrypt encryption key",
		})
	}

	awsAccessKey, err := util.Decrypt(accessKey, key)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not decrypt AWS access key",
		})
	}

	awsSecretKey, err := util.Decrypt(secretKey, key)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not decrypt AWS secret key",
		})
	}

	sess, _ := session.Get("session", c)
	sess.Values["authenticated"] = true
	sess.Values["KEK"] = KEK
	sess.Values["accessKey"] = awsAccessKey
	sess.Values["secretKey"] = awsSecretKey
	sess.Values["defaultRegion"] = models.GetDefaultRegion()

	sess.Save(c.Request(), c.Response())

	pendingJobId := runner.GetCurrentJob()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"hasPending": pendingJobId > 0,
			"jobId":      pendingJobId,
		},
	})
}

func Logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	runner.Stop()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func UserData(c echo.Context) error {
	sess, _ := session.Get("session", c)

	client := awsclient.NewAWSClient(sess.Values["accessKey"].(string), sess.Values["secretKey"].(string), sess.Values["defaultRegion"].(string))
	regions := client.GetRegions()
	log := taskrunner.GetTaskLog()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"accessKey":     sess.Values["accessKey"],
			"secretKey":     sess.Values["secretKey"],
			"defaultRegion": sess.Values["defaultRegion"],
			"regions":       regions,
			"log":           log,
		},
	})
}
