package server

import (
	"errors"
	"fmt"
	"go-htmx-kvstore/internal/web/data"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Server) handleListKV(c echo.Context) error {
	db := s.db

	token, err := c.Cookie("token")
	if err != nil {
		c.Logger().Error(fmt.Errorf("error getting token: %w", err))
		return c.Redirect(http.StatusFound, "/login")
	}

	if token.Value == "" {
		c.Logger().Error(fmt.Errorf("empty token"))
		return c.Redirect(http.StatusFound, "/login")
	}

	var user *data.User
	result := db.Where("token = ?", token.Value).First(&user)
	if result.Error != nil || user == nil {
		c.Logger().Error(fmt.Errorf("error getting user: %w", result.Error))
		return c.Redirect(http.StatusFound, "/login")
	}

	var kvs []data.KeyValue
	// select * from key_values where user_id = ?
	result = db.Where("user_id = ?", user.ID).Find(&kvs)
	if result.Error != nil {
		c.Logger().Error(fmt.Errorf("error getting key values: %w", result.Error))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Render(200, "pages/kv_list.html", echo.Map{
		"Title":     "Key Values",
		"KeyValues": kvs,
	})
}

func (s *Server) handleNewKV(c echo.Context) error {
	return c.Render(200, "pages/kv_view.html", echo.Map{
		"Title": "New Key Value",
	})
}

func (s *Server) handleEditKV(c echo.Context) error {
	key := c.Param("key")
	value := c.QueryParam("value")

	return c.Render(200, "kv_edit", echo.Map{
		"Key":   key,
		"Value": value,
	})
}

func (s *Server) handleViewKV(c echo.Context) error {
	db := s.db

	userID, err := c.Cookie("user_id")
	if err != nil {
		// internal server error
		c.Logger().Error(fmt.Errorf("error getting user id: %w", err))
		return c.NoContent(http.StatusInternalServerError)
	}

	key := c.Param("key")
	var kv data.KeyValue
	result := db.Where("user_id = ? AND key = ?", userID.Value, key).First(&kv)
	if result.Error != nil {
		c.Logger().Error(fmt.Errorf("error getting key value: %w", result.Error))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Render(200, "kv_view", echo.Map{
		"Key":   key,
		"Value": kv.Value,
	})
}

func (s *Server) handleCreateKV(c echo.Context) error {
	db := s.db

	key := c.FormValue("key")
	value := c.FormValue("value")

	if key == "" || value == "" {
		return c.Render(200, "alert_kv_empty", echo.Map{
			"Key":   key,
			"Value": value,
		})
	}

	userID, err := c.Cookie("user_id")
	if err != nil {
		// internal server error
		c.Logger().Error(fmt.Errorf("error getting user id: %w", err))
		return c.NoContent(http.StatusInternalServerError)
	}

	// if key exists, return error
	var kv data.KeyValue
	result := db.Where("key = ? and user_id = ?", key, userID.Value).First(&kv)
	if result.Error == nil {
		return c.Render(200, "alert_kv_exists", echo.Map{
			"Key": key,
		})
	}

	uid, err := strconv.ParseUint(userID.Value, 10, 64)
	if err != nil {
		c.Logger().Error(fmt.Errorf("error converting user id to uint64: %w", err))
		return c.NoContent(http.StatusInternalServerError)
	}

	result = db.Create(&data.KeyValue{
		Key:    key,
		Value:  value,
		UserID: uint(uid),
	})
	if result.Error != nil {
		c.Logger().Error(fmt.Errorf("error creating key value: %w", result.Error))
		return c.Render(200, "alert_generic_error", echo.Map{
			"Error": result.Error.Error(),
		})

	}

	return c.Render(200, "alert_kv_saved", echo.Map{
		"Key":   key,
		"Value": value,
	})
}

func (s *Server) handleUpdateKV(c echo.Context) error {
	db := s.db

	key := c.Param("key")
	value := c.FormValue("value")

	c.Logger().Info(fmt.Sprintf("key: %s, value: %s", key, value))

	if key == "" || value == "" {
		return c.Render(200, "alert_kv_empty", echo.Map{
			"Key":   key,
			"Value": value,
		})
	}

	userID, err := c.Cookie("user_id")
	if err != nil {
		// internal server error
		c.Logger().Error(fmt.Errorf("error getting user id: %w", err))
		return c.NoContent(http.StatusInternalServerError)
	}

	result := db.Model(&data.KeyValue{}).Where("key = ? and user_id = ?", key, userID.Value).Update("value", value)
	if result.Error != nil {
		c.Logger().Error(fmt.Errorf("error updating key value: %w", result.Error))
		return c.Render(200, "alert_generic_error", echo.Map{
			"Error": result.Error.Error(),
		})
	}

	return c.Render(200, "alert_kv_saved", echo.Map{
		"Key":   key,
		"Value": value,
	})
}

func (s *Server) handleDeleteKV(c echo.Context) error {
	db := s.db

	key := c.Param("key")

	userID, err := c.Cookie("user_id")
	if err != nil {
		// internal server error
		c.Logger().Error(fmt.Errorf("error getting user id: %w", err))
		return c.NoContent(http.StatusInternalServerError)
	}

	result := db.Where("key = ? and user_id = ?", key, userID.Value).Delete(&data.KeyValue{})
	if result.Error != nil {
		c.Logger().Error(fmt.Errorf("error deleting key value: %w", result.Error))
		return c.Render(200, "alert_generic_error", echo.Map{
			"Error": result.Error.Error(),
		})
	}

	return c.NoContent(200)
}

func (s *Server) handleViewSignup(c echo.Context) error {
	return c.Render(200, "pages/signup.html", echo.Map{
		"Title": "Sign Up",
	})
}

func (s *Server) handleSignup(c echo.Context) error {
	db := s.db
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Render(200, "alert_user_empty", nil)
	}

	result := db.Create(&data.User{
		Username: username,
		Password: password,
	})
	if result.Error != nil {
		c.Logger().Error(fmt.Errorf("error creating user: %w", result.Error))
		return c.Render(200, "alert_user_exists", nil)
	}

	c.Response().Header().Set("HX-Location", "/login")
	return c.NoContent(200)
}

func (s *Server) handleViewLogin(c echo.Context) error {
	db := s.db
	token, err := c.Cookie("token")
	if err == nil && token.Value != "" {
		var user data.User
		result := db.Where("token = ?", token.Value).First(&user)
		if result.Error != nil {
			c.Render(200, "alert_generic_error", echo.Map{
				"Error": result.Error.Error(),
			})
		}

		if user.Token == token.Value {
			return c.Redirect(http.StatusFound, "/kv")
		}
	}

	return c.Render(200, "pages/login.html", echo.Map{
		"Title": "Login",
	})
}

func (s *Server) handleLogin(c echo.Context) error {
	db := s.db
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Render(200, "alert_user_empty", nil)
	}

	var user data.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		msg := result.Error.Error()
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			msg = "User not found"
		}
		c.Logger().Error(fmt.Errorf("error getting user: %w", result.Error))
		return c.Render(200, "alert_generic_error", echo.Map{
			"Error": msg,
		})
	}

	if user.Password != password {
		return c.Render(200, "alert_wrong_password", echo.Map{
			"Username": username,
		})
	}

	// generate random token
	token := uuid.Must(uuid.NewUUID()).String()
	user.Token = token

	result = db.Save(&user)
	if result.Error != nil {
		return c.Render(200, "alert_generic_error", echo.Map{
			"Error": result.Error.Error(),
		})
	}

	// https://htmx.org/essays/web-security-basics-with-htmx/#secure-your-cookies
	// Set-Cookie header instructs browser to send cookie in subsequent requests
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	// create new cookie for user id
	cookie = &http.Cookie{
		Name:     "user_id",
		Value:    fmt.Sprint(user.ID),
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	c.Response().Header().Set("HX-Location", "/kv")
	return c.NoContent(200)
}

func (s *Server) handleLogout(c echo.Context) error {
	db := s.db
	tok, err := c.Cookie("token")
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	result := db.Where("token = ?", tok.Value).Delete(&data.User{})
	if result.Error != nil {
		return c.Render(200, "alert_generic_error", echo.Map{
			"Error": result.Error.Error(),
		})
	}

	tok = &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(tok)

	userID := &http.Cookie{
		Name:     "user_id",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(userID)

	c.Response().Header().Set("HX-Location", "/login")
	return c.NoContent(200)
}
