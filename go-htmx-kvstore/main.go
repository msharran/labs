package main

import (
	mw "go-htmx-kvstore/internal/middleware"
	"go-htmx-kvstore/web/data"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var inMemKvs = data.KeyValues{
	{Key: "key1", Value: "value1"},
	{Key: "key2", Value: "value2"},
}

// https://go.dev/doc/articles/wiki/
func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	mw.UseNewTemplateRenderer(e,
		"web/tmpl/partials/*.html",
		"web/tmpl/pages/*.html",
	)

	e.GET("/keyvalues", func(c echo.Context) error {
		// if user is already authenticated, redirect to home
		// cookie, err := c.Cookie("username")
		// if err != nil {
		// 	return c.Render(200, "register.html", nil)
		// }
		// username := cookie.Value

		return c.Render(200, "kv_list.html", echo.Map{
			"Title":     "Key Values",
			"KeyValues": inMemKvs,
		})
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
