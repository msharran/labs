package storage

import (
	"go-urlshortner/internal/storage/models"
	"log/slog"

	"gorm.io/gorm"
)

type Storage struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewStorage(db *gorm.DB, log *slog.Logger) *Storage {

	log.Info("Starting migration")
	// Add migration tables here
	db.AutoMigrate(&models.TinyURL{})
	log.Info("Migration done")

	return &Storage{
		db:  db,
		log: log,
	}
}

// define CRUD methods for TinyURL model
func (s *Storage) CreateTinyURL(t *models.TinyURL) error {
	return s.db.Create(t).Error
}

func (s *Storage) GetTinyURLByShort(short string) (*models.TinyURL, error) {
	var t models.TinyURL
	err := s.db.Where("short = ?", short).First(&t).Error
	return &t, err
}

func (s *Storage) GetTinyURLByLong(long string) (*models.TinyURL, error) {
	var t models.TinyURL
	err := s.db.Where("long = ?", long).First(&t).Error
	return &t, err
}

func (s *Storage) UpdateTinyURL(t *models.TinyURL) error {
	return s.db.Save(t).Error
}

func (s *Storage) DeleteTinyURL(t *models.TinyURL) error {
	return s.db.Delete(t).Error
}

func (s *Storage) GetAllTinyURLs() ([]*models.TinyURL, error) {
	var ts []*models.TinyURL
	err := s.db.Find(&ts).Error
	return ts, err
}
