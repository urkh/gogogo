package controllers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"testapp/src/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetUsers(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	page, _ := strconv.Atoi(c.QueryParam("page"))
	offset := limit * (page - 1)

	var users []models.User

	if res := db.Limit(limit).Offset(offset).Find(&users); res.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to fetch users"})
	}

	response := map[string]interface{}{
		"data": users,
	}

	return c.JSON(http.StatusOK, response)
}

func GetUser(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	id := c.Param("id")

	var user = models.User{ID: id}
	if res := db.First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to fetch user"})
	}

	response := map[string]interface{}{
		"data": user,
	}

	return c.JSON(http.StatusOK, response)
}

/*
	func NewUser(c echo.Context) error {
		db := c.Get("db").(*gorm.DB)
		values := new(models.User)

		if err := c.Bind(values); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		user := &models.User{
			Name:     values.Name,
			Email:    values.Email,
			Password: values.Password,
		}

		if res := db.Create(&user); res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res)
		}

		response := map[string]interface{}{
			"data": user,
		}

		return c.JSON(http.StatusCreated, response)
	}
*/
func NewUser(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	// Parse request body into a new User instance
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request payload"})
	}

	// Create the user in the database
	if res := db.Create(&user); res.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to create user"})
	}

	// Prepare and send the response
	response := map[string]interface{}{
		"data": user,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
func UpdateUser(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	id := c.Param("id")
	values := new(models.User)

	if err := c.Bind(values); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user := models.User{ID: id}

	if res := db.First(&user); res.Error != nil {
		msg := map[string]interface{}{
			"message": res.Error.Error(),
		}
		return c.JSON(http.StatusNotFound, msg)
	}

	user.Name = values.Name
	user.Email = values.Email
	user.Password = values.Password

	if res := db.Save(&user); res.Error != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	response := map[string]interface{}{
		"data": user,
	}

	return c.JSON(http.StatusOK, response)
}
*/

func UpdateUser(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	id := c.Param("id")

	// Parse request body into a new User instance
	var updatedUser models.User
	if err := c.Bind(&updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request payload"})
	}

	// Fetch the existing user from the database
	// var existingUser models.User
	user := models.User{ID: id}
	if res := db.First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to fetch user"})
	}

	// Update the user fields
	user.Name = updatedUser.Name
	user.Email = updatedUser.Email
	user.Password = updatedUser.Password

	// Save the changes in the database
	if res := db.Save(&user); res.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to update user"})
	}

	// Prepare and send the response
	response := map[string]interface{}{
		"data": user,
	}

	return c.JSON(http.StatusOK, response)
}

/*
func DeleteUser(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	id := c.Param("id")

	user := models.User{ID: id}

	if res := db.First(&user); res.Error != nil {
		msg := map[string]interface{}{
			"message": res.Error.Error(),
		}
		return c.JSON(http.StatusNotFound, msg)
	}

	if res := db.Delete(&user); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"message": "a user has been deleted",
	}
	return c.JSON(http.StatusOK, response)
}
*/

func DeleteUser(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	id := c.Param("id")

	// Fetch the user from the database
	// var user models.User
	user := models.User{ID: id}
	if res := db.First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to fetch user"})
	}

	// Delete the user from the database
	if res := db.Delete(&user); res.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to delete user"})
	}

	// Prepare and send the response
	response := map[string]interface{}{
		"message": "User has been deleted",
	}

	return c.JSON(http.StatusOK, response)
}

func Save(c echo.Context) error {
	var name string = c.FormValue("name")
	var email string = c.FormValue("email")

	// get image
	picture, err := c.FormFile("picture")
	if err != nil {
		return err
	}

	// source
	src, err := picture.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// destination
	dst, err := os.Create(picture.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusCreated, "name: "+name+" email: "+email)
}

func Show(c echo.Context) error {
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team: "+team+", member: "+member)
}
