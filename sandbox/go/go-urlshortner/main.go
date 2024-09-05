package main

import (
	"go-urlshortner/internal/storage"
	"go-urlshortner/internal/storage/models"
	"os"
	"time"

	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/lmittmann/tint"
	"github.com/oklog/ulid/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dsn         = "postgres://postgres:pO1mtQQz58pK5AUph6IeKoo4lUKtH4OYfdycakEr2sRP2rwxxC2CbtoP76I8KDEd@localhost:5432/postgres"
	tinyURLBase = "http://localhost:3000"
)

func main() {
	log := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.UnixDate,
		}),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	store := storage.NewStorage(db, log)

	app := fiber.New()

	api := app.Group("/api")
	v1 := api.Group("/v1")

	app.Get("/x/:short", func(c *fiber.Ctx) error {
		short := c.Params("short")
		tinyURL, err := store.GetTinyURLByShort(short)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		log.Info("redirecting to long url", "short", short, "long", tinyURL.Long)
		return c.Redirect(tinyURL.Long)
	})

	v1.Get("/tiny-urls", func(c *fiber.Ctx) error {
		tinyURLs, err := store.GetAllTinyURLs()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(tinyURLs)
	})

	v1.Post("/tiny-urls", func(c *fiber.Ctx) error {
		var tinyURL models.TinyURL
		if err := c.BodyParser(&tinyURL); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if tinyURL.Long == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Long is required",
			})
		}

		// generate a tiny url using ulid
		tinyURL.Short = ulid.Make().String()
		if err := store.CreateTinyURL(&tinyURL); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"Short": tinyURLBase + "/x/" + tinyURL.Short,
		})
	})

	if err = app.Listen(":3000"); err != nil {
		log.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
