package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

type DB struct {
	emails map[string]struct{}
}

func main() {
	db := &DB{
		emails: map[string]struct{}{},
	}
	emailsc := make(chan string)

	go func() {
		// update db from emails channel
		for {
			select {
			case e := <-emailsc:
				db.emails[e] = struct{}{}
			}
		}
	}()

	// Create a new engine
	engine := html.New("./views", ".html")

	// Or from an embedded system
	// See github.com/gofiber/embed for examples
	// engine := html.NewFileSystem(http.Dir("./views", ".html"))

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{}, "layouts/main")
	})

	app.Post("/email", func(c *fiber.Ctx) error {
		fmt.Println("post-email req headers", c.GetReqHeaders())

		b := RegisterReq{}
		err := c.BodyParser(&b)
		if err != nil {
			return err
		}

		fmt.Println("post-email before", b.Email, db.emails)
		emailsc <- b.Email
		fmt.Println("post-email after", b.Email, db.emails)

		c.Response().Header.Set("Hx-Trigger", "email-registered")
		return c.SendStatus(200)
	})

	app.Get("/email", func(c *fiber.Ctx) error {
		ee := getEmails(db)
		fmt.Println("get-emails", ee)
		return c.Render("partials/emails", fiber.Map{
			"Emails": ee,
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func getEmails(db *DB) []string {
	ee := make([]string, 0, len(db.emails))
	for k := range db.emails {
		ee = append(ee, k)
	}
	return ee
}

type RegisterReq struct {
	Email string `json:"email"`
}
