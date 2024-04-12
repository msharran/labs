package main

import (
	"errors"
	"fmt"
	mw "go-htmx-kvstore/internal/middleware"
	"go-htmx-kvstore/web/data"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func mustSetupDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("kvstore.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// enable sqlite foreign key support
	result := db.Exec("PRAGMA foreign_keys = ON")
	if result.Error != nil {
		panic(fmt.Errorf("error enabling foreign key support: %w", result.Error))
	}

	err = db.AutoMigrate(&data.User{}, &data.KeyValue{})
	if err != nil {
		panic(fmt.Errorf("failed to migrate database: %w", err))
	}

	return db
}

// https://go.dev/doc/articles/wiki/
func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())
	mw.MustCompileTemplates(e)

	db := mustSetupDb()

	e.GET("/kv", func(c echo.Context) error {
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
		result = db.Find(&kvs)
		if result.Error != nil {
			c.Logger().Error(fmt.Errorf("error getting key values: %w", result.Error))
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.Render(200, "pages/kv_list.html", echo.Map{
			"Title":     "Key Values",
			"KeyValues": kvs,
		})
	})

	e.GET("/kv/new", func(c echo.Context) error {
		return c.Render(200, "pages/kv_view.html", echo.Map{
			"Title": "New Key Value",
		})
	})

	e.GET("/kv/:key/edit", func(c echo.Context) error {
		key := c.Param("key")
		value := c.QueryParam("value")

		return c.Render(200, "kv_edit", echo.Map{
			"Key":   key,
			"Value": value,
		})
	})

	e.GET("/kv/:key/view", func(c echo.Context) error {
		key := c.Param("key")

		var kv data.KeyValue
		result := db.Where("key = ?", key).First(&kv)
		if result.Error != nil {
			c.Logger().Error(fmt.Errorf("error getting key value: %w", result.Error))
			return c.Render(200, "alert_generic_error", echo.Map{
				"Error": result.Error.Error(),
			})
		}

		return c.Render(200, "kv_view", echo.Map{
			"Key":   key,
			"Value": kv.Value,
		})
	})

	e.POST("/kv", func(c echo.Context) error {
		key := c.FormValue("key")
		value := c.FormValue("value")

		if key == "" || value == "" {
			return c.Render(200, "alert_kv_empty", echo.Map{
				"Key":   key,
				"Value": value,
			})
		}

		// if key exists, return error
		var kv data.KeyValue
		result := db.Where("key = ?", key).First(&kv)
		if result.Error == nil {
			return c.Render(200, "alert_kv_exists", echo.Map{
				"Key": key,
			})
		}

		result = db.Create(&data.KeyValue{
			Key:   key,
			Value: value,
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
	})

	e.PUT("/kv/:key", func(c echo.Context) error {
		key := c.Param("key")
		value := c.FormValue("value")

		c.Logger().Info(fmt.Sprintf("key: %s, value: %s", key, value))

		if key == "" || value == "" {
			return c.Render(200, "alert_kv_empty", echo.Map{
				"Key":   key,
				"Value": value,
			})
		}

		// update value when key == key
		result := db.Where("key = ?", key).First(&data.KeyValue{})
		if result.Error != nil {
			c.Logger().Error(fmt.Errorf("error getting key value: %w", result.Error))
			return c.Render(200, "alert_generic_error", echo.Map{
				"Error": result.Error.Error(),
			})
		}

		result = db.Model(&data.KeyValue{}).Where("key = ?", key).Update("value", value)
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
	})

	e.DELETE("/kv/:key", func(c echo.Context) error {
		key := c.Param("key")

		result := db.Where("key = ?", key).Delete(&data.KeyValue{})
		if result.Error != nil {
			c.Logger().Error(fmt.Errorf("error deleting key value: %w", result.Error))
			return c.Render(200, "alert_generic_error", echo.Map{
				"Error": result.Error.Error(),
			})
		}

		return c.NoContent(200)
	})

	// Registering a user
	e.GET("/signup", func(c echo.Context) error {
		return c.Render(200, "pages/signup.html", echo.Map{
			"Title": "Sign Up",
		})
	})

	e.POST("/signup", func(c echo.Context) error {
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
	})

	e.GET("/login", func(c echo.Context) error {
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
	})

	e.POST("/login", func(c echo.Context) error {
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
	})

	e.DELETE("/logout", func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		result := db.Where("token = ?", cookie.Value).Delete(&data.User{})
		if result.Error != nil {
			return c.Render(200, "alert_generic_error", echo.Map{
				"Error": result.Error.Error(),
			})
		}

		cookie = &http.Cookie{
			Name:     "token",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		}
		c.SetCookie(cookie)
		c.Response().Header().Set("HX-Location", "/login")
		return c.NoContent(200)
	})

	e.Logger.Fatal(e.Start("localhost:1323"))
}
