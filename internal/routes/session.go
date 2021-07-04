package routes

import (
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/sereneblue/lakitu/internal/util"
	"github.com/sereneblue/lakitu/models"
)

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
	sess.Values["accessKey"] = awsAccessKey
	sess.Values["secretKey"] = awsSecretKey
	sess.Values["defaultRegion"] = models.GetDefaultRegion()

	sess.Save(c.Request(), c.Response())

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func Logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
