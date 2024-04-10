package main

import (
	"fmt"
	mw "go-htmx-kvstore/internal/middleware"
	"go-htmx-kvstore/web/data"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var registeredUsers = data.Users{
	"admin": {
		Username: "admin",
		Password: "admin",
	},
}
var inMemKvs = data.KeyValues{
	{Key: "key1", Value: "value1"},
	{Key: "key2", Value: "value2"},
}

// https://go.dev/doc/articles/wiki/
func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())
	mw.MustCompileTemplates(e)

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
		for _, u := range registeredUsers {
			if string(u.Token) == token.Value {
				user = u
				break
			}
		}

		if user == nil {
			c.Logger().Error(fmt.Errorf("user not found"))
			return c.Redirect(http.StatusFound, "/login")
		}

		return c.Render(200, "pages/kv_list.html", echo.Map{
			"Title":     "Key Values",
			"KeyValues": inMemKvs,
		})
	})

	e.GET("/kv/new", func(c echo.Context) error {
		return c.Render(200, "pages/kv_view.html", echo.Map{
			"Title": "New Key Value",
		})
	})

	e.GET("/kv/:key", func(c echo.Context) error {
		key := c.Param("key")

		var kv *data.KeyValue
		for _, k := range inMemKvs {
			if k.Key == key {
				kv = k
				break
			}
		}

		if kv == nil {
			return c.Render(404, "pages/not_found.html", echo.Map{
				"Title": "Not Found",
			})
		}

		return c.Render(200, "pages/kv_view.html", echo.Map{
			"Title": "Edit Key Value",
			"Key":   kv.Key,
			"Value": kv.Value,
		})
	})

	e.POST("/kv", func(c echo.Context) error {
		key := c.FormValue("key")
		value := c.FormValue("value")

		if key == "" || value == "" {
			return c.Render(200, "kv_empty", echo.Map{
				"Key":   key,
				"Value": value,
			})
		}

		var exists bool
		for _, kv := range inMemKvs {
			if kv.Key == key {
				kv.Value = value
				exists = true
				break
			}
		}

		if !exists {
			inMemKvs = append(inMemKvs, &data.KeyValue{
				Key:   key,
				Value: value,
			})
		}

		return c.Render(200, "kv_created", nil)
	})

	e.DELETE("/kv/:key", func(c echo.Context) error {
		key := c.Param("key")

		for i, kv := range inMemKvs {
			if kv.Key == key {
				inMemKvs = append(inMemKvs[:i], inMemKvs[i+1:]...)
				break
			}
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

		_, exists := registeredUsers[username]
		if exists {
			return c.Render(200, "alert_user_exists", nil)
		}
		registeredUsers[username] = &data.User{
			Username: username,
			Password: password,
		}

		c.Response().Header().Set("HX-Location", "/login")
		return c.NoContent(200)
	})

	e.GET("/login", func(c echo.Context) error {
		token, err := c.Cookie("token")
		if err == nil && token.Value != "" {
			// check registered users for token
			for _, user := range registeredUsers {
				if string(user.Token) == token.Value {
					return c.Redirect(http.StatusFound, "/kv")
				}
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

		user, exists := registeredUsers[username]
		if !exists || user.Password != password {
			return c.Render(200, "alert_wrong_password", echo.Map{
				"Username": username,
			})
		}

		// generate random token
		token := uuid.Must(uuid.NewUUID()).String()
		user.Token = token

		// https://htmx.org/essays/web-security-basics-with-htmx/#secure-your-cookies
		// Set-Cookie header instructs browser to send cookie in subsequent requests
		cookie := &http.Cookie{
			Name:     "token",
			Value:    token,
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

		for _, user := range registeredUsers {
			if string(user.Token) == cookie.Value {
				user.Token = ""
				break
			}

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

	// e.POST("/start-session", func(c echo.Context) error {
	// 	username := c.FormValue("username")
	//
	// 	_, exists := registeredUsers[username]
	// 	if !exists {
	// 		registeredUsers[username] = &UserData{}
	// 	}
	//
	// 	e.Logger.Infof("user=%s; new=%v\n", username, !exists)
	//
	// 	// https://htmx.org/essays/web-security-basics-with-htmx/#secure-your-cookies
	// 	// Set-Cookie header instructs browser to send cookie in subsequent requests
	// 	cookie := new(http.Cookie)
	// 	cookie.Name = "username"
	// 	cookie.Value = username
	// 	cookie.HttpOnly = true
	// 	c.SetCookie(cookie)
	//
	// 	return c.Render(200, "index.html", echo.Map{
	// 		"Title":           "Home",
	// 		"Username":        username,
	// 		"RegisteredUsers": getRegisteredUsers(),
	// 	})
	// })
	//
	// e.DELETE("/close-session", func(c echo.Context) error {
	//
	// 	user, err := c.Cookie("username")
	// 	if err != nil {
	// 		return c.NoContent(http.StatusInternalServerError)
	// 	}
	// 	delete(registeredUsers, user.Value)
	//
	// 	// delete cookie and return HX header to refresh page
	// 	cookie := new(http.Cookie)
	// 	cookie.Name = "username"
	// 	cookie.Value = ""
	// 	cookie.HttpOnly = true
	// 	cookie.MaxAge = -1
	// 	c.SetCookie(cookie)
	//
	// 	c.Response().Header().Set("HX-Refresh", "true")
	// 	return c.NoContent(http.StatusNoContent)
	// })

	e.Logger.Fatal(e.Start("localhost:1323"))
}

// func getRegisteredUsers() []string {
// 	users := make([]string, 0, len(registeredUsers))
// 	for user := range registeredUsers {
// 		users = append(users, user)
// 	}
// 	return users
// }
