package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"todolist/config"
	"todolist/models"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
	db := config.GetDB()
	rows, err := db.Query("SELECT id, name, email FROM users")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		users = append(users, user)

	}
	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	id := c.Param("id")
	db := config.GetDB()
	var user models.User

	err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, "User not found")
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	db := config.GetDB()
	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	id, _ := result.LastInsertId()
	user.ID = int(id)
	return c.JSON(http.StatusCreated, user)
}

func UpdateUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	db := config.GetDB()
	_, err = db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, id)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}



func DeleteUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	db := config.GetDB()
	result, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	return c.NoContent(http.StatusNoContent)
}